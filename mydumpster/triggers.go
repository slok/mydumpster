package mydumpster

import (
	"fmt"
	"strings"
)

type Trigger struct {
	TableDst      Table
	TableSrcName  string
	TableSrcField string
	TableDstField string
	DumpAll       bool
}

func (t *Trigger) SelectQueryFromRowsStr(rows [][]string, columns []string) string {
	// get position
	pos := SearchStr(columns, t.TableSrcField)

	// Get all the identifiers
	ids := make([]string, len(rows))
	for i, v := range rows {
		ids[i] = v[pos]
	}

	return t.SelectQueryStr(ids)
}

func (t *Trigger) SelectQueryStr(ids []string) string {
	return fmt.Sprintf(IN_FMT, t.TableDstField, strings.Join(ids, ", "))
}
