package types

type DbFilters interface {
	Add(fieldId string,
		value interface{}, operator SQLFilterOperator,
		join FilterJoin)
	Len() int
}
