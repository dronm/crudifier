package fields

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const FormatDateTime = "2006-01-02T15:04:05"

type FieldDateTime struct {
	value   time.Time
	isSet   bool
	notNull bool
}

func NewFieldDateTime(value time.Time, isSet, isNull bool) FieldDateTime {
	return FieldDateTime{value: value, isSet: isSet, notNull: !isNull}
}

func NewFieldDateTimeVal(value time.Time) FieldDateTime {
	return FieldDateTime{value: value, isSet: true, notNull: true}
}

func (f FieldDateTime) GetValue() time.Time {
	return f.value
}

func (f *FieldDateTime) SetValue(v time.Time) {
	f.value = v
	f.isSet = true
	f.notNull = true
}

func (f *FieldDateTime) UnsetValue() {
	f.value = time.Time{} //default
	f.isSet = true
	f.notNull = false
}
func (f FieldDateTime) IsSet() bool {
	return f.isSet
}

func (f FieldDateTime) IsNull() bool {
	return !f.notNull
}

func (v *FieldDateTime) UnmarshalJSON(data []byte) error {
	v.isSet = true
	v.notNull = true
	v.value = time.Time{}

	if ValIsNull(data) {
		v.notNull = false
		return nil
	}

	tmp, err := time.Parse(FormatDateTime, RemoveQuotes(data))
	if err != nil {
		return err
	}
	v.value = tmp

	return nil
}

func (v *FieldDateTime) MarshalJSON() ([]byte, error) {
	if !v.notNull {
		return []byte(jsonNull), nil

	} else {
		return []byte(fmt.Sprintf(`"%s"`, v.value.Format(FormatDateTime))), nil
	}
}

func (v FieldDateTime) String() string {
	if !v.notNull {
		return ""
	}
	return v.value.Format(FormatDateTime)
}

// driver.Scanner, driver.Valuer interfaces
func (v *FieldDateTime) Scan(value any) error {
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
		tmp, err := time.Parse(FormatDateTime, val)
		if err != nil {
			return err
		}
		v.value = tmp
	default:
		return fmt.Errorf(ErrUnsupportedType, "datetime", value)
	}

	return nil
}

func (v FieldDateTime) Value() (driver.Value, error) {
	if !v.notNull {
		return driver.Value(nil), nil
	}
	return driver.Value(v.value), nil
}
