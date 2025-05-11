package pg

import (
	"fmt"
	"strings"

	"github.com/dronm/crudifier/types"
)

const EXPR_PARAM = "{{PARAM}}"

type PgFilter struct {
	fieldID    string
	value      any
	operator   types.SQLFilterOperator
	expression string // validated,sanatized expression
	join       types.FilterJoin
	fieldPref  string
}

func NewPgFilter(fieldId string, value any) *PgFilter {
	return &PgFilter{
		fieldID: fieldId, value: value,
		operator: types.SQL_FILTER_OPERATOR_E,
		join:     types.SQL_FILTER_JOIN_AND,
	}
}

func (f PgFilter) FieldID() string {
	return f.fieldID
}

func (f PgFilter) Value() any {
	return f.value
}

func (f *PgFilter) SetValue(value any) {
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
func (f PgFilter) SQL(queryParams *[]any) string {
	var fieldId string
	if f.fieldPref != "" {
		fieldId = f.fieldPref + "."
	}
	fieldId += f.fieldID

	if f.expression != "" {
		// expression can have a parameter: {{PARAM}}

		if strings.Contains(f.expression, EXPR_PARAM) {
			parInd := 0
			if queryParams != nil {
				parInd = len(*queryParams)
			}
			parInd++
			*queryParams = append(*queryParams, f.value)
			return strings.Replace(f.expression, EXPR_PARAM, fmt.Sprintf("$%d", parInd), 1)
		}

		// no parameter, pure expression
		return f.expression
	}

	if f.value == nil && (f.operator == types.SQL_FILTER_OPERATOR_I || f.operator == types.SQL_FILTER_OPERATOR_IN) {
		return fmt.Sprintf("%s %s NULL", fieldId, f.operator)
	}

	// default
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
