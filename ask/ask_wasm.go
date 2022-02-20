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

type jsask struct{}

func (j jsask) Ask(prompt string) (string, error) {
	result := js.Global().Call("prompt", prompt)
	return result.String(), nil
}

func (j jsask) HiddenAsk(prompt string) (string, error) {
	result := js.Global().Call("prompt", prompt)
	return result.String(), nil
}

var feedbackNotSetWarned = false

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
