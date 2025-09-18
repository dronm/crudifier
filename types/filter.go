package types

type SQLFilterOperator string

const (
	SQL_FILTER_OPERATOR_E        SQLFilterOperator = "="
	SQL_FILTER_OPERATOR_L        SQLFilterOperator = "<"
	SQL_FILTER_OPERATOR_G        SQLFilterOperator = ">"
	SQL_FILTER_OPERATOR_LE       SQLFilterOperator = "<="
	SQL_FILTER_OPERATOR_GE       SQLFilterOperator = ">="
	SQL_FILTER_OPERATOR_LK       SQLFilterOperator = "LIKE"
	SQL_FILTER_OPERATOR_ILK      SQLFilterOperator = "ILIKE"
	SQL_FILTER_OPERATOR_NE       SQLFilterOperator = "<>"
	SQL_FILTER_OPERATOR_I        SQLFilterOperator = "IS"
	SQL_FILTER_OPERATOR_IN       SQLFilterOperator = "IS NOT"
	SQL_FILTER_OPERATOR_INCL     SQLFilterOperator = "IN"
	SQL_FILTER_OPERATOR_ANY      SQLFilterOperator = "ANY"
	SQL_FILTER_OPERATOR_HAS      SQLFilterOperator = "ANY"
	SQL_FILTER_OPERATOR_OVERLAP  SQLFilterOperator = "&&"
	SQL_FILTER_OPERATOR_CONTAINS SQLFilterOperator = "@>"
	SQL_FILTER_OPERATOR_TS       SQLFilterOperator = "@@"
)

type FilterJoin string

const (
	SQL_FILTER_JOIN_AND FilterJoin = "AND"
	SQL_FILTER_JOIN_OR  FilterJoin = "OR"
)

type DbFilter interface {
	FieldID() string
	Value() any
	Operator() SQLFilterOperator
	Expression() string // validated,sanatized expression
	Join() FilterJoin
	FieldPref() string
	SQL() string
}
