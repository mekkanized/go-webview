//go:build linux

package webkitgtk

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"
)

type defaultContext struct {
	// GTK
	gFree                     uintptr
	gIdleAddFull              uintptr
	gSignalConnectData        uintptr
	gtkContainerAdd           uintptr
	gtkInitCheck              uintptr
	gtkMain                   uintptr
	gtkMainQuit               uintptr
	gtkWidgetGrabFocus        uintptr
	gtkWidgetSetSizeRequest   uintptr
	gtkWidgetShowAll          uintptr
	gtkWindowNew              uintptr
	gtkWindowResize           uintptr
	gtkWindowSetGeometryHints uintptr
	gtkWindowSetResizable     uintptr
	gtkWindowSetTitle         uintptr

	// WebKit
	jsCValueToString                                     uintptr
	webKitGetMajorVersion                                uintptr
	webKitGetMinorVersion                                uintptr
	webKitGetMicroVersion                                uintptr
	webKitWebViewNew                                     uintptr
	webKitWebViewGetSettings                             uintptr
	webKitWebViewGetUserContentManager                   uintptr
	webKitWebViewLoadHTML                                uintptr
	webKitWebViewLoadURI                                 uintptr
	webKitWebViewRunJavascript                           uintptr
	webKitJavascriptResultGetJsValue                     uintptr
	webKitUserContentManagerAddScript                    uintptr
	webKitUserContentManagerRegisterScriptMessageHandler uintptr
	webKitUserScriptNew                                  uintptr
	webKitSettingsSetEnableDeveloperExtras               uintptr
	webKitSettingsSetEnableWriteConsoleMessagesToStdout  uintptr
	webKitSettingsSetJavascriptCanAccessClipboard        uintptr
}

func NewDefaultContext() (Context, error) {
	ctx := &defaultContext{}
	if err := ctx.init(); err != nil {
		return nil, fmt.Errorf("failed to initialize default context: %w", err)
	}

	return ctx, nil
}

func (c *defaultContext) GFree(mem uintptr) {
	purego.SyscallN(c.gFree, mem)
}

func (c *defaultContext) GIdleAddFull(priority int, function GSourceFunc, data uintptr, notify GDestroyNotify) {
	var destroyCb uintptr = NULLPTR
	if notify != nil {
		destroyCb = purego.NewCallback(notify)
	}

	purego.SyscallN(c.gIdleAddFull, uintptr(priority), purego.NewCallback(function), data, destroyCb)
}

func (c *defaultContext) GSignalConnectData(instance GtkWidget, detailedSignal string, cHandler GCallback, data uintptr, destroyData GClosureNotify, connectFlags GConnectFlags) uint32 {
	cstrDetailedSignal, free := cStr(detailedSignal)
	defer free()

	var destroyCb uintptr = NULLPTR
	if destroyData != nil {
		destroyCb = purego.NewCallback(destroyData)
	}

	ret, _, _ := purego.SyscallN(c.gSignalConnectData, uintptr(instance), uintptr(unsafe.Pointer(cstrDetailedSignal)), purego.NewCallback(cHandler), data, destroyCb, uintptr(connectFlags))
	return uint32(ret)
}

func (c *defaultContext) GtkContainerAdd(container GtkContainer, widget GtkWidget) {
	purego.SyscallN(c.gtkContainerAdd, uintptr(container), uintptr(widget))
}

func (c *defaultContext) GtkInitCheck() bool {
	ret, _, _ := purego.SyscallN(c.gtkInitCheck)
	return byte(ret) != 0
}

func (c *defaultContext) GtkMain() {
	purego.SyscallN(c.gtkMain)
}

func (c *defaultContext) GtkMainQuit() {
	purego.SyscallN(c.gtkMainQuit)
}

func (c *defaultContext) GtkWidgetGrabFocus(widget GtkWidget) {
	purego.SyscallN(c.gtkWidgetGrabFocus, uintptr(widget))
}

func (c *defaultContext) GtkWidgetSetSizeRequest(widget GtkWidget, width, height int) {
	purego.SyscallN(c.gtkWidgetSetSizeRequest, uintptr(widget), uintptr(width), uintptr(height))
}

func (c *defaultContext) GtkWidgetShowAll(widget GtkWidget) {
	purego.SyscallN(c.gtkWidgetShowAll, uintptr(widget))
}

func (c *defaultContext) GtkWindowNew(windowType GtkWindowType) GtkWidget {
	ret, _, _ := purego.SyscallN(c.gtkWindowNew, uintptr(windowType))
	return GtkWidget(ret)
}

func (c *defaultContext) GtkWindowResize(window GtkWindow, width, height int) {
	purego.SyscallN(c.gtkWindowResize, uintptr(window), uintptr(width), uintptr(height))
}

func (c *defaultContext) GtkWindowSetGeometryHints(window GtkWindow, geometryWidget GtkWidget, geometry GdkGeometry, geomMask GdkWindowHints) {
	purego.SyscallN(c.gtkWindowSetGeometryHints, uintptr(window), uintptr(geometryWidget), uintptr(unsafe.Pointer(&geometry)), uintptr(geomMask))
}

