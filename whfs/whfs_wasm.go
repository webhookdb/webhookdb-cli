//go:build wasm
// +build wasm

package whfs

import (
	"io"
	"io/fs"
	"os"
	"strings"
	"syscall/js"
)

func New() FS {
	return jsfs{}
}

type jsfs struct{}

func (o jsfs) Remove(root string) error {
	return os.Remove(root)
}

func (o jsfs) UserHomeDir() (string, error) {
	return "/home", nil
}

func (o jsfs) Open(path string) (io.ReadCloser, error) {
	return newJSStorageIOReader(path, "sessionStorage")
}

func (o jsfs) CreateWithDirs(path string) (io.WriteCloser, error) {
	return newJSStorageIOWriter(path, "sessionStorage"), nil
}

func newJSStorageIOReader(path, storage string) (io.ReadCloser, error) {
	j := &jsStorageIO{path: path, storage: storage}
	v := j.getItem()
	if v.IsNull() || v.IsUndefined() {
		return nil, fs.ErrNotExist
	}
	j.readBuff = strings.NewReader(v.String())
	return j, nil
}

func newJSStorageIOWriter(path, storage string) io.WriteCloser {
	j := &jsStorageIO{path: path, storage: storage}
	j.writeBuff = &strings.Builder{}
	return j
}

type jsStorageIO struct {
	path      string
	storage   string
	writeBuff *strings.Builder
	readBuff  *strings.Reader
}

func (j *jsStorageIO) fqpath() string {
	return "whdb." + j.path
}

func (j *jsStorageIO) getItem() js.Value {
	return js.Global().Get(j.storage).Call("getItem", j.fqpath())
}

func (j *jsStorageIO) setItem(value string) {
	js.Global().Get(j.storage).Call("setItem", j.fqpath(), value)
}

func (j *jsStorageIO) Read(p []byte) (n int, err error) {
	return j.readBuff.Read(p)
}

func (j *jsStorageIO) Write(p []byte) (n int, err error) {
	return j.writeBuff.Write(p)
}

func (j jsStorageIO) Close() error {
	if j.writeBuff != nil {
		j.setItem(j.writeBuff.String())
	}
	return nil
}
