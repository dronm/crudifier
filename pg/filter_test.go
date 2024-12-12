package pg

import (
	"testing"

	"github.com/dronm/crudifier/types"
)

func TestFiltersSQL(t *testing.T) {
	tests := []struct {
		filters PgFilters
		expSql  string
	}{
		{[]PgFilter{{fieldID: "f1", value: "abc", operator: types.SQL_FILTER_OPERATOR_E}}, " WHERE f1 = $1"},
		{[]PgFilter{{fieldID: "f1", value: "abc", operator: types.SQL_FILTER_OPERATOR_NE}}, " WHERE f1 <> $1"},

		{[]PgFilter{{fieldID: "f1", value: "abc", operator: types.SQL_FILTER_OPERATOR_E},
			{fieldID: "f2", value: 125, operator: types.SQL_FILTER_OPERATOR_E},
		}, " WHERE (f1 = $1) AND (f2 = $2)"},

		{[]PgFilter{{fieldID: "f1", value: "abc", operator: types.SQL_FILTER_OPERATOR_E},
			{fieldID: "f2", value: 125, operator: types.SQL_FILTER_OPERATOR_E, join: types.SQL_FILTER_JOIN_OR},
		}, " WHERE (f1 = $1) OR (f2 = $2)"},
	}
	for _, test := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			params := make([]interface{}, 0)
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
		{PgFilter{fieldID: "f1", value: "abc", operator: types.SQL_FILTER_OPERATOR_E}, "f1 = $1"},
		{PgFilter{fieldID: "f1", value: "abc", operator: types.SQL_FILTER_OPERATOR_LK}, "f1 LIKE $1"},
		{PgFilter{fieldID: "f1", value: "abc", operator: types.SQL_FILTER_OPERATOR_ILK}, "f1 ILIKE $1"},
		{PgFilter{fieldID: "f1", value: "abc", operator: types.SQL_FILTER_OPERATOR_NE}, "f1 <> $1"},
		{PgFilter{fieldID: "f1", value: "abc", operator: types.SQL_FILTER_OPERATOR_L}, "f1 < $1"},
		{PgFilter{fieldID: "f1", value: "abc", operator: types.SQL_FILTER_OPERATOR_LE}, "f1 <= $1"},
		{PgFilter{fieldID: "f1", value: "abc", operator: types.SQL_FILTER_OPERATOR_G}, "f1 > $1"},
		{PgFilter{fieldID: "f1", value: "abc", operator: types.SQL_FILTER_OPERATOR_GE}, "f1 >= $1"},
		{PgFilter{fieldID: "f1", value: nil, operator: types.SQL_FILTER_OPERATOR_IN}, "f1 IS NOT NULL"},
		{PgFilter{fieldID: "f1", value: nil, operator: types.SQL_FILTER_OPERATOR_I}, "f1 IS NULL"},
	}
	for _, test := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			params := make([]interface{}, 0)
			gotSql := test.filter.SQL(&params)
			if test.expSql != gotSql {
				t.Fatalf("expected %s, got %s", test.expSql, gotSql)
			}
		})
	}
}
