package pg

import (
	"fmt"
	"strings"

	"github.com/dronm/crudifier/types"
)

type PgSelect struct {
	model          types.DbAggModel
	filter         *PgFilters
	sorter         *PgSorters
	limit          *PgLimit
	fieldIds       []string
	fieldValues    []interface{}
	aggFields      []string
	aggFieldValues []interface{}
}

func NewPgSelect(model types.DbAggModel, filter *PgFilters, sorter *PgSorters, limit *PgLimit) *PgSelect {
	return &PgSelect{model: model,
		filter: filter,
		sorter: sorter,
		limit:  limit,
	}
}

func (s PgSelect) Model() types.DbAggModel {
	return s.model
}

func (s PgSelect) Filter() types.DbFilters {
	return s.filter
}

func (s *PgSelect) SetFilter(f types.DbFilters) error {
	filters, ok := f.(*PgFilters)
	if !ok {
		return fmt.Errorf("could not assert to *PgFilters")
	}
	s.filter = filters
	return nil
}

func (s PgSelect) Limit() types.DbLimit {
	return s.limit
}

func (s PgSelect) Sorter() types.DbSorters {
	return s.sorter
}

func (s PgSelect) FieldValues() []interface{} {
	return s.fieldValues
}

func (s *PgSelect) AddField(id string, val interface{}) {
	s.fieldIds = append(s.fieldIds, id)
	s.fieldValues = append(s.fieldValues, val)
}

// AddAggField adds aggregate function, fn is the function,
// val is the value for scaning result.
func (s *PgSelect) AddAggField(fn string, val interface{}) {
	s.aggFields = append(s.aggFields, fn)
	s.aggFieldValues = append(s.aggFieldValues, val)
}

func (s PgSelect) SQL(queryParams *[]interface{}) string {
	var filterSQL string
	if s.filter != nil {
		filterSQL = s.filter.SQL(queryParams)
	}
	var sorterSQL string
	if s.sorter != nil {
		sorterSQL = s.sorter.SQL()
	}
	var limitSQL string
	if s.limit != nil {
		limitSQL = s.limit.SQL()
	}
	return fmt.Sprintf("SELECT %s FROM %s%s%s%s",
		strings.Join(s.fieldIds, ","),
		s.model.Relation(),
		filterSQL,
		sorterSQL,
		limitSQL,
	)
}

// CollectionSQL returns two queries: collecion query and aggregation query.
func (s PgSelect) CollectionSQL(queryParams *[]interface{}) (string, string) {
	var filterSQL string
	if s.filter != nil {
		filterSQL = s.filter.SQL(queryParams)
	}
	var sorterSQL string
	if s.sorter != nil {
		sorterSQL = s.sorter.SQL()
	}
	var limitSQL string
	if s.limit != nil {
		limitSQL = s.limit.SQL()
	}

	totQuery := ""
	if len(s.aggFields) > 0 {
		totQuery = fmt.Sprintf("SELECT %s FROM %s%s",
			strings.Join(s.aggFields, ","),
			s.model.Relation(),
			filterSQL,
		)
	}

	return fmt.Sprintf("SELECT %s FROM %s%s%s%s",
			strings.Join(s.fieldIds, ","),
			s.model.Relation(),
			filterSQL,
			sorterSQL,
			limitSQL,
		),
		totQuery
}
