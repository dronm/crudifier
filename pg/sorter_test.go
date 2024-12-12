package pg

import (
	"testing"

	"github.com/dronm/crudifier/types"
)

func TestSorterSQL(t *testing.T) {

	tests := []struct {
		sorter PgSorter
		expSql string
	}{
		{PgSorter{fieldID: "f1", direct: types.SQL_SORT_ASC}, "f1 ASC"},
		{PgSorter{fieldID: "f1", direct: types.SQL_SORT_DESC}, "f1 DESC"},
	}
	for _, test := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			gotSql := test.sorter.SQL()
			if test.expSql != gotSql {
				t.Fatalf("expected %s, got %s", test.expSql, gotSql)
			}
		})
	}
}

func TestSortersSQL(t *testing.T) {
	tests := []struct {
		sorters PgSorters
		expSql  string
	}{
		{[]PgSorter{{fieldID: "f1", direct: types.SQL_SORT_ASC}},
			" ORDER BY f1 ASC",
		},

		{[]PgSorter{{fieldID: "f1", direct: types.SQL_SORT_ASC},
			{fieldID: "f2", direct: types.SQL_SORT_DESC},
		},
			" ORDER BY f1 ASC, f2 DESC",
		},
	}
	for _, test := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			gotSql := test.sorters.SQL()
			if test.expSql != gotSql {
				t.Fatalf("expected %s, got %s", test.expSql, gotSql)
			}
		})
	}
}
