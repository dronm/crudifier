package crudifier

import (
	"testing"
	"time"

	"github.com/dronm/crudifier/fields"
	"github.com/dronm/crudifier/pg"
)

type TableA struct {
	Field1 int    `json:"field1" primaryKey:"true" srvCalc:"true"`
	Field2 string `json:"field2" required:"true"`
}

func (t TableA) Relation() string {
	return "table_a"
}

func (t TableA) CollectionAgg() any {
	return nil
}

type TableB struct {
	ID              fields.FieldInt        `json:"id" primaryKey:"true" srvCalc:"true"` //auto inc field
	TextField       fields.FieldText       `json:"text_field" required:"true" max:"10" min:"1"`
	IntField        fields.FieldInt        `json:"int_field" required:"true" max:"10" min:"1"`
	FloatField      fields.FieldFloat      `json:"float_field" required:"true" max:"10.3" min:"1.5"`
	BoolField       fields.FieldBool       `json:"bool_field" required:"true"`
	DateField       fields.FieldDate       `json:"date_field" required:"true"`
	DateTimeField   fields.FieldDateTime   `json:"date_time_field" required:"true"`
	DateTimeTZField fields.FieldDateTimeTZ `json:"date_timetz_field" required:"true"`
}

func (t TableB) Relation() string {
	return "table_b"
}

func (t TableB) CollectionAgg() any {
	return &struct {
		TotCount fields.FieldInt `json:"tot_count" agg:"count(*)"`
	}{
		fields.NewFieldInt(0, true, false),
	}
}

type TableC struct {
	ID        fields.FieldInt  `json:"id" primaryKey:"true" srvCalc:"true"` //auto inc field
	TextField fields.FieldText `json:"text_field" required:"true"`
	IntField  fields.FieldInt  `json:"int_field" required:"true"`
}

func (t TableC) Relation() string {
	return "table_c"
}

func (t TableC) CollectionAgg() any {
	return &struct {
		TotCount fields.FieldInt `json:"tot_count" agg:"count(*)"`
	}{
		fields.NewFieldInt(0, true, false),
	}
}
func TestPrepareFetchCollection(t *testing.T) {
	tests := []struct {
		keyModel  any
		dbSelect  pg.PgSelect
		parms     CollectionParams
		expSql    string
		expAggSql string
		valid     bool
	}{
		{
			&struct {
				ID fields.FieldInt `json:"id" required:"true"`
			}{ID: fields.NewFieldInt(1, true, false)},
			*pg.NewPgSelect(&TableC{},
				&pg.PgFilters{},
				&pg.PgSorters{},
				&pg.PgLimit{},
			),
			CollectionParams{},
			"SELECT id,text_field,int_field FROM table_c",
			"SELECT count(*) FROM table_c",
			true,
		},
		{
			&struct {
				ID fields.FieldInt `json:"id" required:"true"`
			}{ID: fields.NewFieldInt(1, true, false)},
			*pg.NewPgSelect(&TableC{},
				&pg.PgFilters{},
				&pg.PgSorters{},
				&pg.PgLimit{},
			),
			CollectionParams{Count: 60,
				Sorter: []CollectionSorter{{"text_field", SORT_PAR_ASC}, {"int_field", SORT_PAR_DESC}},
				Filter: []CollectionFilter{{Join: FILTER_PAR_JOIN_AND,
					Fields: map[string]CollectionFilterField{"int_field": {FILTER_OPER_PAR_E, 1}},
				}}},
			"SELECT id,text_field,int_field FROM table_c WHERE int_field = $1 ORDER BY text_field ASC, int_field DESC LIMIT 60",
			"SELECT count(*) FROM table_c WHERE int_field = $1",
			true,
		},
	}
	for _, test := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			err := PrepareFetchModelCollection(&test.dbSelect, test.parms)
			if err != nil && test.valid {
				t.Fatalf("PrepareInsertModel() returned an error: %v, but it should not", err)

			} else if err == nil && !test.valid {
				t.Fatal("PrepareInsertModel() returned no error, but should return one")
			}

			params := make([]any, 0)
			gotSql, gotAggSql := test.dbSelect.CollectionSQL(&params)
			if test.expSql != gotSql {
				t.Fatalf("expected %s, got %s", test.expSql, gotSql)
			}
			if test.expAggSql != gotAggSql {
				t.Fatalf("expected aggregation sql: %s, got: %s", test.expAggSql, gotAggSql)
			}
			// t.Logf("sql: %s, agg: %s", gotSql, gotAggSql)
		})
	}
}

func TestPrepareFetch(t *testing.T) {
	tests := []struct {
		keyModel any
		dbSelect pg.PgDetailSelect
		expSql   string
		valid    bool
	}{
		{
			&struct {
				ID fields.FieldInt `json:"id" required:"true"`
			}{ID: fields.NewFieldInt(1, true, false)},
			*pg.NewPgDetailSelect(&TableB{},
				&pg.PgFilters{},
			),
			"",
			true,
		},
	}
	for _, test := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			err := PrepareFetchModel(test.keyModel, &test.dbSelect)
			if err != nil && test.valid {
				t.Fatalf("PrepareInsertModel() returned an error: %v, but it should not", err)

			} else if err == nil && !test.valid {
				t.Fatal("PrepareInsertModel() returned no error, but should return one")
			}

			params := make([]any, 0)
			gotSql := test.dbSelect.SQL(&params)
			// if test.expSql != gotSql {
			// 	t.Fatalf("expected %s, got %s", test.expSql, gotSql)
			// }
			t.Logf("sql: %s", gotSql)
		})
	}
}

