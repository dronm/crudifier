package crudifier

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dronm/crudifier/metadata"
	"github.com/dronm/crudifier/types"
)

type DbField interface {
	IsSet() bool
	IsNull() bool
}

type ValidationError struct {
	ErrText string
}

func (e *ValidationError) Error() string {
	return e.ErrText
}

func ModelToDbFilters(model interface{}, filters types.DbFilters) error {
	modelVal := reflect.ValueOf(model)
	modelType := reflect.TypeOf(model)
	if modelVal.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
		modelVal = modelVal.Elem()
	}
	if modelVal.Kind() != reflect.Struct {
		return fmt.Errorf(metadata.ER_MODEL_NOT_A_POINTER_OR_STRUCT, "PrepareUpdateModel")
	}

	for i := 0; i < modelVal.NumField(); i++ {
		fieldType := modelType.Field(i)
		fieldId := fieldType.Tag.Get(metadata.FieldAnnotationName)
		if fieldId == "-" || fieldId == "" {
			continue
		}
		field := modelVal.Field(i)
		filters.Add(fieldId, field.Interface(),
			types.SQL_FILTER_OPERATOR_E, types.SQL_FILTER_JOIN_AND)
	}

	return nil
}

// PrepareUpdateModel
func PrepareUpdateModel(keyModel interface{}, dbUpdate types.DbUpdater) error {
	if err := ModelToDbFilters(keyModel, dbUpdate.Filter()); err != nil {
		return err
	}

	model := dbUpdate.Model()
	modelVal := reflect.ValueOf(model)
	modelType := reflect.TypeOf(model)
	if modelVal.Kind() != reflect.Ptr {
		return fmt.Errorf(metadata.ER_MODEL_NOT_A_POINTER_OR_STRUCT, "PrepareUpdateModel")
	}

	modelType = modelType.Elem()
	modelVal = modelVal.Elem()

	modelMd, err := metadata.NewModelMetadata(model)
	if err != nil {
		return err
	}

	var errorList strings.Builder
	for i := 0; i < modelVal.NumField(); i++ {
		fieldType := modelType.Field(i)
		fieldId := fieldType.Tag.Get(metadata.FieldAnnotationName)
		if fieldId == "-" || fieldId == "" {
			continue
		}

		field := modelVal.Field(i)

		fieldMd, ok := modelMd.Fields[fieldId]
		if !ok {
			return fmt.Errorf(ER_NOT_FIELD_IN_MD, "PrepareUpdateModel", fieldId)
		}

		//if value is not present in model or it is not valid - skip field
		if isSet, err := fieldMd.Validate(field); err != nil {
			if errorList.Len() > 0 {
				errorList.WriteString(" ")
			}
			errorList.WriteString(err.Error())
			continue

		} else if !isSet {
			continue
		}

		//value exists and valid
		dbUpdate.AddField(fieldId, field.Interface())
	}

	if errorList.Len() > 0 {
		return &ValidationError{ErrText: errorList.String()}
	}

	return nil
}

// PrepareFetchModelCollection prepares for retrieving a collection of objects
// from database and maps its fields to model.
// It first parses filters, sorters and limit from client qiery parameters.
// Aggregation functions.
func PrepareFetchModelCollection(dbSelect types.DbSelecter, params CollectionParams) error {
	if err := ParseFilterParams(dbSelect.Model(), dbSelect.Filter(), params); err != nil {
		return err
	}

	//Aggregation functions
	aggModel := dbSelect.Model().CollectionAgg()
	if aggModel != nil {
		aggModelVal := reflect.ValueOf(aggModel)
		aggModelType := reflect.TypeOf(aggModel)
		if aggModelVal.Kind() != reflect.Ptr {
			return fmt.Errorf(metadata.ER_MODEL_NOT_A_POINTER, "PrepareFetchModelCollection")
		}
		aggModelType = aggModelType.Elem()
		aggModelVal = aggModelVal.Elem()
		for i := 0; i < aggModelVal.NumField(); i++ {
			aggFieldType := aggModelType.Field(i)
			fieldId := aggFieldType.Tag.Get(metadata.FieldAnnotationName)
			if fieldId == "-" || fieldId == "" {
				return fmt.Errorf(ER_AGG_FIELD_NOT_DEFINED, i)
			}
			aggFunc := aggFieldType.Tag.Get(metadata.ANNOT_TAG_AGG)
			if aggFunc == "" {
				return fmt.Errorf(ER_AGG_FIELD_NO_FUNC, fieldId)
			}
			field := aggModelVal.Field(i)

			//value for scanning
			dbSelect.AddAggField(aggFunc, field.Addr().Interface())
		}
	}

	if err := ParseSorterParams(dbSelect.Model(), dbSelect.Sorter(), params); err != nil {
		return err
	}

	if err := ParseLimitParams(dbSelect.Limit(), params); err != nil {
		return err
	}

	return prepareSelectModel(dbSelect)
}

