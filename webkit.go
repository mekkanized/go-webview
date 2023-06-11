package webkit

import (
	"fmt"
	"unsafe"

	"github.com/Mekkabotics/go-webkit/webkitloader"
)

type webkit struct {
  webview unsafe.Pointer
  window unsafe.Pointer
}

type WebViewOptions struct {
  Debug bool
}

// NewWithOptions creates a new webview using the provided options.
func NewWithOptions(options WebViewOptions) WebView {
  w := &webkit{}

  test := webkitloader.GTKInitCheck()
  fmt.Printf("Loaded: %t\n", test)
  if !test {
    return nil
  }

  if w.window == nil {
    w.window = webkitloader.GTKWindowNew(webkitloader.GTK_WINDOW_TOPLEVEL)
  }

  // TODO: https://github.com/webview/webview/blob/899018ad0e5cc22a18cd734393ccae4d55e3b2b4/webview.h#L525
  // destroyCb := purego.NewCallback(func(unsafe.Pointer, unsafe.Pointer) {
  //   webkitloader.GTKMainQuit()
  // })
  // webkitloader.GTKSignalConnectData(w.window, webkitloader.CString("destroy"), destroyCb, nil, nil, 0)

  // Initialize webview widget
  w.webview = webkitloader.WebkitWebViewNew()
  // manager := webkitloader.WebkitWebViewGetUserContentManager(w.webview) 
  // TODO: https://github.com/webview/webview/blob/899018ad0e5cc22a18cd734393ccae4d55e3b2b4/webview.h#L534

  webkitloader.GTKContainerAdd(w.window, w.webview)
  webkitloader.GTKWidgetGrabFocus(w.webview)

  settings := webkitloader.WebkitWebViewGetSettings(w.webview)
  webkitloader.WebkitSettingsSetJavascriptCanAccessClipboard(settings, true)
  if options.Debug {
    webkitloader.WebkitSettingsSetEnableWriteConsoleMessagesToStdout(settings, true)
    webkitloader.WebkitSettingsSetEnableDeveloperExtras(settings, true)
  }

  webkitloader.GTKWidgetShowAll(w.window)

  return w
}

func (w *webkit) Run() {
  webkitloader.GTKMain()
}

func (w *webkit) Terminate() {

}

func (w *webkit) Dispatch(f func()) {

}

func (w *webkit) Destroy() {

}

func (w *webkit) Window() unsafe.Pointer {
  return nil
}

func (w *webkit) SetTitle(title string) {
  webkitloader.GTKWindowSetTitle(w.window, webkitloader.CString(title))
}

func (w *webkit) SetSize(width int, height int, hint Hint) {
  webkitloader.GTKWindowSetResizable(w.window, hint != HintFixed)
  switch(hint) {
  case HintNone:
    webkitloader.GTKWindowResize(w.window, width, height)
  // case HintFixed:
  // default:
  }
}

func (w *webkit) Navigate(url string) {
  webkitloader.WebkitWebViewLoadURI(w.webview, webkitloader.CString(url))
}

func (w *webkit) SetHtml(html string) {
  webkitloader.WebkitWebViewLoadHTML(w.webview, webkitloader.CString(html), nil)
}

func (w *webkit) Init(js string) {

}

func (w *webkit) Eval(js string) {

}

func (w *webkit) Bind(name string, f interface{}) error {
  return nil
}

