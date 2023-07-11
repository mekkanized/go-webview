package webview

import (
	"unsafe"

	"github.com/jchv/go-webview2"
)

func New(debug bool) WebView {
	return webview2.New(debug)
}

func NewWindow(debug bool, window unsafe.Pointer) WebView {
	return webview2.NewWindow(debug, window)
}

func NewWithOptions(options WebViewOptions) WebView {
	return webview2.NewWithOptions(options)
}
