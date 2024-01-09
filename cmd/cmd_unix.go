//go:build !wasm
// +build !wasm

package cmd

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/go-resty/resty/v2"
	"github.com/lithictech/go-aperitif/logctx"
	"github.com/urfave/cli/v2"
	"github.com/webhookdb/webhookdb-cli/appcontext"
	"github.com/webhookdb/webhookdb-cli/config"
	"github.com/webhookdb/webhookdb-cli/prefs"
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
	if defaultConfig.Privacy {
		return false
	}
	dsn := defaultConfig.SentryDsn
	if dsn == "" && !strings.HasPrefix(defaultConfig.ApiHost, "http://") {
		dsn = config.SentryDsnProd
	} else if !strings.HasPrefix(dsn, "http") {
		return false
	}
	transport := sentry.NewHTTPSyncTransport()
	transport.Timeout = 2 * time.Second
	err := sentry.Init(sentry.ClientOptions{
		Dsn:       dsn,
		Release:   config.BuildSha,
		Debug:     defaultConfig.Debug,
		Transport: transport,
	})
	if err != nil {
		log.Printf("Error starting Sentry: %s", err)
		return false
	}
	sentry.CurrentHub().ConfigureScope(func(sc *sentry.Scope) {
		sc.SetTag("application", "cli")
	})
	return true
}

func wasmUpdateAuthDisplay(_ prefs.Prefs) {}

var IsTty = logctx.IsTty()

func onServerError(c *cli.Context, appCtx appcontext.AppContext, resp *resty.Response) error {
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(sc *sentry.Scope) {
		sc.SetUser(sentry.User{
			Email: appCtx.Prefs.Email,
		})
		sc.SetContext("cli_command", c.Command.FullName())
		sc.SetRequest(resp.Request.RawRequest)
		sc.SetContext("response_code", resp.StatusCode())
		sc.SetContext("response_body", string(resp.Body()))
	})
	hub.CaptureException(serverError{fmt.Sprintf("CLI got an error from the server: %d", resp.StatusCode())})
	return nil
}

type serverError struct{ m string }

func (e serverError) Error() string {
	return e.m
}