func prepareSelectModel(dbSelect types.DbSelecter) error {
	model := dbSelect.Model()
	modelVal := reflect.ValueOf(model)
	modelType := reflect.TypeOf(model)
	if modelVal.Kind() != reflect.Ptr {
		return fmt.Errorf(metadata.ER_MODEL_NOT_A_POINTER, "PrepareFetchModel")
	}

	modelType = modelType.Elem()
	modelVal = modelVal.Elem()

	for i := 0; i < modelVal.NumField(); i++ {
		fieldType := modelType.Field(i)
		fieldId := fieldType.Tag.Get(metadata.FieldAnnotationName)
		if fieldId == "-" || fieldId == "" {
			continue
		}
		field := modelVal.Field(i)
		dbSelect.AddField(fieldId, field.Addr().Interface())
	}
	return nil
}

// PrepareFetchModel prepares for retrieving one object from database and maps its fields
// to model.
func PrepareFetchModel(keyModel interface{}, dbSelect types.DbSelecter) error {
	filters := dbSelect.Filter()
	if err := ModelToDbFilters(keyModel, filters); err != nil {
		return err
	}
	if filters.Len() == 0 {
		return fmt.Errorf(ER_NOT_KEYS, "PrepareFetchModel")
	}

	return prepareSelectModel(dbSelect)
}

// PrepareInsertModel analyses model proviede in dbInsert.
// It validates field value, constructs fields for insertion
// and fields for returing values.
func PrepareInsertModel(dbInsert types.DbInserter) error {
	//server autocalc fields need to be returned
	model := dbInsert.Model()
	modelVal := reflect.ValueOf(model)
	modelType := reflect.TypeOf(model)
	if modelVal.Kind() != reflect.Ptr {
		return fmt.Errorf(metadata.ER_MODEL_NOT_A_POINTER, "PrepareInsertModel")
	}

	modelType = modelType.Elem()
	modelVal = modelVal.Elem()

	modelMd, err := metadata.NewModelMetadata(model)
	if err != nil {
		return err
	}

	var errorList strings.Builder
	for i := 0; i < modelVal.NumField(); i++ {
		fieldType := modelType.Field(i)
		fieldId := fieldType.Tag.Get(metadata.FieldAnnotationName)
		if fieldId == "-" || fieldId == "" {
			continue
		}
		field := modelVal.Field(i)
		if !field.CanInterface() {
			continue
		}

		fieldMd, ok := modelMd.Fields[fieldId]
		if !ok {
			return fmt.Errorf(ER_NOT_FIELD_IN_MD, "PrepareInsertModel", fieldId)
		}

		if err := fieldMd.ValidateRequired(field); err != nil {
			if errorList.Len() > 0 {
				errorList.WriteString(" ")
			}
			errorList.WriteString(err.Error())
			continue
		}

		if fieldMd.SrvCalc() {
			//should be returned
			dbInsert.AddRetField(fieldId, field.Addr().Interface())
			continue
		}

		if present, err := fieldMd.Validate(field); err != nil {
			if errorList.Len() > 0 {
				errorList.WriteString(" ")
			}
			errorList.WriteString(err.Error())
			continue

		} else if !present {
			continue
		}

		//ordinary field, present && valid
		dbInsert.AddField(fieldId, field.Interface())
	}
	if errorList.Len() > 0 {
		return &ValidationError{ErrText: errorList.String()}
	}

	return nil
}
