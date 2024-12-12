package pg

import (
	"fmt"

	"github.com/dronm/crudifier/types"
)

type PgUpdate struct {
	model    types.DbModel
	assigner *PgAssigners
	filter   *PgFilters
	limit    *PgLimit
}

func NewPgUpdate(model types.DbModel) *PgUpdate {
	return &PgUpdate{model: model, filter: &PgFilters{}, assigner: &PgAssigners{}}
}

func (u PgUpdate) Model() types.DbModel {
	return u.model
}

func (u *PgUpdate) AddField(id string, value interface{}) {
	u.assigner.Add(id, value)
}

func (u PgUpdate) Filter() types.DbFilters {
	return u.filter
}

func (s PgUpdate) SQL(queryParams *[]interface{}) string {
	var assignerSQL string
	if s.assigner != nil {
		assignerSQL = s.assigner.SQL(queryParams)
	}
	var filterSQL string
	if s.filter != nil {
		filterSQL = s.filter.SQL(queryParams)
	}
	var limitSQL string
	if s.limit != nil {
		limitSQL = s.limit.SQL()
	}

	return fmt.Sprintf("UPDATE %s SET %s%s%s",
		s.model.Relation(),
		assignerSQL,
		filterSQL,
		limitSQL,
	)
}
