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
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("https://html5test.com/")
	w.Run()
}
