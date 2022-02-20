package whfs

import (
	"io"
)

type FS interface {
	Open(path string) (io.ReadCloser, error)
	CreateWithDirs(path string) (io.WriteCloser, error)
	Remove(root string) error
	UserHomeDir() (string, error)
}
