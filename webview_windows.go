//go:build windows

package webview

import (
	"github.com/jchv/go-webview2"
)

// New calls NewWindow to create a new window and a new webview instance. If debug
// is non-zero - developer tools will be enabled (if the platform supports them).
func New(debug bool) WebView {
	return webview2.New(debug)
}

// NewWithOptions creates a new webview using the provided options.
func NewWithOptions(options WebViewOptions) WebView {
	winOptions := webview2.WebViewOptions{
		Window:    options.Window,
		Debug:     options.Debug,
		DataPath:  options.DataPath,
		AutoFocus: options.AutoFocus,
		WindowOptions: webview2.WindowOptions{
			Title:  options.WindowOptions.Title,
			Width:  options.WindowOptions.Width,
			Height: options.WindowOptions.Height,
			IconId: options.WindowOptions.IconId,
			Center: options.WindowOptions.Center,
		},
	}
	return webview2.NewWithOptions(winOptions)
}
