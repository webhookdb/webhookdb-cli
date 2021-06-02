package prefs

import (
	"encoding/json"
	"github.com/lithictech/go-aperitif/convext"
	"os"
	"path/filepath"
)

type Prefs struct {
	AuthCookie string `json:"auth_cookie"`
	CurrentOrg string `json:"current_org"`
}

func getDir() string {
	home, err := os.UserHomeDir()
	convext.Must(err)
	return filepath.Join(home, ".webhookdb")
}
func getPath() string {
	return filepath.Join(getDir(), "config")
}

func Load() (Prefs, error) {
	p := Prefs{}
	path := getPath()
	f, err := os.Open(path)
	if err != nil && os.IsNotExist(err) {
		return p, nil
	} else if err != nil {
		return p, err
	}
	if err := json.NewDecoder(f).Decode(&p); err != nil {
		return p, err
	}
	return p, nil
}

func Save(p Prefs) error {
	if err := os.MkdirAll(getDir(), os.ModePerm); err != nil {
		return err
	}
	f, err := os.Create(getPath())
	if err != nil {
		return err
	}
	if err := json.NewEncoder(f).Encode(p); err != nil {
		return err
	}
	return nil
}
