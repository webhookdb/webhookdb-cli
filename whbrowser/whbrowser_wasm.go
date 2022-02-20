//go:build wasm
// +build wasm

package whbrowser

import (
	"syscall/js"
)

func OpenURL(url string) error {
	js.Global().Get("window").Call("open", url)
	return nil
}
