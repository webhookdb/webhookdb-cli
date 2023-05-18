//go:build !wasm
// +build !wasm

package whselfupdate

import (
	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

func DetectVersion(slug string, version string) (Release, bool, error) {
	r, ok, err := selfupdate.DetectVersion(slug, version)
	var whr Release
	if r != nil {
		whr = whrelease{r: r}
	}
	return whr, ok, err
}

func UpdateTo(assetURL, cmdPath string) error {
	return selfupdate.UpdateTo(assetURL, cmdPath)
}

type whrelease struct {
	r *selfupdate.Release
}

func (w whrelease) Version() semver.Version {
	return w.r.Version
}

func (w whrelease) AssetURL() string {
	return w.r.AssetURL
}

func init() {
	selfupdate.AssetGoOSAliases["darwin"] = []string{"macos"}
}