func (c *defaultContext) GtkWindowSetResizable(window GtkWindow, resizable bool) {
	purego.SyscallN(c.gtkWindowSetResizable, uintptr(window), uintptr(boolToInt(resizable)))
}

func (c *defaultContext) GtkWindowSetTitle(window GtkWindow, title string) {
	cstrTitle, free := cStr(title)
	defer free()
	purego.SyscallN(c.gtkWindowSetTitle, uintptr(window), uintptr(unsafe.Pointer(cstrTitle)))
}

// WebKit
func (c *defaultContext) JsCValueToString(value JSCValue) string {
	ret, _, _ := purego.SyscallN(c.jsCValueToString, uintptr(value))
	str := goStr(ret)
	c.GFree(ret)
	return str
}

func (c *defaultContext) WebKitGetMajorVersion() uint32 {
	ret, _, _ := purego.SyscallN(c.webKitGetMajorVersion)
	return uint32(ret)
}

func (c *defaultContext) WebKitGetMinorVersion() uint32 {
	ret, _, _ := purego.SyscallN(c.webKitGetMinorVersion)
	return uint32(ret)
}

func (c *defaultContext) WebKitGetMicroVersion() uint32 {
	ret, _, _ := purego.SyscallN(c.webKitGetMicroVersion)
	return uint32(ret)
}

func (c *defaultContext) WebKitWebViewNew() GtkWidget {
	ret, _, _ := purego.SyscallN(c.webKitWebViewNew)
	return GtkWidget(ret)
}

func (c *defaultContext) WebKitWebViewGetUserContentManager(webview WebKitWebView) WebKitUserContentManager {
	ret, _, _ := purego.SyscallN(c.webKitWebViewGetUserContentManager, uintptr(webview))
	return WebKitUserContentManager(ret)
}

func (c *defaultContext) WebKitWebViewGetSettings(webview WebKitWebView) WebKitSettings {
	ret, _, _ := purego.SyscallN(c.webKitWebViewGetSettings, uintptr(webview))
	return WebKitSettings(ret)
}

func (c *defaultContext) WebKitWebViewLoadURI(webview WebKitWebView, uri string) {
	cstrUri, free := cStr(uri)
	defer free()
	purego.SyscallN(c.webKitWebViewLoadURI, uintptr(webview), uintptr(unsafe.Pointer(cstrUri)))
}

func (c *defaultContext) WebKitWebViewLoadHTML(webview WebKitWebView, content string, baseUri string) {
	cstrContent, free := cStr(content)
	defer free()
	baseUriPtr := NULLPTR
	if baseUri != "" {
		cstrBaseUri, free := cStr(baseUri)
		defer free()
		baseUriPtr = uintptr(unsafe.Pointer(cstrBaseUri))
	}
	purego.SyscallN(c.webKitWebViewLoadHTML, uintptr(webview), uintptr(unsafe.Pointer(cstrContent)), baseUriPtr)
}

func (c *defaultContext) WebKitWebViewRunJavascript(webview WebKitWebView, script string, cancellable GCancellable, callback GAsyncReadyCallback, userData uintptr) {
	cstrScript, free := cStr(script)
	defer free()

	var callbackCb uintptr = NULLPTR
	if callback != nil {
		callbackCb = purego.NewCallback(callback)
	}

	purego.SyscallN(c.webKitWebViewRunJavascript, uintptr(webview), uintptr(unsafe.Pointer(cstrScript)), uintptr(cancellable), callbackCb, userData)
}

func (c *defaultContext) WebKitJavascriptResultGetJsValue(jsResult WebKitJavascriptResult) JSCValue {
	ret, _, _ := purego.SyscallN(c.webKitJavascriptResultGetJsValue, uintptr(jsResult))
	return JSCValue(ret)
}

func (c *defaultContext) WebKitUserContentManagerAddScript(manager WebKitUserContentManager, script WebKitUserScript) {
	purego.SyscallN(c.webKitUserContentManagerAddScript, uintptr(manager), uintptr(script))
}

func (c *defaultContext) WebKitUserContentManagerRegisterScriptMessageHandler(manager WebKitUserContentManager, name string) {
	cstrName, free := cStr(name)
	defer free()
	purego.SyscallN(c.webKitUserContentManagerRegisterScriptMessageHandler, uintptr(manager), uintptr(unsafe.Pointer(cstrName)))
}

func (c *defaultContext) WebKitUserScriptNew(source string, injectedFrames WebKitUserContentInjectedFrames, injectionTime WebKitUserScriptInjectionTime, whitelist string, blacklist string) WebKitUserScript {
	cstrSource, free := cStr(source)
	defer free()
	cstrWhitelist, free := cStr(whitelist)
	defer free()
	whitelistPtr := NULLPTR
	if whitelist != "" {
		whitelistPtr = uintptr(unsafe.Pointer(cstrWhitelist))
	}
	cstrBlacklist, free := cStr(blacklist)
	defer free()
	blacklistPtr := NULLPTR
	if blacklist != "" {
		blacklistPtr = uintptr(unsafe.Pointer(cstrBlacklist))
	}
	ret, _, _ := purego.SyscallN(c.webKitUserScriptNew, uintptr(unsafe.Pointer(cstrSource)), uintptr(injectedFrames), uintptr(injectionTime), whitelistPtr, blacklistPtr)
	return WebKitUserScript(ret)
}

