package crudifier

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/dronm/crudifier/metadata"
	"github.com/dronm/crudifier/types"
)

type DBField interface {
	IsSet() bool
	IsNull() bool
}

type ValidationError struct {
	ErrText string
}

func (e *ValidationError) Error() string {
	return e.ErrText
}

func NewValidationError(errText string) *ValidationError {
	return &ValidationError{ErrText: errText}
}

func ModelToDBFilters(model any, filters types.DbFilters, operator types.SQLFilterOperator, join types.FilterJoin) error {
	modelVal := reflect.ValueOf(model)
	modelType := reflect.TypeOf(model)
	if modelVal.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
		modelVal = modelVal.Elem()
	}
	if modelVal.Kind() != reflect.Struct {
		return fmt.Errorf(metadata.ER_MODEL_NOT_A_POINTER_OR_STRUCT, "ModelToDbFilters")
	}

	for i := 0; i < modelVal.NumField(); i++ {
		fieldType := modelType.Field(i)
		fieldID := fieldType.Tag.Get(metadata.FieldAnnotationName)
		if fieldID == "-" || fieldID == "" {
			continue
		}
		field := modelVal.Field(i)
		filters.Add(fieldID, field.Interface(), operator, join)
	}

	return nil
}

// PrepareUpdateModel
func PrepareUpdateModel(keyModel any, dbUpdate types.DbUpdater) error {
	if err := ModelToDBFilters(keyModel, dbUpdate.Filter(), types.SQL_FILTER_OPERATOR_E, types.SQL_FILTER_JOIN_AND); err != nil {
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

		if !field.CanInterface() {
			return fmt.Errorf("reflect.CanInterface() failed for field %s", fieldId)
			// continue
		}

		if !field.IsValid() {
			return fmt.Errorf("reflect.IsValid() failed for field %s", fieldId)
		}

		fieldMd, ok := modelMd.Fields[fieldId]
		if !ok {
			return fmt.Errorf(ER_NO_FIELD_IN_MD, "PrepareUpdateModel", fieldId)
		}

		//if value is not present in model or it is not valid - skip field
		if isSet, err := fieldMd.Validate(field); err != nil {
			if errorList.Len() > 0 {
				errorList.WriteString(" ")
			}
			errorList.WriteString(err.Error())
			continue

		} else if !isSet {
			continue // if no value - skeep required checking
		}

		if err := fieldMd.ValidateRequired(field); err != nil {
			if errorList.Len() > 0 {
				errorList.WriteString(" ")
			}
			errorList.WriteString(err.Error())
			continue
		}

		//value exists and valid
		if fieldMd.DataType() == metadata.FIELD_TYPE_UNDEFINED {
			b, err := json.Marshal(field.Interface())
			if err != nil {
				errorList.WriteString(err.Error())
				continue
			}
			dbUpdate.AddField(fieldId, b)
		}else{
			dbUpdate.AddField(fieldId, field.Interface())
		}
		// fmt.Println("PrepareUpdate fieldId:",fieldId,"value:",field, "fieldMD.DataType():", fieldMd.DataType())
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
		return fmt.Errorf("ParseFilterParams(): %v", err)
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
			fieldID := aggFieldType.Tag.Get(metadata.FieldAnnotationName)
			if fieldID == "-" || fieldID == "" {
				return fmt.Errorf(ER_AGG_FIELD_NOT_DEFINED, i)
			}
			aggFunc := aggFieldType.Tag.Get(metadata.AnnotTagAgg)
			if aggFunc == "" {
				return fmt.Errorf(ER_AGG_FIELD_NO_FUNC, fieldID)
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

	return prepareSelectModel(dbSelect.(types.PrepareModel), dbSelect.Model())
}

func prepareSelectModel(selectModel types.PrepareModel, model any) error {
	modelVal := reflect.ValueOf(model)
	modelType := reflect.TypeOf(model)
	if modelVal.Kind() != reflect.Ptr {
		return fmt.Errorf(metadata.ER_MODEL_NOT_A_POINTER, "PrepareFetchModel")
	}

	modelType = modelType.Elem()
	modelVal = modelVal.Elem()

	for i := 0; i < modelVal.NumField(); i++ {
		fieldType := modelType.Field(i)
		fieldID := fieldType.Tag.Get(metadata.FieldAnnotationName)
		if fieldID == "-" || fieldID == "" {
			continue
		}
		field := modelVal.Field(i)
		selectModel.AddField(fieldID, field.Addr().Interface())
	}
	return nil
}

// PrepareFetchModel prepares for retrieving one object from database and maps its fields
// to model.
func PrepareFetchModel(keyModel any, dbSelect types.DbDetailSelecter) error {
	filters := dbSelect.Filter()
	if err := ModelToDBFilters(keyModel, filters, types.SQL_FILTER_OPERATOR_E, types.SQL_FILTER_JOIN_AND); err != nil {
		return err
	}
	if filters.Len() == 0 {
		return fmt.Errorf(ER_NO_KEYS, "PrepareFetchModel")
	}

	return prepareSelectModel(dbSelect, dbSelect.Model())
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
		fieldID := fieldType.Tag.Get(metadata.FieldAnnotationName)
		if fieldID == "-" || fieldID == "" {
			continue
		}
		field := modelVal.Field(i)
		if !field.CanInterface() {
			return fmt.Errorf("reflect.CanInterface() failed for field %s", fieldID)
			// continue
		}

		if !field.IsValid() {
			return fmt.Errorf("reflect.IsValid() failed for field %s", fieldID)
		}

		fieldMd, ok := modelMd.Fields[fieldID]
		if !ok {
			return fmt.Errorf(ER_NO_FIELD_IN_MD, "PrepareInsertModel", fieldID)
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
			dbInsert.AddRetField(fieldID, field.Addr().Interface())
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
		if fieldMd.DataType() == metadata.FIELD_TYPE_UNDEFINED {
			b, err := json.Marshal(field.Interface())
			if err != nil {
				errorList.WriteString(err.Error())
				continue
			}
			dbInsert.AddField(fieldID, b)
		}else{
			dbInsert.AddField(fieldID, field.Interface())
		}

	}
	if errorList.Len() > 0 {
		return &ValidationError{ErrText: errorList.String()}
	}

	return nil
}

// ValidateModel checks the given model. Validation error is thrown in case of 
// forInsert flag is set to true, a field is required but not set or null.
// If forInsert is false and field is not set, required field annotation is not checked.
// In this latter case the error is thrown only the field value is set to null.
func ValidateModel(model any, forInsert bool) error {
	//server autocalc fields need to be returned
	modelVal := reflect.ValueOf(model)
	modelType := reflect.TypeOf(model)
	if modelVal.Kind() != reflect.Ptr {
		return fmt.Errorf(metadata.ER_MODEL_NOT_A_POINTER, "ValidateModel")
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
			return fmt.Errorf("reflect.CanInterface() failed for field %s", fieldId)
			// continue
		}

		if !field.IsValid() {
			return fmt.Errorf("reflect.IsValid() failed for field %s", fieldId)
		}

		fieldMd, ok := modelMd.Fields[fieldId]
		if !ok {
			return fmt.Errorf(ER_NO_FIELD_IN_MD, "ValidateModel", fieldId)
		}

		if fieldMd.SrvCalc() {
			continue
		}

		if isSet, err := fieldMd.Validate(field); err != nil {
			if errorList.Len() > 0 {
				errorList.WriteString(" ")
			}
			errorList.WriteString(err.Error())

		} else if !isSet && !forInsert {
			continue // if no value - skeep required checking
		}

		if err := fieldMd.ValidateRequired(field); err != nil {
			if errorList.Len() > 0 {
				errorList.WriteString(" ")
			}
			errorList.WriteString(err.Error())
			continue
		}
	}

	if errorList.Len() > 0 {
		return &ValidationError{ErrText: errorList.String()}
	}

	return nil
}
