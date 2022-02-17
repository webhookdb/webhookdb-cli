package whselfupdate

import (
	"errors"
	"github.com/blang/semver"
)

//goland:noinspection GoErrorStringFormat
var ErrUnsupported = errors.New("Sorry, updating is not supported in this environment.")

type Release interface {
	Version() semver.Version
	AssetURL() string
}
