//go:build wasm
// +build wasm

package whfs

import (
	"io"
)

func New() FS {
	return jsfs{}
}

type jsfs struct{}

func (o jsfs) Remove(root string) error {
	panic("implement me")
}

func (o jsfs) UserHomeDir() (string, error) {
	panic("implement me")
}

func (o jsfs) Open(path string) (io.ReadCloser, error) {
	panic("implement me")
}

func (o jsfs) CreateWithDirs(path string) (io.WriteCloser, error) {
	panic("implement me")
}
