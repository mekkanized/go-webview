//go:build linux

package webview

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/mekkanized/go-webview/internal/linux/webkitgtk"
)

var webkit webkitgtk.Context = nil

type webview struct {
	options  WebViewOptions
	bindings map[string]interface{}
	mutex    sync.RWMutex

	webview webkitgtk.WebKitWebView
	window  webkitgtk.GtkWindow
}

// NewWithOptions creates a new webview using the provided options.
func NewWithOptions(options WebViewOptions) WebView {
	w := &webview{
		bindings: make(map[string]interface{}),
		options:  options,
	}

	if webkit == nil {
		var err error
		webkit, err = webkitgtk.NewDefaultContext()
		if err != nil {
			panic(fmt.Errorf("failed to create webkit context: %w", err))
		}
		if err = webkit.LoadFunctions(); err != nil {
			panic(fmt.Errorf("failed to load webkit functions: %w", err))
		}

		fmt.Printf("WebKitGTK version %d.%d.%d\n", webkit.WebKitGetMajorVersion(), webkit.WebKitGetMinorVersion(), webkit.WebKitGetMicroVersion())
	}

	// Initialize GTK
	if !webkit.GtkInitCheck() {
		panic(fmt.Errorf("failed to initialize GTK"))
	}

	w.window = webkitgtk.GtkWindow(options.Window)
	if w.window == webkitgtk.GtkWindow(webkitgtk.NULLPTR) {
		w.window = webkitgtk.GtkWindow(webkit.GtkWindowNew(webkitgtk.GTK_WINDOW_TOPLEVEL))
	}

	webkit.GSignalConnectData(webkitgtk.GtkWidget(w.window), "destroy", func() {
		w.Terminate()
	}, webkitgtk.NULLPTR, nil, webkitgtk.G_CONNECT_DEFAULT)

	// Initialize webview widget
	w.webview = webkitgtk.WebKitWebView(webkit.WebKitWebViewNew())
	manager := webkit.WebKitWebViewGetUserContentManager(w.webview)

	// Setup binding callbacks
	webkit.GSignalConnectData(webkitgtk.GtkWidget(manager), "script-message-received::external", func(manager webkitgtk.WebKitUserContentManager, result webkitgtk.WebKitJavascriptResult, arg uintptr) {
		s, err := getStringFromJsResult(result)
		if err != nil {
			fmt.Printf("RPC call failed: %v\n", fmt.Errorf("failed to get string from js result: %w", err))
		}

		w.onMessage(s)
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
	return *(*unsafe.Pointer)(unsafe.Pointer(&w.window))
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
		geometry := webkitgtk.GdkGeometry{}

		var h webkitgtk.GdkWindowHints
		switch hint {
		default:
			fallthrough
		case HintMin:
			h = webkitgtk.GDK_HINT_MIN_SIZE
			geometry.MinHeight = int32(height)
			geometry.MinWidth = int32(width)
		case HintMax:
			h = webkitgtk.GDK_HINT_MAX_SIZE
			geometry.MaxHeight = int32(height)
			geometry.MaxWidth = int32(width)
		}
		webkit.GtkWindowSetGeometryHints(w.window, webkitgtk.GtkWidget(webkitgtk.NULLPTR), geometry, h)
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

func getStringFromJsResult(r webkitgtk.WebKitJavascriptResult) (string, error) {
	var str string

	if webkit.WebKitGetMajorVersion() >= 2 && webkit.WebKitGetMinorVersion() >= 22 {
		value := webkit.WebKitJavascriptResultGetJsValue(webkitgtk.WebKitJavascriptResult(r))
		str = webkit.JsCValueToString(value)
	} else {
		return "", fmt.Errorf("unsupported webkit version: %d.%d.%d", webkit.WebKitGetMajorVersion(), webkit.WebKitGetMinorVersion(), webkit.WebKitGetMicroVersion())
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
