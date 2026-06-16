package fields

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	FormatDateTimeTZ1 string = "2006-01-02T15:04:05.000-07"
	FormatDateTimeTZ2 string = "2006-01-02T15:04:05-07:00"
	FormatDateTimeTZ3 string = "2006-01-02T15:04:05Z07:00"
)

type FieldDateTimeTZ struct {
	value   time.Time
	isSet   bool
	notNull bool
}

func NewFieldDateTimeTZ(value time.Time, isSet, isNull bool) FieldDateTimeTZ {
	return FieldDateTimeTZ{value: value, isSet: isSet, notNull: !isNull}
}

func NewFieldDateTimeTZVal(value time.Time) FieldDateTimeTZ {
	return FieldDateTimeTZ{value: value, isSet: true, notNull: true}
}

func (f FieldDateTimeTZ) GetValue() time.Time {
	return f.value
}

func (f *FieldDateTimeTZ) SetValue(v time.Time) {
	f.value = v
	f.isSet = true
	f.notNull = true
}

func (f *FieldDateTimeTZ) UnsetValue() {
	f.value = time.Time{} //default
	f.isSet = true
	f.notNull = false
}

func (f FieldDateTimeTZ) IsSet() bool {
	return f.isSet
}

func (f FieldDateTimeTZ) IsNull() bool {
	return !f.notNull
}

func (v *FieldDateTimeTZ) UnmarshalJSON(data []byte) error {
	v.isSet = true
	v.notNull = true
	v.value = time.Time{}

	if ValIsNull(data) {
		v.notNull = false
		return nil
	}

	dataStr := RemoveQuotes(data)

	// layout := time.RFC3339Nano
	var layout string
	if strings.Contains(dataStr, "+") {
		layout = FormatDateTimeTZ2

	} else if strings.Contains(dataStr, "Z") {
		layout = FormatDateTimeTZ3

	} else {
		layout = FormatDateTimeTZ1
	}

	tmp, err := time.Parse(layout, dataStr)
	if err != nil {
		return fmt.Errorf("time.Parse(): %v, with layout: %v, dateStr: %s",err, layout, dataStr)
	}
	v.value = tmp

	return nil
}

func (v *FieldDateTimeTZ) MarshalJSON() ([]byte, error) {
	if !v.notNull {
		return []byte(jsonNull), nil

	} else {
		return json.Marshal(v.value)
	}
}

func (v FieldDateTimeTZ) String() string {
	if !v.notNull {
		return ""
	}
	return v.value.Format(FormatDateTimeTZ1)
}

func (v *FieldDateTimeTZ) Scan(value any) error {
	v.isSet = true
	v.notNull = true

	if value == nil {
		v.notNull = false
		return nil
	}

	switch val := value.(type) {
	case time.Time:
		v.value = val
	case string:
		tmp, err := time.Parse(FormatDateTimeTZ1, val)
		if err != nil {
			return err
		}
		v.value = tmp
	default:
		return fmt.Errorf(ErrUnsupportedType, "datetimetz", value)
	}

	return nil
}

func (v FieldDateTimeTZ) Value() (driver.Value, error) {
	if !v.notNull {
		return driver.Value(nil), nil
	}
	return driver.Value(v.value), nil
}
