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

func castInt(fieldId string, field reflect.Value) (int64, error) {
	var val int64
	var ok bool

	switch field.Type().Kind() {
	case reflect.Int64:
		val, ok = field.Interface().(int64)
		if !ok {
			return 0, fmt.Errorf(ER_VAL_CAST, fieldId, "int64")
		}
	case reflect.Int32:
		val32, ok := field.Interface().(int32)
		if !ok {
			return 0, fmt.Errorf(ER_VAL_CAST, fieldId, "int32")
		}
		val = int64(val32)
	case reflect.Int16:
		val16, ok := field.Interface().(int16)
		if !ok {
			return 0, fmt.Errorf(ER_VAL_CAST, fieldId, "int16")
		}
		val = int64(val16)
	case reflect.Int8:
		val8, ok := field.Interface().(int8)
		if !ok {
			return 0, fmt.Errorf(ER_VAL_CAST, fieldId, "int8")
		}
		val = int64(val8)
	case reflect.Int:
		val0, ok := field.Interface().(int)
		if !ok {
			return 0, fmt.Errorf(ER_VAL_CAST, fieldId, "int")
		}
		val = int64(val0)
	case reflect.Float64:
		valFl, ok := field.Interface().(float64)
		if !ok {
			return 0, fmt.Errorf(ER_VAL_CAST, fieldId, "float64")
		}
		val = int64(valFl)
	default:
		return 0, fmt.Errorf(ER_VAL_CAST, fieldId, "int64")
	}

	return val, nil
}

func (f FieldIntMetadata) Validate(field reflect.Value) (bool, error) {
	var val int64
	modelField, ok := field.Interface().(ModelFieldInt)
	if ok {
		if !modelField.IsSet() {
			return false, nil
		} else if modelField.IsNull() {
			return true, nil
		}
		val = modelField.GetValue()

	} else if field.Kind() == reflect.Ptr && field.IsNil() {
		//standart type, nil pointer
		return true, nil

	} else if field.Kind() == reflect.Ptr {
		var err error

		elem := field.Elem()
		if !elem.IsValid() {
			return true, fmt.Errorf(ER_VAL_CAST, f.ModelID(), "int64")
		}

		val, err = castInt(f.ModelID(), elem)
		if err != nil {
			return true, err
		}

	} else {
		//standart type: int...
		var err error
		val, err = castInt(f.ModelID(), field)
		if err != nil {
			return true, err
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
