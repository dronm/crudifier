package types

type DbModel interface {
	Relation() string
}

type DbAggModel interface {
	DbModel
	CollectionAgg() interface{} // structure with agg:"count(*)"
}

type DbSelecter interface {
	Model() DbAggModel
	Filter() DbFilters
	SetFilter(DbFilters) error
	Limit() DbLimit
	Sorter() DbSorters
	AddField(id string, val interface{})
	AddAggField(aggFn string, val interface{})
}
