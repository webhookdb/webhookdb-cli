package prefs

import (
	"encoding/json"
	"github.com/lithictech/go-aperitif/convext"
	"github.com/lithictech/webhookdb-cli/types"
	"github.com/lithictech/webhookdb-cli/whfs"
	"github.com/pkg/errors"
	"io/fs"
	"io/ioutil"
	"path/filepath"
)

type Namespace string

type GlobalPrefs struct {
	Namespaces map[Namespace]Prefs `json:"namespaces"`
}

func (p *GlobalPrefs) GetNS(namespace string) Prefs {
	return p.Namespaces[Namespace(namespace)]
}

func (p *GlobalPrefs) SetNS(namespace string, prefs Prefs) {
	p.Namespaces[Namespace(namespace)] = prefs
}

func (p *GlobalPrefs) ClearNS(namespace string) {
	delete(p.Namespaces, Namespace(namespace))
}

type Prefs struct {
	AuthToken  types.AuthToken    `json:"auth_token"`
	Email      string             `json:"email"`
	CurrentOrg types.Organization `json:"current_org"`
}

func (p Prefs) ChangeOrg(org types.Organization) Prefs {
	p.CurrentOrg = org
	return p
}

func getDir(pfs whfs.FS) string {
	home, err := pfs.UserHomeDir()
	convext.Must(err)
	return filepath.Join(home, ".webhookdb")
}

func getPath(pfs whfs.FS) string {
	return filepath.Join(getDir(pfs), "config")
}

func Load(pfs whfs.FS) (*GlobalPrefs, error) {
	p := &GlobalPrefs{Namespaces: make(map[Namespace]Prefs, 1)}
	path := getPath(pfs)
	f, err := pfs.Open(path)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		return p, nil
	} else if err != nil {
		return p, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return p, errors.Wrap(err, "reading "+path)
	}
	if err := json.Unmarshal(b, &p); err != nil {
		return p, errors.Wrap(err, "unmarshalling into prefs: "+string(b))
	}
	return p, nil
}

func Save(pfs whfs.FS, p *GlobalPrefs) error {
	f, err := pfs.CreateWithDirs(getPath(pfs))
	if err != nil {
		return err
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(p); err != nil {
		return err
	}
	return nil
}

func DeleteAll(pfs whfs.FS) error {
	err := pfs.Remove(getPath(pfs))
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}
	return err
}