func TestPrepareUpdate(t *testing.T) {
	tests := []struct {
		keyModel any
		update   pg.PgUpdate
		expSql   string
		valid    bool
	}{
		{
			&struct {
				ID fields.FieldInt `json:"id" required:"true"`
			}{ID: fields.NewFieldInt(1, true, false)},
			*pg.NewPgUpdate(&TableB{TextField: fields.NewFieldText("some text", true, false),
				IntField:   fields.NewFieldInt(1, true, false),
				FloatField: fields.NewFieldFloat(3.14, true, false),
				BoolField:  fields.NewFieldBool(false, true, false),
				DateField: fields.NewFieldDate(func() time.Time {
					t, _ := time.Parse(fields.FORMAT_DATE, "2024-12-25")
					return t
				}(), true, false),
				DateTimeField: fields.NewFieldDateTime(func() time.Time {
					t, _ := time.Parse(fields.FORMAT_DATE_TIME, "2024-12-25T07:30:00")
					return t
				}(), true, false),
				DateTimeTZField: fields.NewFieldDateTimeTZ(func() time.Time {
					t, _ := time.Parse(fields.FORMAT_DATE_TIME_TZ1, "2024-12-25T07:30:00.000-05")
					return t
				}(), true, false),
			}),
			"",
			true,
		},
	}
	for _, test := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			err := PrepareUpdateModel(test.keyModel, &test.update)
			if err != nil && test.valid {
				t.Fatalf("PrepareInsertModel() returned an error: %v, but it should not", err)

			} else if err == nil && !test.valid {
				t.Fatal("PrepareInsertModel() returned no error, but should return one")
			}

			params := make([]any, 0)
			gotSql := test.update.SQL(&params)
			// if test.expSql != gotSql {
			// 	t.Fatalf("expected %s, got %s", test.expSql, gotSql)
			// }
			t.Logf("sql: %s", gotSql)
		})
	}
}
func TestPrepareInsert(t *testing.T) {
	tests := []struct {
		insert pg.PgInsert
		expSql string
		valid  bool
	}{
		{*pg.NewPgInsert(&TableB{TextField: fields.NewFieldText("some text", true, false),
			IntField:   fields.NewFieldInt(1, true, false),
			FloatField: fields.NewFieldFloat(3.14, true, false),
			BoolField:  fields.NewFieldBool(false, true, false),
			DateField: fields.NewFieldDate(func() time.Time {
				t, _ := time.Parse(fields.FORMAT_DATE, "2024-12-25")
				return t
			}(), true, false),
			DateTimeField: fields.NewFieldDateTime(func() time.Time {
				t, _ := time.Parse(fields.FORMAT_DATE_TIME, "2024-12-25T07:30:00")
				return t
			}(), true, false),
			DateTimeTZField: fields.NewFieldDateTimeTZ(func() time.Time {
				t, _ := time.Parse(fields.FORMAT_DATE_TIME_TZ1, "2024-12-25T07:30:00.000-05")
				return t
			}(), true, false),
		}),
			"",
			true,
		},
	}
	for _, test := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			err := PrepareInsertModel(&test.insert)
			if err != nil && test.valid {
				t.Fatalf("PrepareInsertModel() returned an error: %v, but it should not", err)

			} else if err == nil && !test.valid {
				t.Fatal("PrepareInsertModel() returned no error, but should return one")
			}

			params := make([]any, 0)
			gotSql := test.insert.SQL(&params)
			// if test.expSql != gotSql {
			// 	t.Fatalf("expected %s, got %s", test.expSql, gotSql)
			// }
			t.Logf("sql: %s", gotSql)
		})
	}
}

func TestDeleteQuery(t *testing.T) {
	tests := []struct {
		delete pg.PgDelete
		expSql string
	}{
		{pg.NewPgDelete(TableA{}, nil), "DELETE FROM table_a"},
		{pg.NewPgDelete(TableA{}, pg.PgFilters{*pg.NewPgFilter("field1", 1)}), "DELETE FROM table_a WHERE field1 = $1"},
		{pg.NewPgDelete(TableB{},
			pg.PgFilters{*pg.NewPgFilter("id", 1),
				*pg.NewPgFilter("text_field", "some_val"),
			}),
			"DELETE FROM table_b WHERE (id = $1) AND (text_field = $2)",
		},
	}
	for _, test := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			params := make([]any, 0)
			gotSql := test.delete.SQL(&params)
			if test.expSql != gotSql {
				t.Fatalf("expected %s, got %s", test.expSql, gotSql)
			}
		})
	}
}