func (c *defaultContext) WebKitSettingsSetEnableDeveloperExtras(settings WebKitSettings, enabled bool) {
	purego.SyscallN(c.webKitSettingsSetEnableDeveloperExtras, uintptr(settings), uintptr(boolToInt(enabled)))
}

func (c *defaultContext) WebKitSettingsSetEnableWriteConsoleMessagesToStdout(settings WebKitSettings, enabled bool) {
	purego.SyscallN(c.webKitSettingsSetEnableWriteConsoleMessagesToStdout, uintptr(settings), uintptr(boolToInt(enabled)))
}

func (c *defaultContext) WebKitSettingsSetJavascriptCanAccessClipboard(settings WebKitSettings, enabled bool) {
	purego.SyscallN(c.webKitSettingsSetJavascriptCanAccessClipboard, uintptr(settings), uintptr(boolToInt(enabled)))
}

func (c *defaultContext) LoadFunctions() error {
	g := &procAddressGetter{ctx: c}

	// GTK
	c.gFree = g.get("g_free")
	c.gIdleAddFull = g.get("g_idle_add_full")
	c.gSignalConnectData = g.get("g_signal_connect_data")
	c.gtkContainerAdd = g.get("gtk_container_add")
	c.gtkInitCheck = g.get("gtk_init_check")
	c.gtkMain = g.get("gtk_main")
	c.gtkMainQuit = g.get("gtk_main_quit")
	c.gtkWidgetGrabFocus = g.get("gtk_widget_grab_focus")
	c.gtkWidgetSetSizeRequest = g.get("gtk_widget_set_size_request")
	c.gtkWidgetShowAll = g.get("gtk_widget_show_all")
	c.gtkWindowNew = g.get("gtk_window_new")
	c.gtkWindowResize = g.get("gtk_window_resize")
	c.gtkWindowSetGeometryHints = g.get("gtk_window_set_geometry_hints")
	c.gtkWindowSetResizable = g.get("gtk_window_set_resizable")
	c.gtkWindowSetTitle = g.get("gtk_window_set_title")

	// WebKit
	c.jsCValueToString = g.get("jsc_value_to_string")
	c.webKitGetMajorVersion = g.get("webkit_get_major_version")
	c.webKitGetMinorVersion = g.get("webkit_get_minor_version")
	c.webKitGetMicroVersion = g.get("webkit_get_micro_version")
	c.webKitWebViewNew = g.get("webkit_web_view_new")
	c.webKitWebViewGetUserContentManager = g.get("webkit_web_view_get_user_content_manager")
	c.webKitWebViewGetSettings = g.get("webkit_web_view_get_settings")
	c.webKitWebViewLoadURI = g.get("webkit_web_view_load_uri")
	c.webKitWebViewLoadHTML = g.get("webkit_web_view_load_html")
	c.webKitWebViewRunJavascript = g.get("webkit_web_view_run_javascript")
	c.webKitJavascriptResultGetJsValue = g.get("webkit_javascript_result_get_js_value")
	c.webKitUserContentManagerAddScript = g.get("webkit_user_content_manager_add_script")
	c.webKitUserContentManagerRegisterScriptMessageHandler = g.get("webkit_user_content_manager_register_script_message_handler")
	c.webKitUserScriptNew = g.get("webkit_user_script_new")
	c.webKitSettingsSetEnableDeveloperExtras = g.get("webkit_settings_set_enable_developer_extras")
	c.webKitSettingsSetEnableWriteConsoleMessagesToStdout = g.get("webkit_settings_set_enable_write_console_messages_to_stdout")
	c.webKitSettingsSetJavascriptCanAccessClipboard = g.get("webkit_settings_set_javascript_can_access_clipboard")

	if g.err != nil {
		return fmt.Errorf("failed to load functions: %w", g.err)
	}

	return nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// cStr takes a Go string (with or without null-termination)
// and returns the C counterpart.
//
// The returned free function must be called once you are done using the string
// in order to free the memory.
func cStr(str string) (cstr *byte, free func()) {
	bs := []byte(str)
	if len(bs) == 0 || bs[len(bs)-1] != 0 {
		bs = append(bs, 0)
	}
	return &bs[0], func() {
		runtime.KeepAlive(bs)
		bs = nil
	}
}

// goStr copies a char* to a Go string.
func goStr(c uintptr) string {
	// We take the address and then dereference it to trick go vet from creating a possible misuse of unsafe.Pointer
	ptr := *(*unsafe.Pointer)(unsafe.Pointer(&c))
	if ptr == nil {
		return ""
	}
	var length int
	for {
		if *(*byte)(unsafe.Add(ptr, uintptr(length))) == '\x00' {
			break
		}
		length++
	}
	return string(unsafe.Slice((*byte)(ptr), length))
}
