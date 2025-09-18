package types

type DbFilters interface {
	Add(fieldID string,
		value any, operator SQLFilterOperator,
		join FilterJoin)
	AddFullTextSearch(fieldID string, value any, join FilterJoin)
	AddArrayInclude(fieldID string, value any, join FilterJoin)
	AddColumnArrayInclude(fieldID string, value any, join FilterJoin)
	Len() int
}
