package metadata

import (
	"fmt"
	"reflect"
)

type FieldDataType byte

const (
	FIELD_TYPE_UNDEFINED FieldDataType = iota
	FIELD_TYPE_BOOL
	FIELD_TYPE_INT
	FIELD_TYPE_DATE
	FIELD_TYPE_DATETIME
	FIELD_TYPE_DATETIMETZ
	FIELD_TYPE_TIME
	FIELD_TYPE_FLOAT
	FIELD_TYPE_TEXT
	FIELD_TYPE_ARRAY  //slice
	FIELD_TYPE_OBJECT //struct
)

func ParseFieldType(fieldTypeName string) FieldDataType {
	switch fieldTypeName {
	case "FieldBool":
		return FIELD_TYPE_BOOL
	case "FieldDate":
		return FIELD_TYPE_DATE
	case "FieldDateTime":
		return FIELD_TYPE_DATETIME
	case "FieldDateTimeTZ":
		return FIELD_TYPE_DATETIMETZ
	case "FieldTime":
		return FIELD_TYPE_TIME
	case "FieldFloat", "float64", "float32":
		return FIELD_TYPE_FLOAT
	case "FieldInt", "int", "int0", "int8", "int16", "int32", "int64":
		return FIELD_TYPE_INT
	case "FieldText", "string":
		return FIELD_TYPE_TEXT
	}
	return FIELD_TYPE_UNDEFINED
}

type FieldMetadata struct {
	modelId    string //real structure ID
	id         string // database field id, json/xml tag
	alias      string
	required   bool
	dataType   FieldDataType
	primaryKey bool
	srvCalc    bool //server calculated field, return to client on insert
	//like auto inc for example
}

func (f FieldMetadata) ModelID() string {
	return f.modelId
}

func (f FieldMetadata) Alias() string {
	return f.alias
}

func (f *FieldMetadata) SetAlias(v string) {
	f.alias = v
}

func (f FieldMetadata) Required() bool {
	return f.required
}

func (f *FieldMetadata) SetRequired(v bool) {
	f.required = v
}

func (f FieldMetadata) PrimaryKey() bool {
	return f.primaryKey
}

func (f *FieldMetadata) SetPrimaryKey(v bool) {
	f.primaryKey = v
}

func (f FieldMetadata) SrvCalc() bool {
	return f.srvCalc
}

func (f *FieldMetadata) SetSrvCalc(v bool) {
	f.srvCalc = v
}
func (f FieldMetadata) ID() string {
	return f.id
}

func (f FieldMetadata) DataType() FieldDataType {
	return f.dataType
}

// Descr returns alias OR id
func (f FieldMetadata) Descr() string {
	if f.alias != "" {
		return f.alias
	}
	return f.id
}

type NullableField interface {
	// GetValue() interface{}
	IsSet() bool
	IsNull() bool
}

func (f FieldMetadata) ValidateRequired(field reflect.Value) error {
	modelField, ok := field.Interface().(NullableField)
	if !ok {
		//standart type, check for ptr
		if field.Kind() == reflect.Ptr && field.IsNil() && f.Required() {
			return fmt.Errorf(ER_VAL_REQUIRED, f.Descr())
		}
		return nil
	}
	if f.Required() && (!modelField.IsSet() || modelField.IsNull()) {
		return fmt.Errorf(ER_VAL_REQUIRED, f.Descr())
	}

	return nil
}

// Validate returns true if value is set.
// Value is set/unset can only be checked for pointers.
func (f FieldMetadata) Validate(field reflect.Value) (bool, error) {
	if field.Type().Kind() == reflect.Ptr && field.IsValid() {
		return !field.IsNil(), nil
	}
	return true, nil
}
