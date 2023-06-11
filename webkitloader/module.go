package webkitloader

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"
)

type GTKWindowType int
const (
  GTK_WINDOW_TOPLEVEL GTKWindowType = iota
  GTK_WINDOW_POPUP
)

type GTKWidget struct {}

var (
  GTKContainerAdd func(unsafe.Pointer, unsafe.Pointer)
  GTKInitCheck func() bool // TODO: Add argc/argv as arguments
  GTKMain func()
  GTKMainQuit func()
  GTKSignalConnectData func(unsafe.Pointer, *byte, unsafe.Pointer, unsafe.Pointer, unsafe.Pointer, int) // TODO: Working cb functions
  GTKWidgetGrabFocus func(unsafe.Pointer)
  GTKWidgetShowAll func(unsafe.Pointer)
  GTKWindowNew func(GTKWindowType) unsafe.Pointer
  GTKWindowResize func(unsafe.Pointer, int, int)
  GTKWindowSetResizable func(unsafe.Pointer, bool)
  GTKWindowSetTitle func(unsafe.Pointer, *byte)

  WebkitGetMajorVersion func() int
  WebkitGetMinorVersion func() int
  WebkitGetMicroVersion func() int

  WebkitWebViewNew func() unsafe.Pointer
  WebkitWebViewGetSettings func(unsafe.Pointer) unsafe.Pointer
  WebkitWebViewGetUserContentManager func(unsafe.Pointer) unsafe.Pointer
  WebkitWebViewLoadHTML func(unsafe.Pointer, *byte, unsafe.Pointer)
  WebkitWebViewLoadURI func(unsafe.Pointer, *byte)

  WebkitSettingsSetEnableDeveloperExtras func(unsafe.Pointer, bool)
  WebkitSettingsSetEnableWriteConsoleMessagesToStdout func(unsafe.Pointer, bool)
  WebkitSettingsSetJavascriptCanAccessClipboard func(unsafe.Pointer, bool)
)

func getSystemLibrary() string {
	switch runtime.GOOS {
	case "darwin":
		return "/usr/lib/libSystem.B.dylib"
	case "linux":
    return "libwebkit2gtk-4.0.so"
	default:
		panic(fmt.Errorf("GOOS=%s is not supported", runtime.GOOS))
	}
}

func init() {
	libwebkit, err := purego.Dlopen(getSystemLibrary(), purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		panic(err)
	}

  purego.RegisterLibFunc(&WebkitGetMajorVersion, libwebkit, "webkit_get_major_version")
  purego.RegisterLibFunc(&WebkitGetMinorVersion, libwebkit, "webkit_get_minor_version")
  purego.RegisterLibFunc(&WebkitGetMicroVersion, libwebkit, "webkit_get_micro_version")
  fmt.Printf("Running WebKit v%d.%d.%d\n", WebkitGetMajorVersion(), WebkitGetMinorVersion(), WebkitGetMicroVersion())

  purego.RegisterLibFunc(&GTKContainerAdd, libwebkit, "gtk_container_add")
  purego.RegisterLibFunc(&GTKInitCheck, libwebkit, "gtk_init_check")
  purego.RegisterLibFunc(&GTKWindowNew, libwebkit, "gtk_window_new")
  purego.RegisterLibFunc(&GTKWidgetShowAll, libwebkit, "gtk_widget_show_all")
  purego.RegisterLibFunc(&GTKMain, libwebkit, "gtk_main")
  purego.RegisterLibFunc(&GTKMainQuit, libwebkit, "gtk_main_quit")
  purego.RegisterLibFunc(&GTKWindowResize, libwebkit, "gtk_window_resize")
  purego.RegisterLibFunc(&GTKWidgetGrabFocus, libwebkit, "gtk_widget_grab_focus")
  purego.RegisterLibFunc(&GTKWindowSetResizable, libwebkit, "gtk_window_set_resizable")
  purego.RegisterLibFunc(&GTKWindowSetTitle, libwebkit, "gtk_window_set_title")
  purego.RegisterLibFunc(&GTKSignalConnectData, libwebkit, "g_signal_connect_data")

  purego.RegisterLibFunc(&WebkitWebViewNew, libwebkit, "webkit_web_view_new")
  purego.RegisterLibFunc(&WebkitWebViewGetUserContentManager, libwebkit, "webkit_web_view_get_user_content_manager")
  purego.RegisterLibFunc(&WebkitWebViewGetSettings, libwebkit, "webkit_web_view_get_settings")
  purego.RegisterLibFunc(&WebkitWebViewLoadURI, libwebkit, "webkit_web_view_load_uri")
  purego.RegisterLibFunc(&WebkitWebViewLoadHTML, libwebkit, "webkit_web_view_load_html")

  purego.RegisterLibFunc(&WebkitSettingsSetEnableDeveloperExtras, libwebkit, "webkit_settings_set_enable_developer_extras")
  purego.RegisterLibFunc(&WebkitSettingsSetEnableWriteConsoleMessagesToStdout, libwebkit, "webkit_settings_set_enable_write_console_messages_to_stdout")
  purego.RegisterLibFunc(&WebkitSettingsSetJavascriptCanAccessClipboard, libwebkit, "webkit_settings_set_javascript_can_access_clipboard")
}

// hasSuffix tests whether the string s ends with suffix.
func hasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

// CString converts a go string to *byte that can be passed to C code.
func CString(name string) *byte {
	if hasSuffix(name, "\x00") {
		return &(*(*[]byte)(unsafe.Pointer(&name)))[0]
	}
	var b = make([]byte, len(name)+1)
	copy(b, name)
	return &b[0]
}

