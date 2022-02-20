//go:build !wasm
// +build !wasm

package cmd

import (
	"github.com/lithictech/webhookdb-cli/prefs"
	"log"
	"os"
)

func Execute() {
	err := BuildApp().Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func wasmUpdateAuthDisplay(_ prefs.Prefs) {}
