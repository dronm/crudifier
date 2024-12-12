package metadata

import (
	"reflect"
	"strconv"
)

var FieldAnnotationName = "json"

// All possible annotations.
const (
	//common
	ANNOT_TAG_ALIAS       = "alias"
	ANNOT_TAG_REQUIRED    = "required"
	ANNOT_TAG_DB_REQUIRED = "dbRequired" //required for Database, can be set with trigger/autoinc
	ANNOT_TAG_PRIM_KEY    = "primaryKey"
	ANNOT_TAG_SRV_CALC    = "srvCalc" //server initialized field on insert

	//
	ANNOT_TAG_AGG = "agg" //aggregation function, like agg:"count(*)" for agg models

	//text
	ANNOT_TAG_MAX_LEN = "max"
	ANNOT_TAG_MIN_LEN = "min"
	ANNOT_TAG_FIX_LEN = "fix"
	ANNOT_TAG_REG_EXP = "regExp"

	//number
	ANNOT_TAG_MAX_VAL = "max"
	ANNOT_TAG_MIN_VAL = "min"
	ANNOT_TAG_FIX_VAL = "fix"
)

func annotationTagBoolVal(fieldType reflect.StructField, tagName string) bool {
	_, present := fieldType.Tag.Lookup(tagName)
	return present
}

func annotationTagStringVal(fieldType reflect.StructField, tagName string) string {
	return fieldType.Tag.Get(tagName)
}

func annotationTagIntVal(fieldType reflect.StructField, tagName string) (*int64, error) {
	tagVal := fieldType.Tag.Get(tagName)
	if tagVal == "" {
		return nil, nil
	}
	iVal, err := strconv.ParseInt(tagVal, 10, 64)
	if err != nil {
		return nil, err
	}
	return &iVal, nil
}

func annotationTagFloatVal(fieldType reflect.StructField, tagName string) (*float64, error) {
	tagVal := fieldType.Tag.Get(tagName)
	if tagVal == "" {
		return nil, nil
	}
	fVal, err := strconv.ParseFloat(tagVal, 64)
	if err != nil {
		return nil, err
	}
	return &fVal, nil
}

func setTextValidatorConstraints(field reflect.StructField, validator *FieldTextMetadata) error {
	tagVal, err := annotationTagIntVal(field, ANNOT_TAG_MAX_LEN)
	if err != nil {
		return err
	}
	validator.maxLength = tagVal

	tagVal, err = annotationTagIntVal(field, ANNOT_TAG_MIN_LEN)
	if err != nil {
		return err
	}
	validator.minLength = tagVal

	tagVal, err = annotationTagIntVal(field, ANNOT_TAG_FIX_LEN)
	if err != nil {
		return err
	}
	validator.fixLength = tagVal

	if tagVal := annotationTagStringVal(field, ANNOT_TAG_REG_EXP); tagVal != "" {
		validator.regExp = tagVal
	}

	return nil
}

func setIntValidatorConstraints(field reflect.StructField, validator *FieldIntMetadata) error {
	tagVal, err := annotationTagIntVal(field, ANNOT_TAG_MIN_VAL)
	if err != nil {
		return err
	}
	validator.minValue = tagVal

	tagVal, err = annotationTagIntVal(field, ANNOT_TAG_MAX_VAL)
	if err != nil {
		return err
	}
	validator.maxValue = tagVal

	return nil
}

func setFloatValidatorConstraints(field reflect.StructField, validator *FieldFloatMetadata) error {
	tagVal, err := annotationTagFloatVal(field, ANNOT_TAG_MIN_VAL)
	if err != nil {
		return err
	}
	validator.minValue = tagVal

	tagVal, err = annotationTagFloatVal(field, ANNOT_TAG_MAX_VAL)
	if err != nil {
		return err
	}
	validator.maxValue = tagVal

	return nil
}
