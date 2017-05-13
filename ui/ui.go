package main

import (
	"github.com/gopherjs/gopherjs/js"
)

func main() {
	js.Global.Get("console").Call("log", struct{ S string }{"Start ui"})

	js.Global.Get("document").Call("addEventListener", "astilectron-ready", func() {
		js.Global.Get("astilectron").Call("listen", func(message string) {
			js.Global.Get("document").Call("getElementById", "message").Set("innerHTML", message)
			js.Global.Get("astilectron").Call("send", "Pong gopherjs!")
		})
	})
}
