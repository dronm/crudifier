package metadata

import (
	"fmt"
	"math"
	"reflect"
)

type FieldFloatMetadata struct {
	FieldMetadata
	maxValue  *float64
	minValue  *float64
	precision int64
}

func NewFieldFloatMedata(modelFieldId, id string) *FieldFloatMetadata {
	return &FieldFloatMetadata{FieldMetadata: FieldMetadata{modelId: modelFieldId, id: id, dataType: FIELD_TYPE_FLOAT}}
}

func (f FieldFloatMetadata) MaxValue() *float64 {
	return f.maxValue
}

func (f FieldFloatMetadata) MinValue() *float64 {
	return f.minValue
}

func (f FieldFloatMetadata) Precision() int64 {
	return f.precision
}

type ModelFieldFloat interface {
	GetValue() float64
	IsSet() bool
	IsNull() bool
}

func castFloat(fieldId string, field reflect.Value) (float64, error) {
	var val float64
	var ok bool

	switch field.Type().Kind() {
	case reflect.Float64:
		val, ok = field.Interface().(float64)
		if !ok {
			return 0, fmt.Errorf(ER_VAL_CAST, fieldId, "float64")
		}
	case reflect.Float32:
		val32, ok := field.Interface().(float32)
		if !ok {
			return 0, fmt.Errorf(ER_VAL_CAST, fieldId, "float32")
		}
		val = float64(val32)
	default:
		return 0, fmt.Errorf(ER_VAL_CAST, fieldId, "float64")
	}

	return val, nil
}

func (f FieldFloatMetadata) Validate(field reflect.Value) (bool, error) {
	var val float64
	modelField, ok := field.Interface().(ModelFieldFloat)
	if ok {
		val = modelField.GetValue()
		if !modelField.IsSet() {
			return false, nil

		} else if modelField.IsNull() {
			return true, nil
		}

	} else if field.Type().Kind() == reflect.Ptr && field.IsNil() {
		return true, nil

	} else if field.Type().Kind() == reflect.Ptr {
		var err error

		elem := field.Elem()
		if !elem.IsValid() {
			return true, fmt.Errorf(ER_VAL_CAST, f.ModelID(), "float64")
		}

		val, err = castFloat(f.ModelID(), elem)
		if err != nil {
			return true, err
		}

	} else {
		var err error
		val, err = castFloat(f.ModelID(), field)
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
	prec := f.Precision()
	if prec > 0 {
		factor := math.Pow(10, float64(prec))
		roundedValue := math.Round(val*factor) / factor
		if !(math.Abs(val-roundedValue) < 1e-9) {
			return true, fmt.Errorf(ER_VAL_PRECISION, f.Descr())
		}
	}

	return true, nil
}
