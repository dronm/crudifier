package fields

import (
	"database/sql/driver"
	"fmt"
)

const (
	jsonTrue  = "true"
	jsonFalse = "false"
)

type FieldBool struct {
	value   bool
	isSet   bool
	notNull bool
}

func NewFieldBool(value bool, isSet, isNull bool) FieldBool {
	return FieldBool{value: value, isSet: isSet, notNull: !isNull}
}

func NewFieldBoolVal(value bool) FieldBool {
	return FieldBool{value: value, isSet: true, notNull: true}
}

func (f FieldBool) GetValue() bool {
	return f.value
}

func (f *FieldBool) SetValue(v bool) {
	f.value = v
	f.isSet = true
	f.notNull = true
}

func (f *FieldBool) UnsetValue() {
	f.value = false //default
	f.isSet = true
	f.notNull = false
}

func (f FieldBool) IsSet() bool {
	return f.isSet
}

func (f FieldBool) IsNull() bool {
	return !f.notNull
}

func (v *FieldBool) UnmarshalJSON(data []byte) error {
	v.isSet = true
	v.notNull = true
	v.value = false

	if ValIsNull(data) {
		v.notNull = false
		return nil
	}

	if string(data) == "true" {
		v.value = true
	}

	return nil
}

func (v *FieldBool) MarshalJSON() ([]byte, error) {
	if !v.notNull {
		return []byte(jsonNull), nil

	} else if v.value {
		return []byte(jsonTrue), nil
	} else {
		return []byte(jsonFalse), nil
	}
}

func (v FieldBool) String() string {
	if !v.notNull {
		return ""
	} else if v.value {
		return jsonTrue
	} else {
		return jsonFalse
	}
}

// driver.Scanner, driver.Valuer interfaces
func (v *FieldBool) Scan(value any) error {
	v.isSet = true
	v.notNull = true

	if value == nil {
		v.notNull = false
		return nil
	}

	val, err := driver.Bool.ConvertValue(value)
	if err != nil {
		return err
	}
	tmp, ok := val.(bool)
	if !ok {
		return fmt.Errorf(ErrUnsupportedType, "bool", value)
	}
	v.value = tmp

	return nil
}

func (v FieldBool) Value() (driver.Value, error) {
	if !v.notNull {
		return driver.Value(nil), nil
	}
	return driver.Value(v.value), nil
}
