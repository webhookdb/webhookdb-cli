//go:build !wasm
// +build !wasm

package cmd

import (
	"log"
	"os"
)

func Execute() {
	err := BuildApp().Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
