package pg

import (
	"fmt"

	"github.com/dronm/crudifier/types"
)

type PgFilter struct {
	fieldID    string
	value      interface{}
	operator   types.SQLFilterOperator
	expression string // validated,sanatized expression
	join       types.FilterJoin
	fieldPref  string
}

func NewPgFilter(fieldId string, value interface{}) *PgFilter {
	return &PgFilter{fieldID: fieldId, value: value,
		operator: types.SQL_FILTER_OPERATOR_E,
		join:     types.SQL_FILTER_JOIN_AND,
	}
}

func (f PgFilter) FieldID() string {
	return f.fieldID
}

func (f PgFilter) Value() interface{} {
	return f.value
}

func (f *PgFilter) SetValue(value interface{}) {
	f.value = value
}

func (f PgFilter) Operator() types.SQLFilterOperator {
	return f.operator
}

func (f *PgFilter) SetOperator(op types.SQLFilterOperator) {
	f.operator = op
}

func (f PgFilter) Expression() string {
	return f.expression
}

func (f *PgFilter) SetExpression(expr string) {
	f.expression = expr
}

func (f PgFilter) Join() types.FilterJoin {
	return f.join
}

func (f *PgFilter) SetJoin(j types.FilterJoin) {
	f.join = j
}

func (f PgFilter) FieldPref() string {
	return f.fieldPref
}

// SQL returns sql string, ready to be used in queries.
// Parameters are added to queryParams slice.
func (f PgFilter) SQL(queryParams *[]interface{}) string {
	var fieldId string
	if f.fieldPref != "" {
		fieldId = f.fieldPref + "."
	}
	fieldId += f.fieldID

	if f.expression != "" {
		return f.expression
	}
	if f.value == nil && (f.operator == types.SQL_FILTER_OPERATOR_I || f.operator == types.SQL_FILTER_OPERATOR_IN) {
		return fmt.Sprintf("%s %s NULL", fieldId, f.operator)
	}

	//default
	if f.operator == "" {
		f.operator = types.SQL_FILTER_OPERATOR_E
	}

	parInd := 0
	if queryParams != nil {
		parInd = len(*queryParams)
	}
	parInd++
	*queryParams = append(*queryParams, f.value)

	return fmt.Sprintf("%s %s $%d", fieldId, f.operator, parInd)
}
