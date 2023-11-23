package webview

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego/objc"
	"github.com/mekkanized/go-webview/internal/darwin/cocoa"
)

func init() {
	runtime.LockOSThread()
}

type webview struct {
	options  WebViewOptions
	bindings map[string]interface{}
	mutex    sync.RWMutex

	webview      cocoa.WKWebView
	window       *cocoa.NSWindow
	parentWindow *cocoa.NSWindow
	manager      cocoa.WKUserContentController
}

func NewWithOptions(options WebViewOptions) WebView {
	w := &webview{
		bindings: make(map[string]interface{}),
		options:  options,
	}
	if options.Window != nil {
		w.parentWindow = &cocoa.NSWindow{ID: objc.ID(options.Window)}
	}

	app := cocoa.NSApplication_GetSharedApplication()
	delegate, err := w.createAppDelegate()
	if err != nil {
		panic(fmt.Errorf("failed to create app delegate: %w", err))
	}
	app.SetDelegate(delegate.ID)
	// TODO: Set associated object
	if w.parentWindow != nil {
		w.onApplicationDidFinishLaunching(delegate.ID, app.ID)
	} else {
		app.Run()
	}

	return w
}

func (w *webview) Run() {
	app := cocoa.NSApplication_GetSharedApplication()
	app.Run()
}

func (w *webview) Terminate() {
	stopRunLoop()
}

func (w *webview) Dispatch(f func()) {
	// TODO: Implement
}

func (w *webview) Destroy() {
	// TODO: Implement
}

func (w *webview) Window() unsafe.Pointer {
	return unsafe.Pointer(w.window.ID)
}

func (w *webview) SetTitle(title string) {
	w.window.SetTitle(title)
}

func (w *webview) SetSize(width int, height int, hint Hint) {
	style := cocoa.NSWindowStyleMaskTitled | cocoa.NSWindowStyleMaskClosable | cocoa.NSWindowStyleMaskMiniaturizable
	if hint != HintFixed {
		style |= cocoa.NSWindowStyleMaskResizable
	}
	w.window.SetStyleMask(style)

	size := cocoa.NSSize{
		Width:  float64(width),
		Height: float64(height),
	}

	switch hint {
	case HintMin:
		w.window.SetContentMinSize(size)
	case HintMax:
		w.window.SetContentMaxSize(size)
	default:
		w.window.SetFrame(cocoa.NSRect{
			Origin: cocoa.NSPoint{
				X: 0,
				Y: 0},
			Size: size,
		}, true, false)
	}
	w.window.Center()
}

func (w *webview) Navigate(url string) {
	pool := cocoa.NSAutoreleasePool_new()
	defer pool.Release()

	wrappedUrl := cocoa.NSURL_URLWithString(url)
	request := cocoa.NSURLRequest_requestWithURL(wrappedUrl)
	w.webview.LoadRequest(request)
}

func (w *webview) SetHtml(html string) {
	pool := cocoa.NSAutoreleasePool_new()
	defer pool.Release()

	// TODO: Figure out why html causes a segmentation violation
	w.webview.LoadHTMLString(html, 0)
}

func (w *webview) Init(js string) {
	script := cocoa.WKUserScript_alloc().
		InitWithSource(
			js,
			cocoa.WKUserScriptInjectionTimeAtDocumentStart,
			true,
		)
	w.manager.AddUserScript(script)
}

func (w *webview) Eval(js string) {
	w.webview.EvaluateJavaScript(js, objc.ID(0))
}

func (w *webview) onApplicationDidFinishLaunching(delegateID objc.ID, appID objc.ID) {
	app := cocoa.NSApplication{ID: appID}
	if w.parentWindow == nil {
		// Stop the main run loop so that we can return
		// from the constructor.
		stopRunLoop()
	}

	// Activate the app if it is not bundled.
	// Bundled apps launched from Finder are activated automatically but
	// otherwise not. Activating the app even when it has been launched from
	// Finder does not seem to be harmful but calling this function is rarely
	// needed as proper activation is normally taken care of for us.
	// Bundled apps have a default activation policy of
	// NSApplicationActivationPolicyRegular while non-bundled apps have a
	// default activation policy of NSApplicationActivationPolicyProhibited.
	if !isAppBundled() {
		// SetActivationPolicy must be invoked before
		// ActivateIgnoringOtherApps for activation to work.
		app.SetActivationPolicy(cocoa.NSApplicationActivationPolicyRegular)
		// Activate the app regardless of other active apps.
		// This can be obtrusive so we only do it when necessary.
		app.ActivateIgnoringOtherApps(true)
	}

	w.window = &cocoa.NSWindow{}
	if w.parentWindow == nil {
		window := cocoa.NSWindow_alloc().InitWithContentRect(cocoa.NSRect{
			Origin: cocoa.NSPoint{X: 0, Y: 0},
			Size:   cocoa.NSSize{Width: 0, Height: 0},
		}, cocoa.NSWindowStyleMaskTitled, cocoa.NSBackingStoreBuffered, false)
		w.window = &window
	} else {
		w.window = w.parentWindow
	}

	config := cocoa.WKWebViewConfiguration_new()
	w.manager = config.UserContentController()
	w.webview = cocoa.WKWebView_alloc().InitWithFrame(
		cocoa.NSRect{
			Origin: cocoa.NSPoint{X: 0, Y: 0},
			Size:   cocoa.NSSize{Width: 0, Height: 0},
		}, config,
	)

	if w.options.Debug {
		config.Preferences().SetValue("developerExtrasEnabled", cocoa.NSNumber_NumberWithBool(true).ID)
	}

	uiDelegate := w.createWebkitUIDelegate()
	w.webview.SetUIDelegate(uiDelegate)

	scriptMessageHandler := w.createScriptMessageHandler()
	w.manager.AddScriptMessageHandler(scriptMessageHandler, "external")

	w.Init(`
		window.external = {
			invoke: function (s) {
				window.webkit.messageHandlers.external.postMessage(s);
			}
		}
	`)
	w.window.SetContentView(w.webview.ID)
	w.window.MakeKeyAndOrderFront(0)
}

