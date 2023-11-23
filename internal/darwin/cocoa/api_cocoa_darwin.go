package cocoa

import (
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/objc"
)

var (
	objc_setAssociatedObject func(obj objc.ID, key objc.ID, value objc.ID, policy uint)
)

var (
	class_NSAutoreleasePool objc.Class
	class_NSApplication     objc.Class
	class_NSBundle          objc.Class
	class_NSEvent           objc.Class
	class_NSInvocation      objc.Class
	Class_NSObject          objc.Class
	Class_NSResponder       objc.Class
	class_NSString          objc.Class
	class_NSMethodSignature objc.Class
	class_NSWindow          objc.Class
	class_NSURL             objc.Class
	class_NSURLRequest      objc.Class

	Class_WebviewAppDelegate objc.Class
)

var (
	sel_activateIgnoringOtherApps                       objc.SEL
	sel_alloc                                           objc.SEL
	Sel_applicationDidFinishLaunching                   objc.SEL
	Sel_applicationShouldTerminateAfterLastWindowClosed objc.SEL
	Sel_applicationShouldTerminate                      objc.SEL
	sel_bundlePath                                      objc.SEL
	sel_initWithContentRect_styleMask_backing_defer_    objc.SEL
	sel_length                                          objc.SEL
	sel_mainBundle                                      objc.SEL
	sel_makeKeyAndOrderFront                            objc.SEL
	Sel_new                                             objc.SEL
	Sel_object                                          objc.SEL
	sel_postEvent_atStart_                              objc.SEL
	sel_run                                             objc.SEL
	sel_setActivationPolicy                             objc.SEL
	sel_setContentView                                  objc.SEL
	sel_setDelegate                                     objc.SEL
	sel_sharedApplication                               objc.SEL
	sel_stop                                            objc.SEL
	sel_UTF8String                                      objc.SEL
	sel_initWithUTF8String                              objc.SEL
	sel_setFrameDisplayAnimate                          objc.SEL
	sel_setContentMinSize                               objc.SEL
	sel_setContentMaxSize                               objc.SEL
	sel_center                                          objc.SEL
	sel_setStyleMask                                    objc.SEL
	sel_loadHTMLStringBaseURL                           objc.SEL
	sel_stringWithUTF8String                            objc.SEL
	sel_release                                         objc.SEL

	sel_invocationWithMethodSignature      objc.SEL
	sel_setSelector                        objc.SEL
	sel_setTarget                          objc.SEL
	sel_setArgumentAtIndex                 objc.SEL
	sel_getReturnValue                     objc.SEL
	sel_invoke                             objc.SEL
	sel_invokeWithTarget                   objc.SEL
	sel_instanceMethodSignatureForSelector objc.SEL
	sel_signatureWithObjCTypes             objc.SEL
)

