//go:build !wasm
// +build !wasm

package whfs

import (
	"io"
	"os"
	"path/filepath"
)

func New() FS {
	return osfs{}
}

type osfs struct{}

func (o osfs) Remove(root string) error {
	return os.Remove(root)
}

func (o osfs) UserHomeDir() (string, error) {
	return os.UserHomeDir()
}

func (o osfs) Open(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

func (o osfs) CreateWithDirs(path string) (io.WriteCloser, error) {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return nil, err
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}
