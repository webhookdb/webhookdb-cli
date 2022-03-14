package formatting

import (
  "fmt"
  "github.com/olekukonko/tablewriter"
)

func NewAutoSizingTableWriter(tw *tablewriter.Table) *AutoSizingTableWriter {
  return &AutoSizingTableWriter{
    Table: tw,
  }
}

type AutoSizingTableWriter struct {
  Table *tablewriter.Table
  minWidths []int
}

func (tw *AutoSizingTableWriter) SetHeader(keys []string) {
  tw.resize(keys)
  tw.Table.SetHeader(keys)
}

func (tw *AutoSizingTableWriter) Append(row []string) {
  tw.resize(row)
  tw.Table.Append(row)
}

func (tw *AutoSizingTableWriter) AppendBulk(rows [][]string) {
  for _, row := range rows {
    tw.Append(row)
  }
}

func (tw *AutoSizingTableWriter) Render() {
  tw.SetSizes()
  tw.Table.Render()
}

func (tw *AutoSizingTableWriter) SetSizes() {
  termwidth := TermWidth()
  if termwidth <= 0 {
    return
  }
  totalWidth := len(tw.minWidths) * 3 + 3 // Each border is up to 3 spaces, add a little buffer too
  for _, w := range tw.minWidths {
    totalWidth += w
  }
  // If term < total, we need to scale down.
  // Otherwise we do not need to scale up/grow.
  scale :=  float64(termwidth) / float64(totalWidth)
  if scale > 1 {
    scale = 1
  }
  adjustedwidth := 0
  for _, w := range tw.minWidths {
    aw := int(float64(w) * scale)
    adjustedwidth += aw
    //tw.Table.SetColMinWidth(i, aw)
  }
  if adjustedwidth >= termwidth {
    panic(fmt.Sprintf("%d > %d", adjustedwidth, termwidth))
  }
  return
}

func (tw *AutoSizingTableWriter) resize(row []string) {
  if tw.minWidths == nil {
    tw.minWidths = make([]int, len(row))
  }
  if len(row) != len(tw.minWidths) {
    return
  }
  for i, c := range row {
    tw.minWidths[i] = intmax(tw.minWidths[i], len(c) + 3) // without the addition we get weird diffs between rows in some terminals
  }
  return
}

func intmax(x, y int) int {
  if x > y {
    return x
  }
  return y
}
