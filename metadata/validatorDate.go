package metadata

import (
	"fmt"
	"reflect"
	"time"
)

type FieldDateMetadata struct {
	FieldMetadata
}

func NewFieldDateMedata(modelFieldID, id string, dataType FieldDataType) *FieldDateMetadata {
	return &FieldDateMetadata{FieldMetadata: FieldMetadata{modelId: modelFieldID, id: id, dataType: dataType}}
}

type ModelFieldDate interface {
	GetValue() time.Time
	IsSet() bool
	IsNull() bool
}

func (f FieldDateMetadata) Validate(field reflect.Value) (bool, error) {
	modelField, ok := field.Interface().(ModelFieldDate)
	if ok {
		if !modelField.IsSet() {
			return false, nil
		}
		return true, nil
	}

	if field.Kind() == reflect.Ptr && field.IsNil() {
		return true, nil

	} else if field.Type().Kind() == reflect.Ptr {
		field = field.Elem()
		if !field.IsValid() {
			return true, fmt.Errorf(ER_VAL_CAST, f.ModelID(), "date")
		}
	}

	txtField, ok := field.Interface().(string)
	if ok {
		if f.DataType() == FIELD_TYPE_DATE && len(txtField) != 10 {
			return true, fmt.Errorf(ER_VAL_CAST, f.ModelID(), "string")

		}else if f.DataType() == FIELD_TYPE_TIME && len(txtField) != 5 {
			return true, fmt.Errorf(ER_VAL_CAST, f.ModelID(), "string")

		}else if f.DataType() == FIELD_TYPE_DATETIME && len(txtField) != 19 {
			return true, fmt.Errorf(ER_VAL_CAST, f.ModelID(), "string")

		}else if f.DataType() == FIELD_TYPE_DATETIMETZ && len(txtField) < 20 {
			return true, fmt.Errorf(ER_VAL_CAST, f.ModelID(), "string")
		}

		return true, nil
	}

	timeField, ok := field.Interface().(time.Time)
	if !ok {
		return true, fmt.Errorf(ER_VAL_CAST, f.ModelID(), "time.Time")
	}
	return timeField.Equal(time.Time{}), nil
}
