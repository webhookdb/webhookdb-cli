//go:build wasm
// +build wasm

package appcontext

import "syscall/js"

const shortSession = "1"

func platformUserAgent() string {
	var nav = js.Global().Get("navigator")
	if !nav.Truthy() {
		return ""
	}
	var ua = nav.Get("userAgent")
	if !ua.Truthy() {
		return ""
	}
	return ua.String()
}
