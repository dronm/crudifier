package pg

import (
	"strings"
)

type PgAssigners []PgAssigner

func (a PgAssigners) SQL(queryParams *[]interface{}) string {
	if len(a) == 0 {
		return ""
	} else if len(a) == 1 {
		return a[0].SQL(queryParams)
	}

	var sqlSt strings.Builder
	for i, assigner := range a {
		if i > 0 {
			sqlSt.WriteString(", ")
		}
		sqlSt.WriteString(assigner.SQL(queryParams))
	}

	return sqlSt.String()
}

func (a *PgAssigners) Add(fieldId string, value interface{}) {
	*a = append(*a, PgAssigner{fieldID: fieldId, value: value})
}

func (a PgAssigners) Len() int {
	return len(a)
}
