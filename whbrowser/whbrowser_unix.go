//go:build !wasm
// +build !wasm

package whbrowser

import "github.com/pkg/browser"

func OpenURL(url string) error {
	return browser.OpenURL(url)
}
