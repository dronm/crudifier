package types

type DbModel interface {
	Relation() string
}

type DbAggModel interface {
	DbModel
	CollectionAgg() any // structure with agg:"count(*)"
}

type PrepareModel interface {
	AddField(id string, val any)
}

// DbDetailSelecter is for detail model.
type DbDetailSelecter interface {
	Model() DbModel
	Filter() DbFilters
	SetFilter(DbFilters) error
	AddField(id string, val any)
}

// DbSelecter is for list model.
type DbSelecter interface {
	Model() DbAggModel
	Filter() DbFilters
	SetFilter(DbFilters) error
	Limit() DbLimit
	Sorter() DbSorters
	AddField(id string, val any)
	AddAggField(aggFn string, val any)
}

