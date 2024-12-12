package pg

import (
	"strings"

	"github.com/dronm/crudifier/types"
)

type PgFilters []PgFilter

// Add(fieldId string, value interface{}, operation types.SQLFilterOperator, join types.FilterJoin)
func (f *PgFilters) Add(fieldId string,
	value interface{}, operator types.SQLFilterOperator,
	join types.FilterJoin) {

	*f = append(*f, PgFilter{fieldID: fieldId, value: value, join: join, operator: operator})
}

func (f PgFilters) Len() int {
	return len(f)
}

func (f PgFilters) SQL(queryParams *[]interface{}) string {
	if len(f) == 0 {
		return ""
	} else if len(f) == 1 {
		return " WHERE " + f[0].SQL(queryParams)
	}

	var sqlSt strings.Builder
	sqlSt.WriteString(" WHERE ")
	for i, filter := range f {
		if i > 0 {
			//default
			if filter.join == "" {
				filter.join = types.SQL_FILTER_JOIN_AND
			}
			sqlSt.WriteString(" " + string(filter.join) + " ")
		}
		sqlSt.WriteString("(" + filter.SQL(queryParams) + ")")
	}
	return sqlSt.String()
}
