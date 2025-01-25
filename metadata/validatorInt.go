package metadata

import (
	"fmt"
	"reflect"
)

type FieldIntMetadata struct {
	FieldMetadata
	maxValue *int64
	minValue *int64
}

func NewFieldIntMedata(modelFieldId, id string) *FieldIntMetadata {
	return &FieldIntMetadata{FieldMetadata: FieldMetadata{modelId: modelFieldId, id: id, dataType: FIELD_TYPE_INT}}
}

func (f FieldIntMetadata) MaxValue() *int64 {
	return f.maxValue
}

func (f FieldIntMetadata) MinValue() *int64 {
	return f.minValue
}

type ModelFieldInt interface {
	GetValue() int64
	IsSet() bool
	IsNull() bool
}

func (f FieldIntMetadata) Validate(field reflect.Value) (bool, error) {
	var val int64
	modelField, ok := field.Interface().(ModelFieldInt)
	if ok {
		val = modelField.GetValue()
		if !modelField.IsSet() || modelField.IsNull() {
			return false, nil
		}

	} else {
		//standart type: int...
		switch field.Type().Kind() {
		case reflect.Int64:
			val, ok = field.Interface().(int64)
			if !ok {
				return true, fmt.Errorf(ER_VAL_CAST, f.ModelID(), "int64")
			}
		case reflect.Int32:
			val32, ok := field.Interface().(int32)
			if !ok {
				return true, fmt.Errorf(ER_VAL_CAST, f.ModelID(), "int32")
			}
			val = int64(val32)
		case reflect.Int16:
			val16, ok := field.Interface().(int16)
			if !ok {
				return true, fmt.Errorf(ER_VAL_CAST, f.ModelID(), "int16")
			}
			val = int64(val16)
		case reflect.Int8:
			val8, ok := field.Interface().(int8)
			if !ok {
				return true, fmt.Errorf(ER_VAL_CAST, f.ModelID(), "int8")
			}
			val = int64(val8)
		case reflect.Int:
			val0, ok := field.Interface().(int)
			if !ok {
				return true, fmt.Errorf(ER_VAL_CAST, f.ModelID(), "int")
			}
			val = int64(val0)
		case reflect.Float64:
			valFl, ok := field.Interface().(float64)
			if !ok {
				return true, fmt.Errorf(ER_VAL_CAST, f.ModelID(), "float64")
			}
			val = int64(valFl)
		default:
			return true, fmt.Errorf(ER_VAL_CAST, f.ModelID(), "int64")
		}
	}

	if f.MinValue() != nil && val < *f.MinValue() {
		return true, fmt.Errorf(ER_VAL_TOO_SMALL, f.Descr())
	}
	if f.MaxValue() != nil && val > *f.MaxValue() {
		return true, fmt.Errorf(ER_VAL_TOO_BIG, f.Descr())
	}

	return true, nil
}
