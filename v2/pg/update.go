package pg

import (
	"fmt"

	"github.com/dronm/crudifier/v2/types"
)

type PgUpdate struct {
	model    types.DBModel
	assigner *PgAssigners
	filter   *PgFilters
	limit    *PgLimit
}

func NewPgUpdate(model types.DBModel) *PgUpdate {
	return &PgUpdate{model: model, filter: &PgFilters{}, assigner: &PgAssigners{}}
}

func (u PgUpdate) Model() types.DBModel {
	return u.model
}

func (u *PgUpdate) AddField(id string, value any) {
	u.assigner.Add(id, value)
}

func (u PgUpdate) AssignerLen() int {
	if u.assigner == nil {
		return 0
	}
	return u.assigner.Len()
}

func (u PgUpdate) Filter() types.DBFilters {
	return u.filter
}

func (s PgUpdate) SQL(queryParams *[]any) string {
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
