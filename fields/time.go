package fields

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const FORMAT_TIME = "15:04:05"

type FieldTime struct {
	value   time.Time
	isSet   bool
	notNull bool
}

func NewFieldTime(value time.Time, isSet, isNull bool) FieldTime {
	return FieldTime{value: value, isSet: isSet, notNull: !isNull}
}

func (f FieldTime) GetValue() time.Time {
	return f.value
}

func (f FieldTime) IsSet() bool {
	return f.isSet
}

func (f FieldTime) IsNull() bool {
	return !f.notNull
}

func (v *FieldTime) UnmarshalJSON(data []byte) error {
	v.isSet = true
	v.notNull = true
	v.value = time.Time{}

	if ValIsNull(data) {
		v.notNull = false
		return nil
	}

	tmp, err := time.Parse(FORMAT_TIME, RemoveQuotes(data))
	if err != nil {
		return err
	}
	v.value = tmp

	return nil
}

func (v *FieldTime) MarshalJSON() ([]byte, error) {
	if !v.notNull {
		return []byte(JSON_NULL), nil

	} else {
		return []byte(fmt.Sprintf(`"%s"`, v.value.Format(FORMAT_TIME))), nil
	}
}

func (v FieldTime) String() string {
	if !v.notNull {
		return ""
	}
	return v.value.Format(FORMAT_TIME)
}

// driver.Scanner, driver.Valuer interfaces
func (v *FieldTime) Scan(value interface{}) error {
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
		tmp, err := time.Parse(FORMAT_TIME, val)
		if err != nil {
			return err
		}
		v.value = tmp
	default:
		return fmt.Errorf(ER_UNSUPPORTED_TYPE, "time", value)
	}

	return nil
}

func (v FieldTime) Value() (driver.Value, error) {
	if !v.notNull {
		return driver.Value(nil), nil
	}
	return driver.Value(v.value), nil
}
