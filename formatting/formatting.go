package formatting

import (
	"encoding/csv"
	"encoding/json"
	"github.com/olekukonko/tablewriter"
	"io"
)

type Format struct {
	FlagValue          string
	ApiRequestValue    string
	ApiResponsePtr     func() interface{}
	WriteApiResponseTo func(o interface{}, w io.Writer) error
	WriteTabular       func(t TabularResponse, w io.Writer) error
}

var JSON = Format{
	FlagValue:       "json",
	ApiRequestValue: "object",
	ApiResponsePtr: func() interface{} {
		return &ObjectResponse{}
	},
	WriteApiResponseTo: func(o interface{}, w io.Writer) error {
		return jsonWriteObj(o, w)
	},
	WriteTabular: func(t TabularResponse, w io.Writer) error {
		m := make([]map[string]interface{}, len(t.Rows))
		for rowIdx, row := range t.Rows {
			rowObj := make(map[string]interface{}, len(t.Headers))
			for headerIdx, header := range t.Headers {
				rowObj[header] = row[headerIdx]
			}
			m[rowIdx] = rowObj
		}
		return jsonWriteObj(m, w)
	},
}

func jsonWriteObj(o interface{}, w io.Writer) error {
	return json.NewEncoder(w).Encode(o)
}

var CSV = Format{
	FlagValue:       "csv",
	ApiRequestValue: "table",
	ApiResponsePtr: func() interface{} {
		return &TabularResponse{}
	},
	WriteApiResponseTo: func(o interface{}, w io.Writer) error {
		t := o.(*TabularResponse)
		return csvWriteTabular(*t, w)
	},
	WriteTabular: csvWriteTabular,
}

func csvWriteTabular(t TabularResponse, w io.Writer) error {
	cw := csv.NewWriter(w)
	if err := cw.Write(t.Headers); err != nil {
		return err
	}
	return cw.WriteAll(t.Rows)
}

var Table = Format{
	FlagValue:       "table",
	ApiRequestValue: "table",
	ApiResponsePtr: func() interface{} {
		return &TabularResponse{}
	},
	WriteApiResponseTo: func(o interface{}, w io.Writer) error {
		t := o.(*TabularResponse)
		return tableWriteTabular(*t, w)
	},
	WriteTabular: tableWriteTabular,
}

func tableWriteTabular(t TabularResponse, w io.Writer) error {
	table := tablewriter.NewWriter(w)
	table.SetHeader(t.Headers)
	ConfigureTableWriter(table)
	for _, row := range t.Rows {
		table.Append(row)
	}
	table.Render()
	return nil
}

type TabularResponse struct {
	Headers []string   `json:"headers"`
	Rows    [][]string `json:"rows"`
}

type ObjectResponse map[string]interface{}

var Formats = []Format{JSON, CSV, Table}

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
