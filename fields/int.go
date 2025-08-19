package fields

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

type FieldInt struct {
	value   int64
	isSet   bool
	notNull bool
}

func NewFieldInt(value int64, isSet, isNull bool) FieldInt {
	return FieldInt{value: value, isSet: isSet, notNull: !isNull}
}

func (f FieldInt) GetValue() int64 {
	return f.value
}

func (f *FieldInt) SetValue(v int64) {
	f.value = v
	f.isSet = true
	f.notNull = true
}

func (f *FieldInt) UnsetValue() {
	f.value = 0 //default
	f.isSet = true
	f.notNull = false
}

func (f FieldInt) IsSet() bool {
	return f.isSet
}

func (f FieldInt) IsNull() bool {
	return !f.notNull
}

func (v *FieldInt) UnmarshalJSON(data []byte) error {
	v.isSet = true
	v.notNull = true

	if ValIsNull(data) {
		v.notNull = false
		return nil
	}

	tmp, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}

	v.value = tmp

	return nil
}

func (v FieldInt) String() string {
	return strconv.FormatInt(v.value, 10)
}

func (v *FieldInt) MarshalJSON() ([]byte, error) {
	if !v.notNull {
		return []byte(JSON_NULL), nil

	} else {
		return []byte(v.String()), nil
	}
}

// driver.Scanner, driver.Valuer interfaces
func (v *FieldInt) Scan(value any) error {
	v.isSet = true
	v.notNull = true

	if value == nil {
		v.notNull = false
		return nil
	}

	if val, err := driver.Int32.ConvertValue(value); err == nil {
		var ok bool
		if v.value, ok = val.(int64); ok {
			return nil
		}
	}

	return fmt.Errorf(ER_UNSUPPORTED_TYPE, "int", value)
}

func (v FieldInt) Value() (driver.Value, error) {
	if !v.notNull {
		return driver.Value(nil), nil
	}
	return driver.Value(v.value), nil
}
