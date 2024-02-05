//go:build !wasm
// +build !wasm

package ask

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func New() Ask {
	a := ioasker{
		ch: make(chan []byte, 4),
	}
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err == nil {
		a.isInteractive = true
		a.input = tty
		a.output = tty
	} else {
		a.isInteractive = false
		a.input = os.Stdin
		a.output = os.Stdout
	}
	a.runReader()
	return &a
}

type ioasker struct {
	isInteractive bool
	input         *os.File
	output        *os.File
	ch            chan []byte
}

func (i *ioasker) Ask(prompt string) (string, error) {
	err := i.print(prompt)
	if err != nil {
		return "", err
	}
	return i.readinput()
}

func (i *ioasker) HiddenAsk(prompt string) (string, error) {
	err := i.stty("-echo")
	if err != nil {
		return "", err
	}
	defer func(i *ioasker, args ...string) {
		_ = i.stty(args...)
	}(i, "echo")
	err = i.print(prompt)
	if err != nil {
		return "", err
	}
	defer func(i *ioasker, str string) {
		_ = i.print(str)
	}(i, "\n")
	return i.readinput()
}

func (i *ioasker) Feedback(line string) {
	fmt.Fprintln(i.output, line)
}

func (i *ioasker) print(str string) error {
	if strings.HasSuffix(str, ":") {
		str += " "
	}
	_, err := fmt.Fprint(i.output, str)
	return err
}

// Reading input is finnicky, because the terminal uses line buffering,
// and Read blocks. So we can't tell the difference between someone pressing Enter,
// and someone pasting (we don't get EOF so even something like io.ReadAll blocks).
//
// But because we have a particular pattern- prompt for input, then look for what is entered-
// we can be clever (oh no...). We run a goroutine that reads from input.
// This readinput function looks for either text, or times out.
// If we time out before getting any text, we keep looping;
// but if we time out after getting some text, assume the goroutine is blocked on Read.
//
// It gets a bit worse. We can finish the input read after the first \n in normal cases,
// but if we're pasting, we may have many newlines. In that case, we only want to break
// after a blank line. Otherwise we can end up with hanging text in the buffer,
// if you paste multiple lines but the final line doesn't end with a newline.
func (i *ioasker) readinput() (string, error) {
	var buffer bytes.Buffer
	lines := 0
readloop:
	for {
		select {
		case readBytes, closed := <-i.ch:
			if !closed {
				// Exit the loop if the reader closed the channel.
				break readloop
			} else {
				// Write the line, and record that we wrote some data.
				buffer.Write(readBytes)
				lines++
			}
		case <-time.After(50 * time.Millisecond):
			// 50ms should be more than enough for the reader goroutine to write to the channel.
			// If we've timed out and only written one line, we can finish the loop.
			// If we've timed out and written multiple lines, only finish the loop if we are ending
			// with a double newline. See docstring.
			// Break if we've read, loop if we haven't
			if lines == 0 {
				// keep going
			} else if lines == 1 {
				break readloop
			} else {
				bufbytes := buffer.Bytes()
				buflen := buffer.Len()
				hasConfirmationNewline := bufbytes[buflen-2] == '\n' && bufbytes[buflen-1] == '\n'
				if hasConfirmationNewline {
					break readloop
				}
			}
		}
	}
	return string(buffer.Bytes()), nil
}

// Code taken and adapted from https://github.com/miquella/ask
func (i *ioasker) stty(args ...string) error {
	// don't do anything if we're non-interactive
	if !i.isInteractive {
		return nil
	}
	cmd := exec.Command("stty", args...)
	// if stty wasn't found in path, try hard-coding it
	if filepath.Base(cmd.Path) == cmd.Path {
		cmd.Path = "/bin/stty"
	}
	cmd.Stdin = i.input
	cmd.Stdout = i.output
	return cmd.Run()
}

func (i *ioasker) runReader() {
	reader := bufio.NewReader(i.input)
	go func() {
		for {
			s, err := reader.ReadBytes('\n')
			if err != nil {
				close(i.ch)
				return
			}
			i.ch <- s
		}
	}()
}
