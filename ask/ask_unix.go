//go:build !wasm
// +build !wasm

package ask

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func New() Ask {
	a := ioasker{}
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
	return &a
}

type ioasker struct {
	isInteractive bool
	input         *os.File
	output        *os.File
}

func (i *ioasker) Ask(prompt string) (string, error) {
	err := i.print(prompt)
	if err != nil {
		return "", err
	}
	return i.readline()
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
	return i.readline()
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

func (i *ioasker) readline() (string, error) {
	var err error
	var buffer bytes.Buffer
	var b [1]byte
	for {
		var n int
		n, err = i.input.Read(b[:])
		if b[0] == '\n' {
			break
		}
		if n > 0 {
			buffer.WriteByte(b[0])
		}
		if n == 0 || err != nil {
			break
		}
	}

	if err != nil && err != io.EOF {
		return "", err
	}
	return string(buffer.Bytes()), err
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
