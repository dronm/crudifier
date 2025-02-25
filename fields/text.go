package fields

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type FieldText struct {
	value   string
	isSet   bool
	notNull bool
}

func NewFieldText(value string, isSet, isNull bool) FieldText {
	return FieldText{value: value, isSet: isSet, notNull: !isNull}
}

func (f FieldText) GetValue() string {
	return f.value
}

func (f FieldText) IsSet() bool {
	return f.isSet
}

func (f FieldText) IsNull() bool {
	return !f.notNull
}

func (v *FieldText) UnmarshalJSON(data []byte) error {
	v.isSet = true
	v.notNull = true

	if ValIsNull(data) {
		v.notNull = false
		return nil
	}

	var tmp string
	if err := json.Unmarshal(data, &tmp); err != nil {
		return fmt.Errorf(ER_UNSUPPORTED_TYPE, "text", string(data))
	}

	v.value = tmp

	return nil
}

func (v FieldText) String() string {
	return v.value
}

func (v *FieldText) MarshalJSON() ([]byte, error) {
	if !v.notNull {
		return []byte(JSON_NULL), nil

	} else {
		return json.Marshal(v.value)
	}
}

// driver.Scanner, driver.Valuer interfaces
func (v *FieldText) Scan(value interface{}) error {
	v.isSet = true
	v.notNull = true

	if value == nil {
		v.notNull = false
		return nil
	}

	val, err := driver.String.ConvertValue(value)
	if err != nil {
		return err
	}

	if valStr, ok := val.(string); ok {
		v.value = valStr

	} else if valB, ok := val.([]byte); ok {
		v.value = string(valB)

	} else {
		return fmt.Errorf(ER_UNSUPPORTED_TYPE, "text", value)
	}

	return nil
}

func (v FieldText) Value() (driver.Value, error) {
	if !v.notNull {
		return driver.Value(nil), nil
	}
	return driver.Value(v.value), nil
}
