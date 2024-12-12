package pg

import (
	"testing"
)

func TestAssignerSQL(t *testing.T) {
	tests := []struct {
		assigner PgAssigner
		expSql   string
	}{
		{PgAssigner{fieldID: "f1", value: "abc"}, "f1 = $1"},
		{PgAssigner{fieldID: "f2", value: 124}, "f2 = $1"},
	}
	for _, test := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			params := make([]interface{}, 0)
			gotSql := test.assigner.SQL(&params)
			if test.expSql != gotSql {
				t.Fatalf("expected %s, got %s", test.expSql, gotSql)
			}
			if len(params) != 1 {
				t.Fatalf("params length expected to be 1, got %d", len(params))
			}
			switch test.assigner.value.(type) {
			case string:
				parVal, ok := params[0].(string)
				if !ok {
					t.Fatalf("param value is expected to be a string")
				}
				if test.assigner.value != parVal {
					t.Fatalf("parameter value expected %s, got %s", test.assigner.value, parVal)
				}
			case int:
				parVal, ok := params[0].(int)
				if !ok {
					t.Fatalf("param value is expected to be an int")
				}
				if test.assigner.value != parVal {
					t.Fatalf("parameter value expected %d, got %d", test.assigner.value, parVal)
				}
			}
		})
	}
}

func TestAssignersSQL(t *testing.T) {
	tests := []struct {
		assigners PgAssigners
		expSql    string
	}{
		{[]PgAssigner{{fieldID: "f1", value: "abc"}}, "f1 = $1"},
		{[]PgAssigner{{fieldID: "f1", value: "abc"},
			{fieldID: "f2", value: 123},
		}, "f1 = $1, f2 = $2"},
	}
	for _, test := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			params := make([]interface{}, 0)
			gotSql := test.assigners.SQL(&params)
			if test.expSql != gotSql {
				t.Fatalf("expected %s, got %s", test.expSql, gotSql)
			}
			parmCnt := len(test.assigners)
			if len(params) != parmCnt {
				t.Fatalf("params length expected to be %d, got %d", parmCnt, len(params))
			}
			for i, assigner := range test.assigners {
				switch assigner.value.(type) {
				case string:
					parVal, ok := params[i].(string)
					if !ok {
						t.Fatalf("param value is expected to be a string")
					}
					if assigner.value != parVal {
						t.Fatalf("parameter value expected %s, got %s", assigner.value, parVal)
					}
				case int:
					parVal, ok := params[i].(int)
					if !ok {
						t.Fatalf("param value is expected to be an int")
					}
					if assigner.value != parVal {
						t.Fatalf("parameter value expected %d, got %d", assigner.value, parVal)
					}
				}
			}
		})
	}
}
