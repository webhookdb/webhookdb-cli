package ask

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
