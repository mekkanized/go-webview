package main

import (
	"log"

	"github.com/Mekkabotics/go-webkit"
)

func main() {
  w := webkit.NewWithOptions(webkit.WebViewOptions{
    Debug: true,
  })
  if w == nil {
    log.Fatalln("Failed to load webkit.")
  }
  defer w.Destroy()

  w.SetTitle("Basic Example")
  w.SetSize(480, 320, webkit.HintNone)
  w.SetHtml("Thanks for using webview!")

  w.Run()
}

