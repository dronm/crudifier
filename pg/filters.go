package pg

import (
	"fmt"
	"strings"

	"github.com/dronm/crudifier/types"
)

type PgFilters []PgFilter

func (f *PgFilters) Add(fieldID string,
	value any, operator types.SQLFilterOperator,
	join types.FilterJoin) {

	*f = append(*f, PgFilter{fieldID: fieldID, value: value, join: join, operator: operator})
}

func (f *PgFilters) AddFullTextSearch(fieldID string, value any, join types.FilterJoin) {
	*f = append(*f, PgFilter{
		value: value,
		join: join,
		expression: fmt.Sprintf("%s @@ to_tsquery('russian', {{PARAM}})", fieldID),
	})
}

func (f *PgFilters) AddArrayInclude(fieldID string, value any, join types.FilterJoin) {
	*f = append(*f, PgFilter{
		value: value,
		join: join,
		expression: fmt.Sprintf("%s = ANY({{PARAM}})", fieldID),
	})
}

func (f *PgFilters) AddColumnArrayInclude(fieldID string, value any, join types.FilterJoin) {
	*f = append(*f, PgFilter{
		value: value,
		join: join,
		expression: fmt.Sprintf("{{PARAM}} = ANY(%s)", fieldID),
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
