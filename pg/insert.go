package pg

import (
	"fmt"
	"strings"

	"github.com/dronm/crudifier/types"
)

type PgField struct {
	ID    string
	Value any
}

type PgInsert struct {
	model          types.DbModel
	values         []any
	fields         []PgField
	retFieldIds    []string
	retFieldValues []any
}

func NewPgInsert(model types.DbModel) *PgInsert {
	return &PgInsert{model: model}
}

func (s PgInsert) Model() types.DbModel {
	return s.model
}

func (s *PgInsert) AddRetField(id string, val any) {
	s.retFieldIds = append(s.retFieldIds, id)
	s.retFieldValues = append(s.retFieldValues, val)
}

func (s PgInsert) RetFieldIds() []string {
	return s.retFieldIds
}

func (s PgInsert) RetFieldValues() []any {
	return s.retFieldValues
}

func (s PgInsert) RetFields() map[string]any {
	res := make(map[string]any, len(s.retFieldIds))
	for i, f := range s.retFieldIds {
		res[f] = s.retFieldValues[i]
	}
	return res
}

func (s *PgInsert) AddField(fieldId string, val any) {
	s.fields = append(s.fields, PgField{ID: fieldId, Value: val})
}

func (s PgInsert) SQL(queryParams *[]any) string {
	//values to a string
	paramInd := len(*queryParams)
	var fieldIds strings.Builder
	var fieldVals strings.Builder
	for _, field := range s.fields {
		if fieldIds.Len() > 0 {
			fieldIds.WriteString(",")
			fieldVals.WriteString(",")
		}
		paramInd++
		fieldVals.WriteString(fmt.Sprintf("$%d", paramInd))
		fieldIds.WriteString(field.ID)
		*queryParams = append(*queryParams, field.Value)
	}

	retFields := ""
	if len(s.retFieldIds) > 0 {
		retFields = " RETURNING " + strings.Join(s.retFieldIds, ",")
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)%s",
		s.model.Relation(),
		fieldIds.String(),
		fieldVals.String(),
		retFields,
	)
}

// // InsertModel insert model data to database ans returns server init field values and
// // primary keys.
// func InsertModel(ctx context.Context, conn *pgx.Conn, model DbModel) (any, error) {
// 	dbInsert := NewPgInsert(model)
// 	if err := qclauses.PrepareInsertModel(dbInsert); err != nil {
// 		return nil, err
// 	}
//
// 	queryParams := make([]any, 0)
// 	query := dbInsert.SQL(&queryParams)
// 	if err := conn.QueryRow(ctx, query, queryParams...).Scan(dbInsert.retFieldValues...); err != nil {
// 		return nil, err
// 	}
//
// 	return dbInsert.retFieldValues, nil
// }
