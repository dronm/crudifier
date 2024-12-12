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

func (f FieldFloatMetadata) Validate(field reflect.Value) (bool, error) {
	var val float64
	modelField, ok := field.Interface().(ModelFieldFloat)
	if ok {
		val = modelField.GetValue()
		if !modelField.IsSet() || modelField.IsNull() {
			return false, nil
		}

	} else {
		//standart type: int...
		switch field.Type().Kind() {
		case reflect.Float64:
			val, ok = field.Interface().(float64)
			if !ok {
				return true, &IntervalServerError{ErrText: fmt.Sprintf(ER_VAL_CAST, f.ModelID(), "float64")}
			}
		case reflect.Float32:
			val32, ok := field.Interface().(float32)
			if !ok {
				return true, &IntervalServerError{ErrText: fmt.Sprintf(ER_VAL_CAST, f.ModelID(), "float32")}
			}
			val = float64(val32)
		default:
			return true, &IntervalServerError{ErrText: fmt.Sprintf(ER_VAL_CAST, f.ModelID(), "float64")}
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
