package types

type DbFilters interface {
	Add(fieldId string,
		value any, operator SQLFilterOperator,
		join FilterJoin)
	AddFullTextSearch(fieldId string, value any, join FilterJoin)
	AddArrayInclude(fieldId string, value any, join FilterJoin)
	Len() int
}