func init() {
	_, err := purego.Dlopen("/System/Library/Frameworks/Cocoa.framework/Cocoa", purego.RTLD_LAZY|purego.RTLD_GLOBAL)
	if err != nil {
		panic(err)
	}
	_, err = purego.Dlopen("/System/Library/Frameworks/WebKit.framework/WebKit", purego.RTLD_LAZY|purego.RTLD_GLOBAL)
	if err != nil {
		panic(err)
	}

	class_NSAutoreleasePool = objc.GetClass("NSAutoreleasePool")
	class_NSApplication = objc.GetClass("NSApplication")
	class_NSBundle = objc.GetClass("NSBundle")
	class_NSEvent = objc.GetClass("NSEvent")
	class_NSInvocation = objc.GetClass("NSInvocation")
	Class_NSObject = objc.GetClass("NSObject")
	Class_NSResponder = objc.GetClass("NSResponder")
	class_NSString = objc.GetClass("NSString")
	class_NSMethodSignature = objc.GetClass("NSMethodSignature")
	class_NSWindow = objc.GetClass("NSWindow")
	class_NSURL = objc.GetClass("NSURL")
	class_NSURLRequest = objc.GetClass("NSURLRequest")

	sel_activateIgnoringOtherApps = objc.RegisterName("activateIgnoringOtherApps:")
	sel_alloc = objc.RegisterName("alloc")
	Sel_applicationDidFinishLaunching = objc.RegisterName("applicationDidFinishLaunching:")
	Sel_applicationShouldTerminateAfterLastWindowClosed = objc.RegisterName("applicationShouldTerminateAfterLastWindowClosed:")
	Sel_applicationShouldTerminate = objc.RegisterName("applicationShouldTerminate:")
	sel_bundlePath = objc.RegisterName("bundlePath")
	sel_initWithContentRect_styleMask_backing_defer_ = objc.RegisterName("initWithContentRect:styleMask:backing:defer:")
	sel_length = objc.RegisterName("length")
	sel_mainBundle = objc.RegisterName("mainBundle")
	sel_makeKeyAndOrderFront = objc.RegisterName("makeKeyAndOrderFront:")
	Sel_new = objc.RegisterName("new")
	Sel_object = objc.RegisterName("object")
	sel_postEvent_atStart_ = objc.RegisterName("postEvent:atStart:")
	sel_run = objc.RegisterName("run")
	sel_setActivationPolicy = objc.RegisterName("setActivationPolicy:")
	sel_setContentView = objc.RegisterName("setContentView:")
	sel_setDelegate = objc.RegisterName("setDelegate:")
	sel_sharedApplication = objc.RegisterName("sharedApplication")
	sel_stop = objc.RegisterName("stop:")
	sel_UTF8String = objc.RegisterName("UTF8String")
	sel_initWithUTF8String = objc.RegisterName("initWithUTF8String:")
	sel_setFrameDisplayAnimate = objc.RegisterName("setFrame:display:animate:")
	sel_setContentMinSize = objc.RegisterName("setContentMinSize:")
	sel_setContentMaxSize = objc.RegisterName("setContentMaxSize:")
	sel_center = objc.RegisterName("center")
	sel_setStyleMask = objc.RegisterName("setStyleMask:")
	sel_stringWithUTF8String = objc.RegisterName("stringWithUTF8String:")
	sel_release = objc.RegisterName("release")

	sel_invocationWithMethodSignature = objc.RegisterName("invocationWithMethodSignature:")
	sel_setSelector = objc.RegisterName("setSelector:")
	sel_setTarget = objc.RegisterName("setTarget:")
	sel_setArgumentAtIndex = objc.RegisterName("setArgument:atIndex:")
	sel_getReturnValue = objc.RegisterName("getReturnValue:")
	sel_invoke = objc.RegisterName("invoke")
	sel_invokeWithTarget = objc.RegisterName("invokeWithTarget:")
	sel_instanceMethodSignatureForSelector = objc.RegisterName("instanceMethodSignatureForSelector:")
	sel_signatureWithObjCTypes = objc.RegisterName("signatureWithObjCTypes:")

	class_WKWebView = objc.GetClass("WKWebView")
	class_WKWebViewConfiguration = objc.GetClass("WKWebViewConfiguration")
	class_WKUserScript = objc.GetClass("WKUserScript")

	sel_UserContentController = objc.RegisterName("userContentController")
}

type CGFloat float64

type NSInteger = int
type NSUInteger = uint

type NSApplicationActivationPolicy NSInteger

var (
	NSApplicationActivationPolicyRegular NSApplicationActivationPolicy = 0
)

type NSWindowStyleMask NSUInteger

var (
	NSWindowStyleMaskTitled         NSWindowStyleMask = 1 << 0
	NSWindowStyleMaskClosable       NSWindowStyleMask = 1 << 1
	NSWindowStyleMaskMiniaturizable NSWindowStyleMask = 1 << 2
	NSWindowStyleMaskResizable      NSWindowStyleMask = 1 << 3
)

type NSApplication struct {
	objc.ID
}

type NSBackingStoreType NSUInteger

type NSModalResponse NSInteger

const (
	NSModalResponseOK NSModalResponse = 1
)

const (
	NSBackingStoreBuffered NSBackingStoreType = 2
)

type WKUserScriptInjectionTime NSInteger

const (
	WKUserScriptInjectionTimeAtDocumentStart WKUserScriptInjectionTime = 0
)

func NSApplication_GetSharedApplication() NSApplication {
	return NSApplication{objc.ID(class_NSApplication).Send(sel_sharedApplication)}
}

func (app NSApplication) ActivateIgnoringOtherApps(flag bool) {
	app.Send(sel_activateIgnoringOtherApps, flag)
}

func (app NSApplication) SetDelegate(delegate objc.ID) {
	app.Send(sel_setDelegate, delegate)
}

func (app NSApplication) Run() {
	app.Send(sel_run)
}

func (app NSApplication) SetActivationPolicy(policy NSApplicationActivationPolicy) {
	app.Send(sel_setActivationPolicy, policy)
}

func (app NSApplication) Stop(sender objc.ID) {
	app.Send(sel_stop, sender)
}

func (app NSApplication) PostEvent(event NSEvent, atStart bool) {
	app.Send(sel_postEvent_atStart_, event, atStart)
}

