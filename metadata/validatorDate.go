package metadata

import (
	"fmt"
	"reflect"
	"time"
)

type FieldDateMetadata struct {
	FieldMetadata
}

func NewFieldDateMedata(modelFieldId, id string, dataType FieldDataType) *FieldDateMetadata {
	return &FieldDateMetadata{FieldMetadata: FieldMetadata{modelId: modelFieldId, id: id, dataType: dataType}}
}

type ModelFieldDate interface {
	GetValue() time.Time
	IsSet() bool
	IsNull() bool
}

func (f FieldDateMetadata) Validate(field reflect.Value) (bool, error) {
	modelField, ok := field.Interface().(ModelFieldDate)
	if ok {
		if !modelField.IsSet() || modelField.IsNull() {
			return false, nil
		}
		return true, nil
	}

	timeField, ok := field.Interface().(time.Time)
	if !ok {
		return true, &IntervalServerError{ErrText: fmt.Sprintf(ER_VAL_CAST, f.ModelID(), "time.Time")}
	}
	return timeField == time.Time{}, nil
}
