package metadata

import (
	"fmt"
	"testing"
)

// "github.com/dronm/sqlutils/fields"
func TestValidateModel(t *testing.T) {

	Enums = map[string][]string{"sex": {"male", "female"}}

	tests := []struct {
		name      string
		expErr    bool
		expErrStr string
		model     interface{}
	}{
		{"enum", true,
			fmt.Sprintf(ER_VAL_VAL_LIST, "b"),
			&struct {
				B string `json:"b" required:"true" enum:"sex"`
			}{B: "neuter"},
		},
		{"value list", true,
			fmt.Sprintf(ER_VAL_VAL_LIST, "b"),
			&struct {
				B string `json:"b" required:"true" valList:"aaa@@bbb@@ccc"`
			}{B: "ddd"},
		},
		{"max value length", true,
			fmt.Sprintf(ER_VAL_LEN_TOO_LONG, "a"),
			&struct {
				A string `json:"a" required:"true" max:"5"`
			}{A: "0123456789"},
		},
		{"min value length", true,
			fmt.Sprintf(ER_VAL_LEN_TOO_SHORT, "b"),
			&struct {
				B string `json:"b" required:"true" min:"5"`
			}{B: "012"},
		},
		{"fix value length", true,
			fmt.Sprintf(ER_VAL_LEN_NOT_FIX, "b"),
			&struct {
				B string `json:"b" required:"true" fix:"5"`
			}{B: "0123456789"},
		},
		{"min value int", true,
			fmt.Sprintf(ER_VAL_TOO_SMALL, "b"),
			&struct {
				B int `json:"b" required:"true" min:"25"`
			}{B: 15},
		},
		{"min value *int", true,
			fmt.Sprintf(ER_VAL_TOO_SMALL, "b"),
			&struct {
				B *int `json:"b" required:"true" min:"25"`
			}{B: &[]int{15}[0]},
		},
		{"max value int", true,
			fmt.Sprintf(ER_VAL_TOO_BIG, "b"),
			&struct {
				B int `json:"b" required:"true" max:"25"`
			}{B: 75},
		},
		{"max value *int", true,
			fmt.Sprintf(ER_VAL_TOO_BIG, "b"),
			&struct {
				B *int `json:"b" required:"true" max:"25"`
			}{B: &[]int{75}[0]},
		},
		{"min value float32", true,
			fmt.Sprintf(ER_VAL_TOO_SMALL, "b"),
			&struct {
				B float32 `json:"b" required:"true" min:"5.75"`
			}{B: 3.14},
		},
		{"min value *float32", true,
			fmt.Sprintf(ER_VAL_TOO_SMALL, "b"),
			&struct {
				B *float32 `json:"b" required:"true" min:"5.75"`
			}{B: &[]float32{3.14}[0]},
		},
		{"max value float32", true,
			fmt.Sprintf(ER_VAL_TOO_BIG, "b"),
			&struct {
				B float32 `json:"b" required:"true" max:"3.14"`
			}{B: 3.15},
		},
		{"max value *float32", true,
			fmt.Sprintf(ER_VAL_TOO_BIG, "b"),
			&struct {
				B *float32 `json:"b" required:"true" max:"3.14"`
			}{B: &[]float32{3.15}[0]},
		},
		{"min value float64", true,
			fmt.Sprintf(ER_VAL_TOO_SMALL, "b"),
			&struct {
				B float64 `json:"b" required:"true" min:"5.75"`
			}{B: 3.14},
		},
		{"min value *float64", true,
			fmt.Sprintf(ER_VAL_TOO_SMALL, "b"),
			&struct {
				B *float64 `json:"b" required:"true" min:"5.75"`
			}{B: &[]float64{3.14}[0]},
		},
		{"max value float64", true,
			fmt.Sprintf(ER_VAL_TOO_BIG, "b"),
			&struct {
				B float64 `json:"b" required:"true" max:"3.14"`
			}{B: 3.15},
		},
		{"max value *float64", true,
			fmt.Sprintf(ER_VAL_TOO_BIG, "b"),
			&struct {
				B *float64 `json:"b" required:"true" max:"3.14"`
			}{B: &[]float64{3.15}[0]},
		},
		{"required string", true,
			fmt.Sprintf(ER_VAL_REQUIRED, "b"),
			&struct {
				B *string `json:"b" required:"true"`
			}{B: nil},
		},
		{"required int", true,
			fmt.Sprintf(ER_VAL_REQUIRED, "b"),
			&struct {
				B *int `json:"b" required:"true"`
			}{B: nil},
		},
		{"required float64", true,
			fmt.Sprintf(ER_VAL_REQUIRED, "b"),
			&struct {
				B *float64 `json:"b" required:"true"`
			}{B: nil},
		},
		{"required bool", true,
			fmt.Sprintf(ER_VAL_REQUIRED, "b"),
			&struct {
				B *bool `json:"b" required:"true"`
			}{B: nil},
		},
	}
	fieldTagName := "json"
	for _, ts := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			gotErr := ValidateModel(ts.model, fieldTagName)
			t.Logf("got error: %v", gotErr)
			if ts.expErr && gotErr == nil {
				t.Fatalf("test: %s, expected error %s, got none", ts.name, ts.expErrStr)
			}
			if !ts.expErr && gotErr != nil {
				t.Fatalf("test: %s, expected no error, got one: %v", ts.name, gotErr)
			}
			if ts.expErr && ts.expErrStr != gotErr.Error() {
				t.Fatalf("test: %s, expected error: %s, got: %s", ts.name, ts.expErrStr, gotErr.Error())
			}
		})
	}
}
