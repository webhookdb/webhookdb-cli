//go:build !wasm
// +build !wasm

package cmd

import (
	"github.com/getsentry/sentry-go"
	"github.com/lithictech/go-aperitif/logctx"
	"github.com/lithictech/webhookdb-cli/config"
	"github.com/lithictech/webhookdb-cli/prefs"
	"log"
	"os"
	"strings"
	"time"
)

// Let this get defaulted to the executed command name
var helpName = ""

func Execute() {
	if configureSentry() {
		defer func() {
			if r := recover(); r != nil {
				sentry.CurrentHub().Recover(r)
				sentry.Flush(time.Second * 2)
				panic(r)
			}
		}()
	}
	err := BuildApp().Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func configureSentry() bool {
	defaultConfig := config.LoadConfig()
	dsn := defaultConfig.SentryDsn
	if dsn == "" && !strings.HasPrefix(defaultConfig.ApiHost, "http://") {
		dsn = "https://3e125fd192c34979b2f1a4a5ceb9abd6@o292308.ingest.sentry.io/6224206"
	} else if !strings.HasPrefix(dsn, "http") {
		return false
	}
	err := sentry.Init(sentry.ClientOptions{
		Dsn:     dsn,
		Release: config.BuildSha,
		Debug:   defaultConfig.Debug,
	})
	if err != nil {
		log.Printf("Error starting Sentry: %s", err)
		return false
	}
	return true
}

func wasmUpdateAuthDisplay(_ prefs.Prefs) {}

var IsTty = logctx.IsTty()
