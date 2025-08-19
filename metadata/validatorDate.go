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

	timeField, ok := field.Interface().(time.Time)
	if !ok {
		return true, fmt.Errorf(ER_VAL_CAST, f.ModelID(), "time.Time")
	}
	return timeField == time.Time{}, nil
}
