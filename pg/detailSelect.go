package pg

import (
	"fmt"
	"strings"

	"github.com/dronm/crudifier/types"
)

type PgDetailSelect struct {
	model          types.DbModel
	filter         *PgFilters
	fieldIds       []string
	fieldValues    []any
}

func NewPgDetailSelect(model types.DbModel, filter *PgFilters) *PgDetailSelect {
	return &PgDetailSelect{
		model: model,
		filter: filter,
	}
}

func (s PgDetailSelect) Model() types.DbModel {
	return s.model
}

func (s PgDetailSelect) Filter() types.DbFilters {
	return s.filter
}

func (s *PgDetailSelect) SetFilter(f types.DbFilters) error {
	filters, ok := f.(*PgFilters)
	if !ok {
		return fmt.Errorf("could not assert to *PgFilters")
	}
	s.filter = filters
	return nil
}

func (s PgDetailSelect) FieldValues() []any {
	return s.fieldValues
}

func (s *PgDetailSelect) AddField(id string, val any) {
	s.fieldIds = append(s.fieldIds, id)
	s.fieldValues = append(s.fieldValues, val)
}

func (s PgDetailSelect) SQL(queryParams *[]any) string {
	var filterSQL string
	if s.filter != nil {
		filterSQL = s.filter.SQL(queryParams)
	}
	return fmt.Sprintf("SELECT %s FROM %s%s",
		strings.Join(s.fieldIds, ","),
		s.model.Relation(),
		filterSQL,
	)
}

