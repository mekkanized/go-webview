//go:build !windows

package webview

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/Mekkanized/go-webview/internal/webkitgtk"
	"github.com/pkg/errors"
)

var webkit webkitgtk.Context = nil

type webview struct {
	webview webkitgtk.WebKitWebView
	window  webkitgtk.GtkWindow
}

// New calls NewWindow to create a new window and a new webview instance. If debug
// is non-zero - developer tools will be enabled (if the platform supports them).
func New(debug bool) WebView {
	return NewWithOptions(WebViewOptions{
		Debug: debug,
	})
}

// NewWindow creates a new webview using an existing window.
//
// Deprecated: Use NewWithOptions.
func NewWindow(debug bool, window unsafe.Pointer) WebView {
	return NewWithOptions(WebViewOptions{
		Debug:  debug,
		Window: window,
	})
}

// NewWithOptions creates a new webview using the provided options.
func NewWithOptions(options WebViewOptions) WebView {
	w := &webview{}

	if webkit == nil {
		var err error
		webkit, err = webkitgtk.NewDefaultContext()
		if err != nil {
			panic(errors.Wrap(err, "failed to create webkit context"))
		}
		if err = webkit.LoadFunctions(); err != nil {
			panic(errors.Wrap(err, "failed to load webkit functions"))
		}

		fmt.Printf("WebKitGTK version %d.%d.%d\n", webkit.WebKitGetMajorVersion(), webkit.WebKitGetMinorVersion(), webkit.WebKitGetMicroVersion())
	}

	// Initialize GTK
	if !webkit.GtkInitCheck() {
		panic(errors.New("failed to initialize GTK"))
	}

	w.window = webkitgtk.GtkWindow(options.Window)
	if w.window == webkitgtk.GtkWindow(webkitgtk.NULLPTR) {
		w.window = webkitgtk.GtkWindow(webkit.GtkWindowNew(webkitgtk.GTK_WINDOW_TOPLEVEL))
	}
	webkit.GSignalConnectData(webkitgtk.GtkWidget(w.window), "destroy", func() {
		webkit.GtkMainQuit()
	}, webkitgtk.NULLPTR, nil, webkitgtk.G_CONNECT_DEFAULT)

	// Initialize webview widget
	w.webview = webkitgtk.WebKitWebView(webkit.WebKitWebViewNew())
	manager := webkit.WebKitWebViewGetUserContentManager(w.webview)

	webkit.GSignalConnectData(webkitgtk.GtkWidget(manager), "script-message-received::external", func(manager uintptr, result uintptr, arg uintptr) {
		s, err := getStringFromJsResult(result)
		if err != nil {
			// TODO: Handle this
			panic(errors.Wrap(err, "script-message-received::external failed"))
		}

		// TODO: Use message appropriately
		fmt.Printf("Received message from script-message-received::external -> \"%s\"", s)
	}, webkitgtk.NULLPTR, nil, webkitgtk.G_CONNECT_DEFAULT)

	webkit.WebKitUserContentManagerRegisterScriptMessageHandler(manager, "external")

	w.Init("window.external={invoke:function(s){window.webkit.messageHandlers.external.postMessage(s);}}")

	webkit.GtkContainerAdd(webkitgtk.GtkContainer(w.window), webkitgtk.GtkWidget(w.webview))
	webkit.GtkWidgetGrabFocus(webkitgtk.GtkWidget(w.webview))

	settings := webkit.WebKitWebViewGetSettings(w.webview)
	webkit.WebKitSettingsSetJavascriptCanAccessClipboard(settings, true)
	if options.Debug {
		webkit.WebKitSettingsSetEnableWriteConsoleMessagesToStdout(settings, true)
		webkit.WebKitSettingsSetEnableDeveloperExtras(settings, true)
	}

	webkit.GtkWidgetShowAll(webkitgtk.GtkWidget(w.window))

	return w
}

func (w *webview) Run() {
	webkit.GtkMain()
}

func (w *webview) Terminate() {
	webkit.GtkMainQuit()
}

func (w *webview) Dispatch(f func()) {
	webkit.GIdleAddFull(webkitgtk.G_PRIORITY_HIGH_IDLE, func(userData uintptr) bool {
		f()
		return webkitgtk.G_SOURCE_REMOVE
	}, webkitgtk.NULLPTR, nil)
}

func (w *webview) Destroy() {
	webkit = nil
}

func (w *webview) Window() unsafe.Pointer {
	return unsafe.Pointer(w.window)
}

func (w *webview) SetTitle(title string) {
	webkit.GtkWindowSetTitle(w.window, title)
}

func (w *webview) SetSize(width int, height int, hint Hint) {
	webkit.GtkWindowSetResizable(w.window, hint != HintFixed)
	switch hint {
	case HintNone:
		webkit.GtkWindowResize(w.window, width, height)
	case HintFixed:
		webkit.GtkWidgetSetSizeRequest(webkitgtk.GtkWidget(w.window), width, height)
	default:
		// geometry := webkitgtk.GDKGeometry{
		// 	MinWidth:  int32(width),
		// 	MaxWidth:  int32(width),
		// 	MinHeight: int32(height),
		// 	MaxHeight: int32(height),
		// }

		// h := webkitloader.GDKHintMaxSize
		// if hint == HintMin {
		// 	h = webkitloader.GDKHintMinSize
		// }
		// webkitloader.GTKWindowSetGeometryHints(w.window, nil, unsafe.Pointer(&geometry), h)
	}
}

func (w *webview) Navigate(url string) {
	webkit.WebKitWebViewLoadURI(w.webview, url)
}

func (w *webview) SetHtml(html string) {
	webkit.WebKitWebViewLoadHTML(w.webview, html, "")
}

func (w *webview) Init(js string) {
	manager := webkit.WebKitWebViewGetUserContentManager(w.webview)

	script := webkit.WebKitUserScriptNew(js, webkitgtk.WEBKIT_USER_CONTENT_INJECT_TOP_FRAME, webkitgtk.WEBKIT_USER_SCRIPT_INJECT_AT_DOCUMENT_START, "", "")
	webkit.WebKitUserContentManagerAddScript(manager, script)
}

func (w *webview) Eval(js string) {
	webkit.WebKitWebViewRunJavascript(w.webview, js, webkitgtk.GCancellable(webkitgtk.NULLPTR), nil, webkitgtk.NULLPTR)
}

func (w *webview) Bind(name string, f interface{}) error {
	v := reflect.ValueOf(f)
	if v.Kind() != reflect.Func {
		return errors.New("only functions can be bound")
	}

	if n := v.Type().NumOut(); n > 2 {
		return errors.New("function may only return a value or a value+error")
	}
	return nil
}

func getStringFromJsResult(r uintptr) (string, error) {
	var str string

	if webkit.WebKitGetMajorVersion() >= 2 && webkit.WebKitGetMinorVersion() >= 22 {
		value := webkit.WebKitJavascriptResultGetJsValue(webkitgtk.WebKitJavascriptResult(r))
		str = webkit.JsCValueToString(value)
	} else {
		return "", errors.New("unsupported webkit version")
		// ctx := webkitloader.WebkitJavascriptResultGetGlobalContext(r)
		// value := webkitloader.WebkitJavascriptResultGetValue(r)
		// c_str := webkitloader.JSValueToStringCopy(ctx, value, nil)
		// _ = webkitloader.JSStringGetMaximumUTF8CStringSize(c_str)
		// TODO: Fix and test
		// s := webkitloader.GNew(char, n)
		// webkitloader.JSStringGetUTF8CString(c_str, s, n)
		// str = webkitloader.GoString(s)
	}

	return str, nil
}
