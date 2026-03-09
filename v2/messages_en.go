package crudifier

const (
	ErrNoFieldInMD        = "%s() failed: field %s not found in metadata"
	ErrInvalidFieldExpr   = "%s() failed: invalid field expression %s"
	ErrNoInsertFields     = "%s() failed: no fields to insert"
	ErrNoUpdateFields     = "%s() failed: no fields to update"
	ErrNoKeys             = "%s() failed: keys not found"
	ErrAggFieldNotDefined = "PrepareFetchModelCollection() failed: aggregation field not defined for index %d"
	ErrAggFieldNoFnc      = "PrepareFetchModelCollection() failed: aggregation function not defined for field: %s"
	ErrFilterNotInit      = "%s() failed: Filter should be initialized"
	ErrSorterNotInit      = "%s() failed: Sorter should be initialized"
	ErrLimitNotInit       = "%s() failed: Limit should be initialized"
)
