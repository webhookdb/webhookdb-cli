package formatting

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/lithictech/webhookdb-cli/types"
	"github.com/olekukonko/tablewriter"
	"io"
)

func CollectionResponseHeaders(cr types.CollectionResponse) ([]types.KeyAndName, []string) {
	hdrs := cr.DisplayHeaders()
	names := make([]string, len(hdrs))
	for i, h := range hdrs {
		names[i] = h.Name
	}
	return hdrs, names
}

func CollectionResponseStringRow(headers []types.KeyAndName, item map[string]interface{}, row []string) {
	for i, h := range headers {
		row[i] = ToString(item[h.Key])
	}
}

type Format struct {
	// The string value of this format to use in a CLI flag, like 'csv' or 'json'.
	FlagValue string
	// Write the API collection response with "display_headers" and "items" to w.
	WriteCollection func(io.Writer, types.CollectionResponse) error
}

var JSON = Format{
	FlagValue: "json",
	WriteCollection: func(w io.Writer, r types.CollectionResponse) error {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(r.Items())
	},
}

var CSV = Format{
	FlagValue: "csv",
	WriteCollection: func(w io.Writer, cr types.CollectionResponse) error {
		headers, headerNames := CollectionResponseHeaders(cr)
		cw := csv.NewWriter(w)
		if err := cw.Write(headerNames); err != nil {
			return err
		}
		row := make([]string, len(headers))
		for _, item := range cr.Items() {
			CollectionResponseStringRow(headers, item, row)
			if err := cw.Write(row); err != nil {
				return err
			}
		}
		cw.Flush()
		return cw.Error()
	},
}

var Table = Format{
	FlagValue: "table",
	WriteCollection: func(w io.Writer, cr types.CollectionResponse) error {
		items := cr.Items()
		if len(items) == 0 {
			return nil
		}
		table := tablewriter.NewWriter(w)
		headers, headerNames := CollectionResponseHeaders(cr)
		table.SetHeader(headerNames)
		ConfigureTableWriter(table)
		row := make([]string, len(headers))
		for _, item := range items {
			CollectionResponseStringRow(headers, item, row)
			table.Append(row)
		}
		table.Render()
		return nil
	},
}

var Formats = []Format{JSON, CSV, Table}

type TabularData struct {
	Headers []string   `json:"headers"`
	Rows    [][]string `json:"rows"`
}

func (t TabularData) Write(w io.Writer) error {
	table := tablewriter.NewWriter(w)
	table.SetHeader(t.Headers)
	ConfigureTableWriter(table)
	for _, row := range t.Rows {
		table.Append(row)
	}
	table.Render()
	return nil
}

func LookupByFlag(v string) (Format, bool) {
	for _, f := range Formats {
		if f.FlagValue == v {
			return f, true
		}
	}
	return Format{}, false
}

func FormatFlagValues() []string {
	values := make([]string, len(Formats))
	for i, f := range Formats {
		values[i] = f.FlagValue
	}
	return values
}

func ConfigureTableWriter(table *tablewriter.Table) {
	table.SetBorder(false)
	table.SetRowSeparator("")
	table.SetColumnSeparator("")
	table.SetCenterSeparator("")
	table.SetHeaderLine(false)
}

func ToString(i interface{}) string {
	if str, ok := i.(fmt.Stringer); ok {
		return str.String()
	}
	if str, ok := i.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", i)
}

func FprintlnNonempty(w io.Writer, a ...interface{}) {
	if len(a) == 0 || (len(a) == 1 && a[0] == "") {
		return
	}
	fmt.Fprintln(w, a...)
}
