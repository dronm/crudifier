package pg

import (
	"fmt"
	"strings"

	"github.com/dronm/crudifier/types"
)

type PgFilters []PgFilter

// Add(fieldId string, value any, operation types.SQLFilterOperator, join types.FilterJoin)
func (f *PgFilters) Add(fieldId string,
	value any, operator types.SQLFilterOperator,
	join types.FilterJoin) {

	*f = append(*f, PgFilter{fieldID: fieldId, value: value, join: join, operator: operator})
}

func (f *PgFilters) AddFullTextSearch(fieldId string, value any, join types.FilterJoin) {
	*f = append(*f, PgFilter{
		value: value,
		join: join,
		expression: fmt.Sprintf("%s @@ to_tsquery('russian', {{PARAM}})", fieldId),
	})
}

func (f *PgFilters) AddArrayInclude(fieldId string, value any, join types.FilterJoin) {
	*f = append(*f, PgFilter{
		value: value,
		join: join,
		expression: fmt.Sprintf("%s = ANY({{PARAM}})", fieldId),
	})
}

func (f PgFilters) Len() int {
	return len(f)
}

func (f PgFilters) SQL(queryParams *[]any) string {
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
