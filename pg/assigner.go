package pg

import (
	"fmt"
)

type PgAssigner struct {
	fieldID string
	value   interface{}
}

func (a PgAssigner) FieldID() string {
	return a.fieldID
}

func (a PgAssigner) SQL(queryParams *[]interface{}) string {
	parInd := 0
	if queryParams != nil {
		parInd = len(*queryParams)
	}
	parInd++
	*queryParams = append(*queryParams, a.value)

	return fmt.Sprintf("%s = $%d", a.fieldID, parInd)
}
