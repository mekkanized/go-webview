package main

import (
	"fmt"
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
	w.SetTitle("Demo Example")
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("https://html5test.com/")
	w.Bind("myfunc", func(s string) {
		fmt.Printf("test: %s\n", s)
	})
	w.Run()
}
