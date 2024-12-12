package pg

import (
	"fmt"

	"github.com/dronm/crudifier/types"
)

type PgDelete struct {
	model  types.DbModel
	filter PgFilters
}

func NewPgDelete(model types.DbModel, filter PgFilters) PgDelete {
	return PgDelete{model: model, filter: filter}
}

func (d PgDelete) Model() types.DbModel {
	return d.model
}

func (d PgDelete) Filter() PgFilters {
	return d.filter
}

func (s PgDelete) SQL(queryParams *[]interface{}) string {
	return fmt.Sprintf("DELETE FROM %s%s",
		s.model.Relation(),
		s.filter.SQL(queryParams),
	)
}
