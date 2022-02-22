//go:build wasm
// +build wasm

package ask

import (
	"fmt"
	"syscall/js"
)

func New() Ask {
	a := jsask{}
	return &a
}

var feedbackNotSetWarned = false
var promptNotSetWarned = false

type jsask struct{}

func (j jsask) Ask(prompt string) (string, error) {
	return j.ask(prompt, false)
}

func (j jsask) HiddenAsk(prompt string) (string, error) {
	return j.ask(prompt, true)
}

func (j jsask) ask(prompt string, hidden bool) (string, error) {
	wasmPrompt := js.Global().Get("wasmPropt")
	if !wasmPrompt.Truthy() {
		if promptNotSetWarned {
			promptNotSetWarned = true
			fmt.Println("wasmPrompt global function is not set, falling back to prompt builtin")
		}
		result := js.Global().Call("prompt", prompt)
		if !result.Truthy() {
			return "", nil
		}
		return result.String(), nil
	}
	ch := make(chan string)
	cb := func(this js.Value, args []js.Value) interface{} {
		ch <- args[0].String()
		return nil
	}
	go func() {
		wasmPrompt.Invoke(prompt, hidden, js.FuncOf(cb))
	}()
	return <-ch, nil
}

func (j jsask) Feedback(line string) {
	fb := js.Global().Get("wasmFeedback")
	if fb.IsNull() || fb.IsUndefined() {
		if feedbackNotSetWarned {
			feedbackNotSetWarned = true
			fmt.Println("wasmFeedback global function is not set, printing feedback to console")
		}
		fmt.Println(line)
		return
	}
	fb.Invoke(line)
}