type NSBundle struct {
	objc.ID
}

func NSBundle_GetMainBundle() NSBundle {
	return NSBundle{objc.ID(class_NSBundle).Send(sel_mainBundle)}
}

func (bundle NSBundle) BundlePath() NSString {
	return NSString{bundle.Send(sel_bundlePath)}
}

type NSEvent struct {
	objc.ID
}

// NSInvocation is being used to call functions that can't be called directly with purego.SyscallN.
// See the downsides of that function for what it cannot do.
type NSInvocation struct {
	objc.ID
}

func NSInvocation_invocationWithMethodSignature(sig NSMethodSignature) NSInvocation {
	return NSInvocation{objc.ID(class_NSInvocation).Send(sel_invocationWithMethodSignature, sig.ID)}
}

func (i NSInvocation) SetSelector(cmd objc.SEL) {
	i.Send(sel_setSelector, cmd)
}

func (i NSInvocation) SetTarget(target objc.ID) {
	i.Send(sel_setTarget, target)
}

func (i NSInvocation) SetArgumentAtIndex(arg unsafe.Pointer, idx int) {
	i.Send(sel_setArgumentAtIndex, arg, idx)
}

func (i NSInvocation) GetReturnValue(ret unsafe.Pointer) {
	i.Send(sel_getReturnValue, ret)
}

func (i NSInvocation) Invoke() {
	i.Send(sel_invoke)
}

func (i NSInvocation) InvokeWithTarget(target objc.ID) {
	i.Send(sel_invokeWithTarget, target)
}

type NSMethodSignature struct {
	objc.ID
}

func NSMethodSignature_instanceMethodSignatureForSelector(self objc.ID, cmd objc.SEL) NSMethodSignature {
	return NSMethodSignature{self.Send(sel_instanceMethodSignatureForSelector, cmd)}
}

// NSMethodSignature_signatureWithObjCTypes takes a string that represents the type signature of a method.
// It follows the encoding specified in the Apple Docs.
//
// [Apple Docs]: https://developer.apple.com/library/archive/documentation/Cocoa/Conceptual/ObjCRuntimeGuide/Articles/ocrtTypeEncodings.html#//apple_ref/doc/uid/TP40008048-CH100
func NSMethodSignature_signatureWithObjCTypes(types string) NSMethodSignature {
	return NSMethodSignature{objc.ID(class_NSMethodSignature).Send(sel_signatureWithObjCTypes, types)}
}

type NSString struct {
	objc.ID
}

func NSString_alloc() NSString {
	return NSString{objc.ID(class_NSString).Send(sel_alloc)}
}

func (s NSString) InitWithUTF8String(utf8 string) NSString {
	return NSString{ID: s.Send(sel_initWithUTF8String, utf8)}
}

func StringWithUTF8String(utf8 string) NSString {
	return NSString{objc.ID(class_NSString).Send(sel_stringWithUTF8String, utf8)}
}

func (s NSString) String() string {
	return string(unsafe.Slice((*byte)(unsafe.Pointer(s.Send(sel_UTF8String))), s.Send(sel_length)))
}

type NSWindow struct {
	objc.ID
}

func NSWindow_alloc() NSWindow {
	return NSWindow{objc.ID(class_NSWindow).Send(sel_alloc)}
}

func (window NSWindow) InitWithContentRect(rect NSRect, styleMask NSWindowStyleMask, backing NSBackingStoreType, deferCreation bool) NSWindow {
	// inv := NSInvocation_invocationWithMethodSignature(NSMethodSignature_signatureWithObjCTypes("v@:{CGRect={CGPoint=dd}{CGSize=dd}}QcB"))
	inv := NSInvocation_invocationWithMethodSignature(NSMethodSignature_signatureWithObjCTypes("{NSWindow=#}@:{CGRect={CGPoint=dd}{CGSize=dd}}QcB"))
	inv.SetTarget(window.ID)
	inv.SetSelector(sel_initWithContentRect_styleMask_backing_defer_)
	inv.SetArgumentAtIndex(unsafe.Pointer(&rect), 2)
	inv.SetArgumentAtIndex(unsafe.Pointer(&styleMask), 3)
	inv.SetArgumentAtIndex(unsafe.Pointer(&backing), 4)
	inv.SetArgumentAtIndex(unsafe.Pointer(&deferCreation), 5)
	inv.Invoke()

	var ret NSWindow
	inv.GetReturnValue(unsafe.Pointer(&ret))
	return ret
}

func (window NSWindow) SetContentView(view objc.ID) {
	window.Send(sel_setContentView, view)
}

