package metadata

import (
	"fmt"
	"reflect"
)

type FieldBoolMetadata struct {
	FieldMetadata
}

func NewFieldBoolMedata(modelFieldId, id string) *FieldBoolMetadata {
	return &FieldBoolMetadata{FieldMetadata: FieldMetadata{modelId: modelFieldId, id: id, dataType: FIELD_TYPE_BOOL}}
}

type ModelFieldBool interface {
	GetValue() bool
	IsSet() bool
	IsNull() bool
}

func (f FieldBoolMetadata) Validate(field reflect.Value) (bool, error) {
	modelField, ok := field.Interface().(ModelFieldBool)
	if ok {
		if !modelField.IsSet() || modelField.IsNull() {
			return false, nil
		}
		return true, nil
	}

	boolField, ok := field.Interface().(bool)
	if !ok {
		return true, fmt.Errorf(ER_VAL_CAST, f.ModelID(), "bool")
	}
	return boolField, nil
}
