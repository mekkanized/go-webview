//go:build windows

package webview

import (
	"unsafe"

	"github.com/jchv/go-webview2"
)

// New calls NewWindow to create a new window and a new webview instance. If debug
// is non-zero - developer tools will be enabled (if the platform supports them).
func New(debug bool) WebView {
	return webview2.New(debug)
}

// NewWithOptions creates a new webview using the provided options.
func NewWithOptions(options WebViewOptions) WebView {
	return webview2.NewWithOptions(options)
}