func (window NSWindow) MakeKeyAndOrderFront(sender objc.ID) {
	window.Send(sel_makeKeyAndOrderFront, sender)
}

func (window NSWindow) SetStyleMask(style NSWindowStyleMask) {
	window.Send(sel_setStyleMask, style)
}

func (window NSWindow) SetFrame(frameRect NSRect, display bool, animate bool) {
	inv := NSInvocation_invocationWithMethodSignature(NSMethodSignature_signatureWithObjCTypes("{NSWindow=#}@:{CGRect={CGPoint=dd}{CGSize=dd}}BB"))
	inv.SetTarget(window.ID)
	inv.SetSelector(sel_setFrameDisplayAnimate)
	inv.SetArgumentAtIndex(unsafe.Pointer(&frameRect), 2)
	inv.SetArgumentAtIndex(unsafe.Pointer(&display), 3)
	inv.SetArgumentAtIndex(unsafe.Pointer(&animate), 4)
	inv.Invoke()
}

func (window NSWindow) SetContentMinSize(size NSSize) {
	inv := NSInvocation_invocationWithMethodSignature(NSMethodSignature_signatureWithObjCTypes("v@:{CGSize=dd}"))
	inv.SetTarget(window.ID)
	inv.SetSelector(sel_setContentMinSize)
	inv.SetArgumentAtIndex(unsafe.Pointer(&size), 2)
	inv.Invoke()
}

func (window NSWindow) SetContentMaxSize(size NSSize) {
	inv := NSInvocation_invocationWithMethodSignature(NSMethodSignature_signatureWithObjCTypes("v@:{CGSize=dd}"))
	inv.SetTarget(window.ID)
	inv.SetSelector(sel_setContentMaxSize)
	inv.SetArgumentAtIndex(unsafe.Pointer(&size), 2)
	inv.Invoke()
}

func (window NSWindow) Center() {
	window.Send(sel_center)
}

func (window NSWindow) SetTitle(title string) {
	wrappedTitle := NSString_alloc().InitWithUTF8String(title)
	window.Send(objc.RegisterName("setTitle:"), wrappedTitle.ID)
}

type NSOpenPanel struct {
	objc.ID
}

func NSOpenPanel_openPanel() NSOpenPanel {
	return NSOpenPanel{objc.ID(objc.GetClass("NSOpenPanel")).Send(objc.RegisterName("openPanel"))}
}

func (panel NSOpenPanel) SetCanChooseFiles(flag bool) {
	panel.Send(objc.RegisterName("setCanChooseFiles:"), flag)
}

func (panel NSOpenPanel) SetCanChooseDirectories(flag bool) {
	panel.Send(objc.RegisterName("setCanChooseDirectories:"), flag)
}

func (panel NSOpenPanel) SetAllowsMultipleSelection(flag bool) {
	panel.Send(objc.RegisterName("setAllowsMultipleSelection:"), flag)
}

func (panel NSOpenPanel) RunModal() NSModalResponse {
	return NSModalResponse(panel.Send(objc.RegisterName("runModal")))
}

func (panel NSOpenPanel) URLs() objc.ID {
	return panel.Send(objc.RegisterName("URLs"))
}

type NSNumber struct {
	objc.ID
}

func NSNumber_NumberWithBool(value bool) NSNumber {
	return NSNumber{objc.ID(objc.GetClass("NSNumber")).Send(objc.RegisterName("numberWithBool:"), value)}
}

type CGPoint struct {
	X float64
	Y float64
}
type NSPoint = CGPoint

type CGSize struct {
	Width  float64
	Height float64
}
type NSSize = CGSize

type CGRect struct {
	Origin CGPoint
	Size   CGSize
}
type NSRect = CGRect

type NSURL struct {
	objc.ID
}

func NSURL_URLWithString(url string) NSURL {
	wrappedUrl := NSString_alloc().InitWithUTF8String(url)
	return NSURL{objc.ID(class_NSURL).Send(objc.RegisterName("URLWithString:"), wrappedUrl.ID)}
}

type NSURLRequest struct {
	objc.ID
}

func NSURLRequest_requestWithURL(url NSURL) NSURLRequest {
	return NSURLRequest{objc.ID(class_NSURLRequest).Send(objc.RegisterName("requestWithURL:"), url.ID)}
}

type WebviewAppDelegate struct {
	objc.ID
}

// WebView

var (
	class_WKUserScript           objc.Class
	class_WKWebView              objc.Class
	class_WKWebViewConfiguration objc.Class

	sel_UserContentController objc.SEL
)

