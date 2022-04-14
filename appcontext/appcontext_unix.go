//go:build !wasm
// +build !wasm

package appcontext

const shortSession = ""

func platformUserAgent() string {
	return ""
}