func isAppBundled() bool {
	bundle := cocoa.NSBundle_GetMainBundle()
	if bundle.ID == 0 {
		return false
	}

	bundlePath := bundle.BundlePath().String()
	return strings.HasSuffix(bundlePath, ".app")
}

func stopRunLoop() {
	app := cocoa.NSApplication_GetSharedApplication()
	app.Stop(0)

	// TODO: Figure out how to create an NSEvent
}

func (w *webview) createAppDelegate() (*cocoa.WebviewAppDelegate, error) {
	// Note: Avoid registering the class name "AppDelegate" as it is the
	// default name in projects created with Xcode, and using the same name
	// causes registerClassPair to crash.
	class, err := objc.RegisterClass(
		"WebviewAppDelegate",
		cocoa.Class_NSResponder,
		[]*objc.Protocol{
			objc.GetProtocol("NSTouchBarProvider"),
		},
		[]objc.FieldDef{},
		[]objc.MethodDef{
			{
				Cmd: cocoa.Sel_applicationShouldTerminateAfterLastWindowClosed,
				Fn: func(id objc.ID, cmd objc.SEL, notification objc.ID) bool {
					return true
				},
			},
			{
				Cmd: cocoa.Sel_applicationShouldTerminate,
				Fn: func(id objc.ID, cmd objc.SEL, sender objc.ID) int {
					// return 2 // NSTerminateLater
					return 1 // NSTerminateNow
				},
			},
			// TODO: Only register the following if we did not initialize with an existing window
			{
				Cmd: cocoa.Sel_applicationDidFinishLaunching,
				Fn: func(self objc.ID, cmd objc.SEL, notification objc.ID) {
					fmt.Println("Received")
					app := notification.Send(cocoa.Sel_object)
					// TODO: Use get_associated_webview instead of taking object from caller
					w.onApplicationDidFinishLaunching(self, app)
				},
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register app delegate class: %w", err)
	}

	res := &cocoa.WebviewAppDelegate{
		ID: objc.ID(class).Send(cocoa.Sel_new),
	}

	return res, nil
}

func (w *webview) createWebkitUIDelegate() objc.ID {
	var res objc.ID
	var err error

	class, err := objc.RegisterClass(
		"WebviewUIDelegate",
		cocoa.Class_NSObject,
		[]*objc.Protocol{
			objc.GetProtocol("WKUIDelegate"),
		},
		[]objc.FieldDef{},
		[]objc.MethodDef{
			{
				Cmd: objc.RegisterName("webView:runOpenPanelWithParameters:initiatedByFrame:completionHandler:"),
				Fn: func(_ objc.ID, _ objc.SEL, _ objc.ID, parameters objc.ID, _ objc.ID, completionHandler objc.ID) {
					allowsMultipleSelection := objc.ID(parameters).Send(objc.RegisterName("allowsMultipleSelection")) != 0
					allowsDirectories := objc.ID(parameters).Send(objc.RegisterName("allowsDirectories")) != 0

					panel := cocoa.NSOpenPanel_openPanel()
					panel.SetCanChooseFiles(true)
					panel.SetCanChooseDirectories(allowsDirectories)
					panel.SetAllowsMultipleSelection(allowsMultipleSelection)
					modalResponse := panel.RunModal()

					var urls objc.ID
					if modalResponse == cocoa.NSModalResponseOK {
						urls = panel.URLs()
					}

					sig := cocoa.NSMethodSignature_signatureWithObjCTypes("v@?@")
					invocation := cocoa.NSInvocation_invocationWithMethodSignature(sig)
					invocation.Send(objc.RegisterName("setTarget"), completionHandler)
					invocation.Send(objc.RegisterName("setArgument:atIndex:"), urls, 1)
					invocation.Send(objc.RegisterName("invoke"))
				},
			},
		})
	if err != nil {
		panic(fmt.Errorf("failed to register webkit ui delegate class: %w", err))
	}
	res = objc.ID(class).Send(cocoa.Sel_new)

	return res
}

func (w *webview) createScriptMessageHandler() objc.ID {
	class, err := objc.RegisterClass(
		"WebviewWKScriptMessageHandler",
		cocoa.Class_NSResponder,
		[]*objc.Protocol{
			objc.GetProtocol("WKScriptMessageHandler"),
		},
		[]objc.FieldDef{},
		[]objc.MethodDef{
			{
				Cmd: objc.RegisterName("userContentController:didReceiveScriptMessage:"),
				Fn: func(self objc.ID, cmd objc.SEL, _ objc.ID, msg objc.ID) {
					body := msg.Send(objc.RegisterName("body"))
					// str := body.Send(objc.RegisterName("UTF8String"))
					w.onMessage(cocoa.NSString{ID: body}.String())
				},
			},
		})
	if err != nil {
		panic(fmt.Errorf("failed to register webkit ui delegate class: %w", err))
	}
	res := objc.ID(class).Send(cocoa.Sel_new)

	return res
}