type WKWebViewConfiguration struct {
	objc.ID
}

func WKWebViewConfiguration_new() WKWebViewConfiguration {
	return WKWebViewConfiguration{objc.ID(class_WKWebViewConfiguration).Send(Sel_new)}
}

func (c WKWebViewConfiguration) UserContentController() WKUserContentController {
	return WKUserContentController{c.Send(sel_UserContentController)}
}

func (c WKWebViewConfiguration) Preferences() WKPreferences {
	return WKPreferences{c.Send(objc.RegisterName("preferences"))}
}

type WKPreferences struct {
	objc.ID
}

func (p WKPreferences) SetValue(key string, value objc.ID) {
	wrappedString := NSString_alloc().InitWithUTF8String(key)
	p.Send(objc.RegisterName("setValue:forKey:"), value, wrappedString.ID)
}

type WKUserContentController struct {
	objc.ID
}

func (c WKUserContentController) AddScriptMessageHandler(handler objc.ID, name string) {
	wrappedName := NSString_alloc().InitWithUTF8String(name)
	c.Send(objc.RegisterName("addScriptMessageHandler:name:"), handler, wrappedName.ID)
}

func (c WKUserContentController) AddUserScript(script WKUserScript) {
	c.Send(objc.RegisterName("addUserScript:"), script.ID)
}

type WKWebView struct {
	objc.ID
}

func WKWebView_alloc() WKWebView {
	return WKWebView{objc.ID(class_WKWebView).Send(sel_alloc)}
}

func (w WKWebView) InitWithFrame(frame NSRect, configuration WKWebViewConfiguration) WKWebView {
	inv := NSInvocation_invocationWithMethodSignature(NSMethodSignature_signatureWithObjCTypes("{WKWebView=#}@:{CGRect={CGPoint=dd}{CGSize=dd}}@"))
	inv.SetTarget(w.ID)
	inv.SetSelector(objc.RegisterName("initWithFrame:configuration:"))
	inv.SetArgumentAtIndex(unsafe.Pointer(&frame), 2)
	inv.SetArgumentAtIndex(unsafe.Pointer(&configuration), 3)
	inv.Invoke()

	var ret WKWebView
	inv.GetReturnValue(unsafe.Pointer(&ret))
	return ret
}

func (w WKWebView) SetUIDelegate(delegate objc.ID) {
	w.Send(objc.RegisterName("setUIDelegate:"), delegate)
}

func (w WKWebView) LoadRequest(request NSURLRequest) {
	w.Send(objc.RegisterName("loadRequest:"), request.ID)
}

func (w WKWebView) LoadHTMLString(html string, baseURL objc.ID) {
	wrappedHtml := StringWithUTF8String(html)
	w.Send(sel_loadHTMLStringBaseURL, wrappedHtml.ID)
}

func (w WKWebView) EvaluateJavaScript(js string, completionHandler objc.ID) {
	wrappedJs := StringWithUTF8String(js)
	w.Send(objc.RegisterName("evaluateJavaScript:completionHandler:"), wrappedJs.ID, completionHandler)
}

type WKUserScript struct {
	objc.ID
}

func WKUserScript_alloc() WKUserScript {
	return WKUserScript{
		objc.ID(class_WKUserScript).Send(sel_alloc),
	}
}

func (s WKUserScript) InitWithSource(source string, injectionTime WKUserScriptInjectionTime, forMainFrameOnly bool) WKUserScript {
	sourceWrapped := NSString_alloc().InitWithUTF8String(source)

	inv := NSInvocation_invocationWithMethodSignature(NSMethodSignature_signatureWithObjCTypes("{WKUserScript=#}@:@QcB"))
	inv.SetTarget(s.ID)
	inv.SetSelector(objc.RegisterName("initWithSource:injectionTime:forMainFrameOnly:"))
	inv.SetArgumentAtIndex(unsafe.Pointer(&sourceWrapped), 2)
	inv.SetArgumentAtIndex(unsafe.Pointer(&injectionTime), 3)
	inv.SetArgumentAtIndex(unsafe.Pointer(&forMainFrameOnly), 4)
	inv.Invoke()

	var ret WKUserScript
	inv.GetReturnValue(unsafe.Pointer(&ret))
	return ret
}

type NSAutoreleasePool struct {
	objc.ID
}

func NSAutoreleasePool_new() NSAutoreleasePool {
	return NSAutoreleasePool{objc.ID(class_NSAutoreleasePool).Send(Sel_new)}
}

func (p NSAutoreleasePool) Release() {
	p.Send(sel_release)
}
