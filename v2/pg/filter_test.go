package pg

import (
	"testing"

	"github.com/dronm/crudifier/v2/types"
)

func TestFiltersSQL(t *testing.T) {
	tests := []struct {
		filters PgFilters
		expSql  string
	}{
		{[]PgFilter{{fieldID: "f1", value: "abc", operator: types.SQLFilterOperatorEq}}, " WHERE f1 = $1"},
		{[]PgFilter{{fieldID: "f1", value: "abc", operator: types.SQLFilterOperatorNotEq}}, " WHERE f1 <> $1"},

		{[]PgFilter{{fieldID: "f1", value: "abc", operator: types.SQLFilterOperatorEq},
			{fieldID: "f2", value: 125, operator: types.SQLFilterOperatorEq},
		}, " WHERE (f1 = $1) AND (f2 = $2)"},

		{[]PgFilter{{fieldID: "f1", value: "abc", operator: types.SQLFilterOperatorEq},
			{fieldID: "f2", value: 125, operator: types.SQLFilterOperatorEq, join: types.SQLFilterJoinOr},
		}, " WHERE (f1 = $1) OR (f2 = $2)"},
	}
	for _, test := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			params := make([]any, 0)
			gotSql := test.filters.SQL(&params)
			if test.expSql != gotSql {
				t.Fatalf("expected %s, got %s", test.expSql, gotSql)
			}
		})
	}
}

func TestFilterSQL(t *testing.T) {

	tests := []struct {
		filter PgFilter
		expSql string
	}{
		{PgFilter{fieldID: "f1", value: "abc", operator: types.SQLFilterOperatorEq}, "f1 = $1"},
		{PgFilter{fieldID: "f1", value: "abc", operator: types.SQLFilterOperatorLk}, "f1 LIKE $1"},
		{PgFilter{fieldID: "f1", value: "abc", operator: types.SQLFilterOperatorILk}, "f1 ILIKE $1"},
		{PgFilter{fieldID: "f1", value: "abc", operator: types.SQLFilterOperatorNotEq}, "f1 <> $1"},
		{PgFilter{fieldID: "f1", value: "abc", operator: types.SQLFilterOperatorLess}, "f1 < $1"},
		{PgFilter{fieldID: "f1", value: "abc", operator: types.SQLFilterOperatorLessEq}, "f1 <= $1"},
		{PgFilter{fieldID: "f1", value: "abc", operator: types.SQLFilterOperatorGr}, "f1 > $1"},
		{PgFilter{fieldID: "f1", value: "abc", operator: types.SQLFilterOperatorGrEq}, "f1 >= $1"},
		{PgFilter{fieldID: "f1", value: nil, operator: types.SQLFilterOperatorIn}, "f1 IS NOT NULL"},
		{PgFilter{fieldID: "f1", value: nil, operator: types.SQLFilterOperatorIs}, "f1 IS NULL"},
		{PgFilter{fieldID: "f1", value: []int{1, 2}, operator: types.SQLFilterOperatorAny}, "f1 = ANY($1)"},
	}
	for _, test := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			params := make([]any, 0)
			gotSql := test.filter.SQL(&params)
			if test.expSql != gotSql {
				t.Fatalf("expected %s, got %s", test.expSql, gotSql)
			}
		})
	}
}
