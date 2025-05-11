package fields

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const FORMAT_DATE = "2006-01-02"

type FieldDate struct {
	value   time.Time
	isSet   bool
	notNull bool
}

func NewFieldDate(value time.Time, isSet, isNull bool) FieldDate {
	return FieldDate{value: value, isSet: isSet, notNull: !isNull}
}

func (f FieldDate) GetValue() time.Time {
	return f.value
}

func (f *FieldDate) SetValue(v time.Time) {
	f.value = v
	f.isSet = true
	f.notNull = true
}

func (f *FieldDate) UnsetValue() {
	f.value = time.Time{} //default
	f.isSet = true
	f.notNull = false
}

func (f FieldDate) IsSet() bool {
	return f.isSet
}

func (f FieldDate) IsNull() bool {
	return !f.notNull
}

func (v *FieldDate) UnmarshalJSON(data []byte) error {
	v.isSet = true
	v.notNull = true
	v.value = time.Time{}

	if ValIsNull(data) {
		v.notNull = false
		return nil
	}

	dataStr := RemoveQuotes(data)
	if len(dataStr) > len(FORMAT_DATE) {
		dataStr = dataStr[0:len(FORMAT_DATE)]
	}
	tmp, err := time.Parse(FORMAT_DATE, dataStr)
	if err != nil {
		return err
	}
	v.value = tmp

	return nil
}

func (v *FieldDate) MarshalJSON() ([]byte, error) {
	if !v.notNull {
		return []byte(JSON_NULL), nil

	} else {
		return []byte(fmt.Sprintf(`"%s"`, v.value.Format(FORMAT_DATE))), nil
	}
}

func (v FieldDate) String() string {
	if !v.notNull {
		return ""
	}
	return v.value.Format(FORMAT_DATE)
}

// driver.Scanner, driver.Valuer interfaces
func (v *FieldDate) Scan(value any) error {
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
		tmp, err := time.Parse(FORMAT_DATE, val)
		if err != nil {
			return err
		}
		v.value = tmp
	default:
		return fmt.Errorf(ER_UNSUPPORTED_TYPE, "date", value)
	}

	return nil
}

func (v FieldDate) Value() (driver.Value, error) {
	if !v.notNull {
		return driver.Value(nil), nil
	}
	return driver.Value(v.value), nil
}

