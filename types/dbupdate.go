package types

type DbUpdater interface {
	Model() DbModel
	Filter() DbFilters
	AddField(id string, value any)
}
