package metadata

import (
	"fmt"
	"reflect"
)

type FieldBoolMetadata struct {
	FieldMetadata
}

func NewFieldBoolMedata(modelFieldID, id string) *FieldBoolMetadata {
	return &FieldBoolMetadata{FieldMetadata: FieldMetadata{modelId: modelFieldID, id: id, dataType: FieldTypeBool}}
}

type ModelFieldBool interface {
	GetValue() bool
	IsSet() bool
	IsNull() bool
}

func (f FieldBoolMetadata) Validate(field reflect.Value) (bool, error) {
	modelField, ok := field.Interface().(ModelFieldBool)
	if ok {
		//no farther validation
		return modelField.IsSet(), nil
	}

	var val bool
	if field.Kind() == reflect.Ptr && field.IsNil() {
		//standart type, nil pointer
		return true, nil
	} else if field.Kind() == reflect.Ptr {
		elem := field.Elem()
		if !elem.IsValid() {
			return true, fmt.Errorf(ErrValCast, f.ModelID(), "bool")
		}
		val, ok = elem.Interface().(bool)
		if !ok {
			return true, fmt.Errorf(ErrValCast, f.ModelID(), "bool")
		}
	} else {
		val, ok = field.Interface().(bool)
		if !ok {
			return true, fmt.Errorf(ErrValCast, f.ModelID(), "bool")
		}
	}

	return val, nil
}
