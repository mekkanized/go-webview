//go:build !windows

package webview

import "unsafe"

// This is copied from webview/webview.
// The documentation is included for convenience.

// Hint is used to configure window sizing and resizing behavior.
type Hint int

type WebViewOptions struct {
	Window unsafe.Pointer
	Debug  bool

	DisableHardwareAcceleration bool
}
