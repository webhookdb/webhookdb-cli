//go:build wasm
// +build wasm

package whselfupdate

func DetectVersion(slug string, version string) (Release, bool, error) {
	return nil, false, ErrUnsupported
}

func UpdateTo(assetURL, cmdPath string) error {
	return ErrUnsupported
}
