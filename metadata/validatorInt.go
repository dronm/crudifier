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

func NewFieldIntMedata(modelFieldID, id string) *FieldIntMetadata {
	return &FieldIntMetadata{FieldMetadata: FieldMetadata{modelId: modelFieldID, id: id, dataType: FIELD_TYPE_INT}}
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

func castInt(fieldID string, field reflect.Value) (int64, error) {
	var val int64
	var ok bool

	// unwrap if it's interface{} holding a value
	if field.Kind() == reflect.Interface {
		field = field.Elem()
	}

	switch field.Type().Kind() {
	case reflect.Int64:
		val, ok = field.Interface().(int64)
		if !ok {
			return 0, fmt.Errorf(ER_VAL_CAST, fieldID, "int64")
		}
	case reflect.Int32:
		val32, ok := field.Interface().(int32)
		if !ok {
			return 0, fmt.Errorf(ER_VAL_CAST, fieldID, "int32")
		}
		val = int64(val32)
	case reflect.Int16:
		val16, ok := field.Interface().(int16)
		if !ok {
			return 0, fmt.Errorf(ER_VAL_CAST, fieldID, "int16")
		}
		val = int64(val16)
	case reflect.Int8:
		val8, ok := field.Interface().(int8)
		if !ok {
			return 0, fmt.Errorf(ER_VAL_CAST, fieldID, "int8")
		}
		val = int64(val8)
	case reflect.Int:
		val0, ok := field.Interface().(int)
		if !ok {
			return 0, fmt.Errorf(ER_VAL_CAST, fieldID, "int")
		}
		val = int64(val0)
	case reflect.Float64:
		valFl, ok := field.Interface().(float64)
		if !ok {
			return 0, fmt.Errorf(ER_VAL_CAST, fieldID, "float64")
		}
		val = int64(valFl)
	default:
		return 0, fmt.Errorf(ER_VAL_CAST, fieldID, "int64")
	}

	return val, nil
}

func (f FieldIntMetadata) ValidateSlice(fieldID string, field reflect.Value) error {
	if field.Kind() != reflect.Slice {
		return fmt.Errorf("field %s is not a slice", fieldID)
	}

	for i := 0; i < field.Len(); i++ {
		elem := field.Index(i)
		val, err := castInt(fmt.Sprintf("%s[%d]", fieldID, i), elem)
		if err != nil {
			return err
		}
		if err := f.CheckValue(field, val); err != nil {
			return err
		}
	}

	return nil
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
		// standart type, nil pointer
		return true, nil
	} else if field.Kind() == reflect.Ptr {
		var err error

		elem := field.Elem()
		if !elem.IsValid() {
			return true, fmt.Errorf(ER_VAL_CAST, f.ModelID(), "int64")
		}

		if field.Type().Kind() == reflect.Slice {
			// check every element
			if err := f.ValidateSlice(f.ModelID(), field); err != nil {
				return true, err
			}
			return true, nil
		}

		val, err = castInt(f.ModelID(), elem)
		if err != nil {
			return true, err
		}

	} else {
		// standart type: int...
		if field.Type().Kind() == reflect.Slice {
			// check every element
			if err := f.ValidateSlice(f.ModelID(), field); err != nil {
				return true, err
			}
			return true, nil
		}

		var err error
		val, err = castInt(f.ModelID(), field)
		if err != nil {
			return true, err
		}
	}

	return true, f.CheckValue(field, val)
}

// CheckValue does the actual validation of the given int64 value.
func (f FieldIntMetadata) CheckValue(field reflect.Value, val int64) error {
	if f.MinValue() != nil && val < *f.MinValue() {
		return fmt.Errorf(ER_VAL_TOO_SMALL, f.Descr())
	}
	if f.MaxValue() != nil && val > *f.MaxValue() {
		return fmt.Errorf(ER_VAL_TOO_BIG, f.Descr())
	}

	return nil
}
