module github.com/mekkanized/go-webview

go 1.18

replace github.com/progrium/macdriver => ../darwinkit

require (
	github.com/ebitengine/purego v0.8.2
	github.com/jchv/go-webview2 v0.0.0-20221223143126-dc24628cff85
)

require (
	github.com/jchv/go-winloader v0.0.0-20200815041850-dec1ee9a7fd5 // indirect
	golang.org/x/sys v0.9.0 // indirect
)
