package main

import (
	"log"

	"github.com/mekkanized/go-webview"
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
	w.SetSize(480, 320, webview.HintFixed)
	w.SetHtml("Thanks for using webview!")

	w.Run()
}
