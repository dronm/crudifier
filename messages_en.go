package crudifier

const (
	ER_NOT_FIELD_IN_MD       = "%s() failed: field %s not found in metadata"
	ER_NOT_KEYS              = "%s() failed: keys not found"
	ER_AGG_FIELD_NOT_DEFINED = "PrepareFetchModelCollection() failed: aggregation field not defined for index %d"
	ER_AGG_FIELD_NO_FUNC     = "PrepareFetchModelCollection() failed: aggregation function not defined for field: %s"
	ER_FILTER_NOT_INIT       = "%s() failed: Filter should be initialized"
	ER_SORTER_NOT_INIT       = "%s() failed: Sorter should be initialized"
	ER_LIMIT_NOT_INIT        = "%s() failed: Limit should be initialized"
)
