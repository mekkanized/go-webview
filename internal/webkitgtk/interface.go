package webkitgtk

type (
	GAsyncResult uintptr
	GCancellable uintptr
	GObject      uintptr
	GtkContainer uintptr
	GtkWidget    uintptr
	GtkWindow    uintptr

	GAsyncReadyCallback func(sourceObject GObject, res GAsyncResult, userData uintptr)
	GDestroyNotify      func(data uintptr)
	GSourceFunc         func(userData uintptr) bool
	GCallback           interface{}
	GClosureNotify      func(data uintptr, closure uintptr)

	JSCValue                 uintptr
	JSContextRef             uintptr
	JSValueRef               uintptr
	WebKitJavascriptResult   uintptr
	WebKitSettings           uintptr
	WebKitUserContentManager uintptr
	WebKitUserScript         uintptr
	WebKitWebView            uintptr
)

const (
	G_SOURCE_CONTINUE       = true
	G_SOURCE_REMOVE         = false
	G_PRIORITY_HIGH         = -100
	G_PRIORITY_DEFAULT      = 0
	G_PRIORITY_HIGH_IDLE    = 100
	G_PRIORITY_DEFAULT_IDLE = 200
	G_PRIORITY_LOW          = 300

	NULLPTR uintptr = 0
)

type GdkGravity uint

const (
	GDK_GRAVITY_NORTH_WEST GdkGravity = 1 << iota
	GDK_GRAVITY_NORTH
	GDK_GRAVITY_NORTH_EAST
	GDK_GRAVITY_WEST
	GDK_GRAVITY_CENTER
	GDK_GRAVITY_EAST
	GDK_GRAVITY_SOUTH_WEST
	GDK_GRAVITY_SOUTH
	GDK_GRAVITY_SOUTH_EAST
	GDK_GRAVITY_STATIC
)

type GdkGeometry struct {
	MinWidth   int32
	MinHeight  int32
	MaxWidth   int32
	MaxHeight  int32
	BaseWidth  int32
	BaseHeight int32
	WidthInc   int32
	HeightInc  int32
	MinAspect  float64
	MaxAspect  float64
	WinGravity GdkGravity
}

type GdkWindowHints uint

const (
	GDK_HINT_POS         GdkWindowHints = 1 << 0
	GDK_HINT_MIN_SIZE    GdkWindowHints = 1 << 1
	GDK_HINT_MAX_SIZE    GdkWindowHints = 1 << 2
	GDK_HINT_BASE_SIZE   GdkWindowHints = 1 << 3
	GDK_HINT_ASPECT      GdkWindowHints = 1 << 4
	GDK_HINT_RESIZE_INC  GdkWindowHints = 1 << 5
	GDK_HINT_WIN_GRAVITY GdkWindowHints = 1 << 6
	GDK_HINT_USER_POS    GdkWindowHints = 1 << 7
	GDK_HINT_USER_SIZE   GdkWindowHints = 1 << 8
)

type GtkWindowType uint

const (
	GTK_WINDOW_TOPLEVEL GtkWindowType = iota
	GTK_WINDOW_POPUP
)

type GConnectFlags uint

const (
	G_CONNECT_DEFAULT GConnectFlags = iota
	G_CONNECT_AFTER
	G_CONNECT_SWAPPED
)

type WebKitHardwareAccelerationPolicy uint

const (
	WEBKIT_HARDWARE_ACCELERATION_POLICY_ON_DEMAND WebKitHardwareAccelerationPolicy = iota
	WEBKIT_HARDWARE_ACCELERATION_POLICY_ALWAYS
	WEBKIT_HARDWARE_ACCELERATION_POLICY_NEVER
)

type WebKitUserContentInjectedFrames uint

const (
	WEBKIT_USER_CONTENT_INJECT_ALL_FRAMES WebKitUserContentInjectedFrames = iota
	WEBKIT_USER_CONTENT_INJECT_TOP_FRAME
)

type WebKitUserScriptInjectionTime uint

const (
	WEBKIT_USER_SCRIPT_INJECT_AT_DOCUMENT_START WebKitUserScriptInjectionTime = iota
	WEBKIT_USER_SCRIPT_INJECT_AT_DOCUMENT_END
)

type Context interface {
	LoadFunctions() error

	// GLib
	GFree(mem uintptr)
	GIdleAddFull(priority int, function GSourceFunc, data uintptr, notify GDestroyNotify)
	GSignalConnectData(instance GtkWidget, detailedSignal string, cHandler GCallback, data uintptr, destroyData GClosureNotify, connectFlags GConnectFlags) uint32
	GtkContainerAdd(container GtkContainer, widget GtkWidget)
	GtkInitCheck() bool
	GtkMain()
	GtkMainQuit()
	GtkWidgetGrabFocus(widget GtkWidget)
	GtkWidgetSetSizeRequest(widget GtkWidget, width, height int)
	GtkWidgetShowAll(widget GtkWidget)
	GtkWindowNew(windowType GtkWindowType) GtkWidget
	GtkWindowResize(window GtkWindow, width, height int)
	GtkWindowSetGeometryHints(window GtkWindow, geometryWidget GtkWidget, geometry GdkGeometry, geomMask GdkWindowHints)
	GtkWindowSetResizable(window GtkWindow, resizable bool)
	GtkWindowSetTitle(window GtkWindow, title string)

	// WebKit
	JsCValueToString(value JSCValue) string
	WebKitGetMajorVersion() uint32
	WebKitGetMinorVersion() uint32
	WebKitGetMicroVersion() uint32
	WebKitWebViewNew() GtkWidget
	WebKitWebViewGetUserContentManager(webview WebKitWebView) WebKitUserContentManager
	WebKitWebViewGetSettings(webview WebKitWebView) WebKitSettings
	WebKitWebViewLoadURI(webview WebKitWebView, uri string)
	WebKitWebViewLoadHTML(webview WebKitWebView, content string, baseUri string)
	WebKitWebViewRunJavascript(webview WebKitWebView, script string, cancellable GCancellable, callback GAsyncReadyCallback, userData uintptr)
	WebKitJavascriptResultGetJsValue(jsResult WebKitJavascriptResult) JSCValue
	WebKitUserContentManagerAddScript(manager WebKitUserContentManager, script WebKitUserScript)
	WebKitUserContentManagerRegisterScriptMessageHandler(manager WebKitUserContentManager, name string)
	WebKitUserScriptNew(source string, injectedFrames WebKitUserContentInjectedFrames, injectionTime WebKitUserScriptInjectionTime, whitelist string, blacklist string) WebKitUserScript
	WebKitSettingsSetEnableDeveloperExtras(settings WebKitSettings, enabled bool)
	WebKitSettingsSetEnableWriteConsoleMessagesToStdout(settings WebKitSettings, enabled bool)
	WebKitSettingsSetJavascriptCanAccessClipboard(settings WebKitSettings, enabled bool)
}
