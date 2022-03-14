package formatting

import (
	"encoding/json"
	"fmt"
	"github.com/lithictech/go-aperitif/convext"
	"io"
)

type Block struct {
	Type     string          `json:"type"`
	RawValue json.RawMessage `json:"value"`
}

const (
	BlockLine  = "line"
	BlockTable = "table"
)

func (b Block) LineValue() (s string) {
	convext.MustUnmarshal(b.RawValue, &s)
	return
}

func (b Block) TableValue() (t TabularResponse) {
	convext.MustUnmarshal(b.RawValue, &t)
	return
}

type Blocks []Block

func (bs Blocks) WriteTo(w io.Writer) (int64, error) {
	for _, b := range bs {
		if b.Type == BlockLine {
			fmt.Fprintln(w, b.LineValue())
		} else if b.Type == BlockTable {
			if err := tableWriteTabular(b.TableValue(), w); err != nil {
				return 0, err
			}
		}
	}
	return 1, nil
}
