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
	"sync"
	"time"
)

func New() Ask {
	globalIOOnce.Do(func() {
		g := globalIO{
			Chan: make(chan []byte, 4),
		}
		tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
		if err == nil {
			g.IsInteractive = true
			g.Input = tty
			g.Output = tty
		} else {
			g.IsInteractive = false
			g.Input = os.Stdin
			g.Output = os.Stdout
		}
		g.run()
		globalIOInst = &g
	})
	a := ioasker{
		g: globalIOInst,
	}
	return &a
}

type ioasker struct {
	g *globalIO
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
	fmt.Fprintln(i.g.Output, line)
}

func (i *ioasker) print(str string) error {
	if strings.HasSuffix(str, ":") {
		str += " "
	}
	_, err := fmt.Fprint(i.g.Output, str)
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
		case readBytes, closed := <-i.g.Chan:
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
	if !i.g.IsInteractive {
		return nil
	}
	cmd := exec.Command("stty", args...)
	// if stty wasn't found in path, try hard-coding it
	if filepath.Base(cmd.Path) == cmd.Path {
		cmd.Path = "/bin/stty"
	}
	cmd.Stdin = i.g.Input
	cmd.Stdout = i.g.Output
	return cmd.Run()
}

var globalIOInst *globalIO
var globalIOOnce = &sync.Once{}

// We can only set up ONE reader for the input, since it uses ReadBytes which is blocking.
// If we set up multiple readers in ReadBytes, they end up all listening for the prompt;
// then for example by the third Ask call, we have a third reader,
// and the input is lost.
// So we need a single global reader for whatever we're using for input/output.
// This should be global for the process anyway (ask.New() takes no arguments),
// so is not a big deal.
type globalIO struct {
	IsInteractive bool
	Input         *os.File
	Output        *os.File
	Chan          chan []byte
}

func (g *globalIO) run() {
	reader := bufio.NewReader(g.Input)
	go func() {
		for {
			s, err := reader.ReadBytes('\n')
			if err != nil {
				close(g.Chan)
				return
			}
			g.Chan <- s
		}
	}()
}
