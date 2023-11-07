package webview

import (
	"unsafe"
)

type webview struct {
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
	w := &webview{}

	return w
}

func (w *webview) Run() {
}

func (w *webview) Terminate() {
}

func (w *webview) Dispatch(f func()) {
}

func (w *webview) Destroy() {
}

func (w *webview) Window() unsafe.Pointer {
	return nil
}

func (w *webview) SetTitle(title string) {
}

func (w *webview) SetSize(width int, height int, hint Hint) {
}

func (w *webview) Navigate(url string) {
}

func (w *webview) SetHtml(html string) {
}

func (w *webview) Init(js string) {
}

func (w *webview) Eval(js string) {
}

func (w *webview) Bind(name string, f interface{}) error {
	return nil
}
