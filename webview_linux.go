//go:build linux

package webview

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sync"
	"unsafe"

	"github.com/mekkanized/go-webview/internal/webkitgtk"
	"github.com/pkg/errors"
)

var webkit webkitgtk.Context = nil

type webview struct {
	webview webkitgtk.WebKitWebView
	window  webkitgtk.GtkWindow

	bindings map[string]interface{}
	mutex    sync.RWMutex
}

// New calls NewWindow to create a new window and a new webview instance. If debug
// is non-zero - developer tools will be enabled (if the platform supports them).
func New(debug bool) WebView {
	return NewWithOptions(WebViewOptions{
		Debug: debug,
	})
}

// NewWithOptions creates a new webview using the provided options.
func NewWithOptions(options WebViewOptions) WebView {
	w := &webview{
		bindings: make(map[string]interface{}),
	}

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
		w.Terminate()
	}, webkitgtk.NULLPTR, nil, webkitgtk.G_CONNECT_DEFAULT)

	// Initialize webview widget
	w.webview = webkitgtk.WebKitWebView(webkit.WebKitWebViewNew())
	manager := webkit.WebKitWebViewGetUserContentManager(w.webview)

	// Setup binding callbacks
	webkit.GSignalConnectData(webkitgtk.GtkWidget(manager), "script-message-received::external", func(manager webkitgtk.WebKitUserContentManager, result webkitgtk.WebKitJavascriptResult, arg uintptr) {
		s, err := getStringFromJsResult(result)
		if err != nil {
			fmt.Printf("RPC call failed: %v\n", errors.Wrap(err, "failed to get string from js result"))
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

func (w *webview) Bind(name string, f interface{}) error {
	v := reflect.ValueOf(f)
	if v.Kind() != reflect.Func {
		return errors.New("only functions can be bound")
	}

	if n := v.Type().NumOut(); n > 2 {
		return errors.New("function may only return a value or a value+error")
	}

	w.mutex.Lock()
	w.bindings[name] = f
	w.mutex.Unlock()

	js := fmt.Sprintf(`(function() { var name = '%s';
		var RPC = window._rpc = (window._rpc || {nextSeq: 1});
		window[name] = function() {
			var seq = RPC.nextSeq++;
			var promise = new Promise(function(resolve, reject) {
				RPC[seq] = {
					resolve: resolve,
					reject: reject
				};
			});
			window.external.invoke(JSON.stringify({
				id: seq,
				method: name,
				params: Array.prototype.slice.call(arguments),
			}));
			return promise;
		};
	})();`, name)
	w.Init(js)
	w.Eval(js)

	return nil
}

type rpcMessage struct {
	ID     int               `json:"id"`
	Method string            `json:"method"`
	Params []json.RawMessage `json:"params"`
}

func (w *webview) onMessage(msg string) {
	var req rpcMessage
	if err := json.Unmarshal([]byte(msg), &req); err != nil {
		log.Printf("invalid RPC message: %v", err)
		return
	}

	defer w.Eval(fmt.Sprintf(`delete window._rpc[%d];`, req.ID))

	res, err := w.callBinding(req)
	if err != nil {
		w.Eval(fmt.Sprintf(`window._rpc[%d].reject(%s);`, req.ID, err.Error()))
		return
	}

	serRes, err := json.Marshal(res)
	if err != nil {
		w.Eval(fmt.Sprintf(`window._rpc[%d].reject(%s);`, req.ID, err.Error()))
		return
	}

	w.Eval(fmt.Sprintf(`window._rpc[%d].resolve(%s);`, req.ID, serRes))
}

func (w *webview) callBinding(req rpcMessage) (interface{}, error) {
	w.mutex.RLock()
	f, ok := w.bindings[req.Method]
	w.mutex.RUnlock()
	if !ok {
		return nil, nil
	}

	v := reflect.ValueOf(f)
	isVariadic := v.Type().IsVariadic()
	numIn := v.Type().NumIn()
	if (isVariadic && len(req.Params) < numIn-1) || (!isVariadic && len(req.Params) != numIn) {
		return nil, errors.New("function arguments mismatch")
	}

	args := []reflect.Value{}
	for i := range req.Params {
		var arg reflect.Value
		if isVariadic && i >= numIn-1 {
			arg = reflect.New(v.Type().In(numIn - 1).Elem())
		} else {
			arg = reflect.New(v.Type().In(i))
		}
		if err := json.Unmarshal(req.Params[i], arg.Interface()); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal argument")
		}
		args = append(args, arg.Elem())
	}

	errorType := reflect.TypeOf((*error)(nil)).Elem()
	res := v.Call(args)
	switch len(res) {
	case 0:
		// No results from the function, just return nil
		return nil, nil
	case 1:
		// One result may be a value, or an error
		if res[0].Type().Implements(errorType) {
			if res[0].Interface() != nil {
				return nil, res[0].Interface().(error)
			}
			return nil, nil
		}
		return res[0].Interface(), nil
	case 2:
		// Two results: first one is value, second one is error
		if !res[1].Type().Implements(errorType) {
			return nil, errors.New("second return value must be an error")
		}
		if res[1].Interface() != nil {
			return res[0].Interface(), nil
		}
		return res[0].Interface(), res[1].Interface().(error)
	default:
		return nil, errors.New("unexpected number of return values")
	}
}

func getStringFromJsResult(r webkitgtk.WebKitJavascriptResult) (string, error) {
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
