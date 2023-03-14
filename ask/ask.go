package ask

import "errors"

// ErrBreak represents a 'break' in environments that do not have
// a real break, like our WASM intepreter.
// Basically what happens in all cases would be:
//   - The Ask implementation panics with ErrBreak
//   - The handler for this invocation recovers from panics,
//     and checks if recover() == ErrBreak.
//   - If it is, we should indicate to the caller (the web terminal)
//     that the client used a break.
//   - The client should interpret this how it needs,
//     which is going to emulate a normal terminal.
var ErrBreak = errors.New("psuedo-break error")
var BreakSentinel = "__whdb_break"

// Ask encapsulates the prompting and feedback of asking for input.
type Ask interface {
	// Ask prints the given prompt and asks for input.
	// It returns the input, or an error.
	Ask(prompt string) (string, error)
	// HiddenAsk is like Ask, but the input is hidden if possible.
	HiddenAsk(prompt string) (string, error)
	// Feedback prints the given feedback.
	// This is usually used to print any output after the result of an Ask.
	Feedback(line string)
}

func HiddenPrompt(prefix string) string {
	return prefix + " *** "
}
