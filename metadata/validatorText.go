package metadata

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
)

type FieldTextMetadata struct {
	FieldMetadata
	maxLength *int64
	minLength *int64
	fixLength *int64
	regExp    string
	valList   []string
}

func NewFieldTextMedata(modelFieldId, id string) *FieldTextMetadata {
	return &FieldTextMetadata{FieldMetadata: FieldMetadata{modelId: modelFieldId, id: id, dataType: FIELD_TYPE_TEXT}}
}

func (f FieldTextMetadata) MaxLength() *int64 {
	return f.maxLength
}

func (f FieldTextMetadata) MinLength() *int64 {
	return f.minLength
}

func (f FieldTextMetadata) FixLength() *int64 {
	return f.fixLength
}

func (f FieldTextMetadata) RegExp() string {
	return f.regExp
}

func (f FieldTextMetadata) ValList() []string {
	return f.valList
}

type ModelFieldText interface {
	GetValue() string
	IsSet() bool
	IsNull() bool
}

// Validate validates a text field, returns flag
// indicating that the value is set (bool) and error
func (f FieldTextMetadata) Validate(field reflect.Value) (bool, error) {
	var val string
	textField, ok := field.Interface().(ModelFieldText)
	if ok {
		val = textField.GetValue()
		if !textField.IsSet() { //|| textField.IsNull()
			return false, nil
		} else if textField.IsNull() {
			return true, nil
		}

	} else if field.Kind() == reflect.Ptr && field.IsNil() {
		return true, nil

	} else {
		//standart type, need type assertion
		valIntf := field.Interface()
		val, ok = valIntf.(string)
		if !ok {
			return true, fmt.Errorf(ER_VAL_CAST, f.ModelID(), "string")
		}
	}

	valLen := int64(len([]rune(val)))
	if f.MinLength() != nil && valLen < *f.MinLength() {
		return true, fmt.Errorf(ER_VAL_LEN_TOO_SHORT, f.Descr())
	}

	if f.MaxLength() != nil && valLen > *f.MaxLength() {
		return true, fmt.Errorf(ER_VAL_LEN_TOO_LONG, f.Descr())
	}

	if f.FixLength() != nil && valLen != *f.FixLength() {
		return true, fmt.Errorf(ER_VAL_LEN_NOT_FIX, f.Descr())
	}

	if f.RegExp() != "" {
		match, err := regexp.MatchString(f.RegExp(), val)
		if err != nil {
			return true, err
		}
		if !match {
			return true, fmt.Errorf(ER_VAL_REG_EXP, f.Descr())
		}
	}

	if len(f.valList) > 0 {
		res := slices.Contains(f.valList, val)
		if !res {
			return true, fmt.Errorf(ER_VAL_VAL_LIST, f.Descr())
		}
	}

	return true, nil
}
