package pg

import (
	"fmt"
)

type PgAssigner struct {
	fieldID string
	value   any
}

func (a PgAssigner) FieldID() string {
	return a.fieldID
}

func (a PgAssigner) SQL(queryParams *[]any) string {
	parInd := 0
	if queryParams != nil {
		parInd = len(*queryParams)
	}
	parInd++

	//this is a patch for json fields.
	//if parameter is a struct, then marshal it to json,
	//do not rely on pgx.
	// t := reflect.TypeOf(a.value)
	// isStruct := t.Kind() == reflect.Struct || (t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct)
	// if isStruct {
	// 	fmt.Println("converting struct to JSON for field",a.fieldID)
	// 	*queryParams = append(*queryParams, a.value)
	// 	// b, _ := json.Marshal(a.value)
	// 	// *queryParams = append(*queryParams, b)
	// }else{
	// }
	*queryParams = append(*queryParams, a.value)

	return fmt.Sprintf("%s = $%d", a.fieldID, parInd)
}
