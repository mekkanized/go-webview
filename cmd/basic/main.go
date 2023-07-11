package main

import (
	"log"

	"github.com/Mekkanized/go-webview"
)

func main() {
	w := webview.NewWithOptions(webview.WebViewOptions{
		Debug: true,
	})
	if w == nil {
		log.Fatalln("Failed to load webview.")
	}
	defer w.Destroy()

	w.SetTitle("Basic Example")
	w.SetSize(480, 320, webview.HintNone)
	w.SetHtml("Thanks for using webview!")

	w.Run()
}
