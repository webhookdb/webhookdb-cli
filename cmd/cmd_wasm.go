//go:build wasm
// +build wasm

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/urfave/cli/v2"
	"github.com/webhookdb/webhookdb-cli/appcontext"
	"github.com/webhookdb/webhookdb-cli/ask"
	"github.com/webhookdb/webhookdb-cli/prefs"
	"os"
	"runtime/debug"
	"strings"
	"syscall/js"
)

// Set this explicitly since os.Args[0] will always be 'js'
var helpName = "webhookdb"

func Execute() {
	c := make(chan struct{}, 0)
	js.Global().Set("webhookdbRunGo", js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		go webhookdbRunGo(args)
		return nil
	}))
	<-c
}

// Wrap a CLI run in the work to go from/to JS
func webhookdbRunGo(arguments []js.Value) {
	onComplete := arguments[1]
	defer func() {
		if r := recover(); r != nil {
			if r == ask.ErrBreak {
				onComplete.Invoke(fmt.Sprintf(`{"break":true}`))
			} else {
				fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
				onComplete.Invoke(fmt.Sprintf(`{"stderr":"Sorry, something went wrong. Check the console for more information, or please try again.", "panic":true}`))
			}
		}
	}()
	if err := webhookdbSetenv(arguments[0]); err != nil {
		onComplete.Invoke(fmt.Sprintf(`{"stderr":"During setenv: %s"}`, err.Error()))
		return
	}
	jsCliArgs := arguments[2:]
	strCliArgs := make([]string, len(jsCliArgs))
	for i, arg := range jsCliArgs {
		strCliArgs[i] = arg.String()
	}
	println(fmt.Sprintf("Running: %s", strings.Join(strCliArgs, " ")))
	stdout, stderr := handleRunGo(strCliArgs)
	res := map[string]string{
		"stdout": stdout,
		"stderr": stderr,
	}
	j, err := json.Marshal(res)
	if err != nil {
		onComplete.Invoke(fmt.Sprintf(`{"stderr": "During marshaling: %s"}`, err.Error()))
		return
	}
	onComplete.Invoke(string(j))
}

// Run a CLI app with in-memory things and return stdout/stderr
func handleRunGo(arguments []string) (string, string) {
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	app := BuildApp()
	app.Writer = stdout
	app.ErrWriter = stderr
	cli.ErrWriter = stderr
	cli.OsExiter = func(c int) {
		return
	}
	appErr := app.Run(arguments)
	// Log the app error to the console, not return it.
	// It's output to stderr already anyway.
	if appErr != nil {
		fmt.Println(appErr)
	}
	return stdout.String(), stderr.String()
}

// unmarshal a marshaled environment and update the env,
// so config can be loaded from it.
func webhookdbSetenv(v js.Value) error {
	var env map[string]string
	if err := json.Unmarshal([]byte(v.String()), &env); err != nil {
		return err
	}
	for k, v := range env {
		_ = os.Setenv(k, v)
	}
	return nil
}

func wasmUpdateAuthDisplay(p prefs.Prefs) {
	cb := js.Global().Get("webhookdbOnAuthed")
	if !cb.Truthy() {
		fmt.Println("webhookdbOnAuthed not registered, cannot call back.")
		return
	}
	j, err := json.Marshal(map[string]interface{}{"email": p.Email, "org_name": p.CurrentOrg.Name, "org_key": p.CurrentOrg.Key})
	if err != nil {
		fmt.Println("json marshal error:", err)
		return
	}
	cb.Invoke(string(j))
}

var IsTty = false

func onServerError(*cli.Context, appcontext.AppContext, *resty.Response) error {
	// Eventually we may want to use Sentry from the browser.
	return nil
}
