package fields

import (
	"database/sql/driver"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type FieldFloat struct {
	value   float64
	isSet   bool
	notNull bool
}

func NewFieldFloat(value float64, isSet, isNull bool) FieldFloat {
	return FieldFloat{value: value, isSet: isSet, notNull: !isNull}
}

func NewFieldFloatVal(value float64) FieldFloat {
	return FieldFloat{value: value, isSet: true, notNull: true}
}

func (f FieldFloat) GetValue() float64 {
	return f.value
}

func (f *FieldFloat) SetValue(v float64) {
	f.value = v
	f.isSet = true
	f.notNull = true
}

func (f *FieldFloat) UnsetValue() {
	f.value = 0 //default
	f.isSet = true
	f.notNull = false
}

func (f FieldFloat) IsSet() bool {
	return f.isSet
}

func (f FieldFloat) IsNull() bool {
	return !f.notNull
}

func (v *FieldFloat) UnmarshalJSON(data []byte) error {
	v.isSet = true
	v.notNull = true

	if ValIsNull(data) {
		v.notNull = false
		return nil
	}
	// dataStr := strings.Replace(string(data), ",", ".", 1)
	tmp, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return err
	}

	v.value = tmp

	return nil
}

func (v FieldFloat) String() string {
	return strconv.FormatFloat(v.value, 'f', -1, 64)
}

func (v *FieldFloat) MarshalJSON() ([]byte, error) {
	if !v.notNull {
		return []byte(jsonNull), nil

	} else {
		return []byte(v.String()), nil
	}
}

// driver.Scanner, driver.Valuer interfaces
func (v *FieldFloat) Scan(value any) error {
	v.isSet = true
	v.notNull = true

	if value == nil {
		v.notNull = false
		return nil
	}

	switch val := value.(type) {
	case float64:
		v.value = val

	case float32:
		v.value = float64(val)

	case int64:
		v.value = float64(val)

	case string:
		//0e0=0 1035e-2=10,35
		valStr := string(val)
		if isNan := strings.Index(strings.ToLower(valStr), "nan"); isNan >= 0 {
			v.notNull = false

		} else if expPos := strings.Index(valStr, "e"); expPos == -1 {
			//no exponent part
			var err error
			v.value, err = strconv.ParseFloat(valStr, 64)
			if err != nil {
				return err
			}

		} else {
			num, err := strconv.ParseInt(valStr[:expPos], 10, 64)
			if err != nil {
				return err
			}
			exp, err := strconv.ParseInt(valStr[expPos+1:], 10, 64)
			if err != nil {
				return err
			}
			v.value = float64(num) * math.Pow(10.0, float64(exp))
		}
	default:
		return fmt.Errorf(ErrUnsupportedType, "float", value)
	}

	return nil
}

func (v FieldFloat) Value() (driver.Value, error) {
	if !v.notNull {
		return driver.Value(nil), nil
	}
	return driver.Value(v.value), nil
}
