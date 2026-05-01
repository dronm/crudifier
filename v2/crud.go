// Package crudifier manages database operations, like
// insert, update, delete.
package crudifier

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/dronm/crudifier/v2/metadata"
	"github.com/dronm/crudifier/v2/types"
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

func ModelToDBFilters(model any, filters types.DBFilters, operator types.SQLFilterOperator, join types.FilterJoin, table string) error {
	modelVal := reflect.ValueOf(model)
	modelType := reflect.TypeOf(model)
	if modelVal.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
		modelVal = modelVal.Elem()
	}
	if modelVal.Kind() != reflect.Struct {
		return fmt.Errorf(metadata.ErrModelNotAPointerOrStruct, "ModelToDbFilters")
	}

	for i := 0; i < modelVal.NumField(); i++ {
		fieldType := modelType.Field(i)
		fieldID := fieldType.Tag.Get(metadata.FieldAnnotationName)
		if fieldID == "-" || fieldID == "" {
			continue
		}
		field := modelVal.Field(i)
		filters.Add(table, fieldID, field.Interface(), operator, join)
	}

	return nil
}

func PrepareUpdateModel(keyModel any, dbUpdate types.DBUpdater) error {
	if err := ModelToDBFilters(keyModel, dbUpdate.Filter(), types.SQLFilterOperatorEq, types.SQLFilterJoinAnd, ""); err != nil {
		return err
	}

	model := dbUpdate.Model()
	modelVal := reflect.ValueOf(model)
	modelType := reflect.TypeOf(model)
	if modelVal.Kind() != reflect.Ptr {
		return fmt.Errorf(metadata.ErrModelNotAPointerOrStruct, "PrepareUpdateModel")
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
			return fmt.Errorf(ErrNoFieldInMD, "PrepareUpdateModel", fieldID)
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
		if fieldMd.DataType() == metadata.FieldTypeUndefined {
			b, err := json.Marshal(field.Interface())
			if err != nil {
				errorList.WriteString(err.Error())
				continue
			}
			dbUpdate.AddField(fieldID, b)
		} else {
			dbUpdate.AddField(fieldID, field.Interface())
		}
		// fmt.Println("PrepareUpdate fieldID:",fieldID,"value:",field, "fieldMD.DataType():", fieldMd.DataType())
	}

	if errorList.Len() > 0 {
		return &ValidationError{ErrText: errorList.String()}
	}

	if assigner, ok := dbUpdate.(interface{ AssignerLen() int }); ok && assigner.AssignerLen() == 0 {
		return fmt.Errorf(ErrNoUpdateFields, "PrepareUpdateModel")
	}

	return nil
}

// PrepareFetchModelCollection prepares for retrieving a collection of objects
// from database and maps its fields to a model.
// It first parses filters, sorters and limit from client qiery parameters.
// Aggregation functions.
func PrepareFetchModelCollection(dbSelect types.DBSelecter, params CollectionParams) error {
	if err := ParseFilterParams(dbSelect.Model(), dbSelect.Filter(), params, ""); err != nil {
		return fmt.Errorf("ParseFilterParams(): %v", err)
	}

	//Aggregation functions
	aggModel := dbSelect.Model().CollectionAgg()
	if aggModel != nil {
		aggModelVal := reflect.ValueOf(aggModel)
		aggModelType := reflect.TypeOf(aggModel)
		if aggModelVal.Kind() != reflect.Ptr {
			return fmt.Errorf(metadata.ErrModelNotAPointer, "PrepareFetchModelCollection")
		}
		aggModelType = aggModelType.Elem()
		aggModelVal = aggModelVal.Elem()
		for i := 0; i < aggModelVal.NumField(); i++ {
			aggFieldType := aggModelType.Field(i)
			fieldID := aggFieldType.Tag.Get(metadata.FieldAnnotationName)
			if fieldID == "-" || fieldID == "" {
				return fmt.Errorf(ErrAggFieldNotDefined, i)
			}
			aggFunc := aggFieldType.Tag.Get(metadata.AnnotTagAgg)
			if aggFunc == "" {
				return fmt.Errorf(ErrAggFieldNoFnc, fieldID)
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
		return fmt.Errorf(metadata.ErrModelNotAPointer, "PrepareFetchModel")
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
// to a model.
func PrepareFetchModel(keyModel any, dbSelect types.DBDetailSelecter) error {
	filters := dbSelect.Filter()
	if err := ModelToDBFilters(keyModel, filters, types.SQLFilterOperatorEq, types.SQLFilterJoinAnd, ""); err != nil {
		return err
	}
	if filters.Len() == 0 {
		return fmt.Errorf(ErrNoKeys, "PrepareFetchModel")
	}

	return prepareSelectModel(dbSelect, dbSelect.Model())
}

// PrepareInsertModel analyses model proviede in dbInsert.
// It validates field value, constructs fields for insertion
// and fields for returing values.
func PrepareInsertModel(dbInsert types.DBInserter) error {
	//server autocalc fields need to be returned
	model := dbInsert.Model()
	modelVal := reflect.ValueOf(model)
	modelType := reflect.TypeOf(model)
	if modelVal.Kind() != reflect.Ptr {
		return fmt.Errorf(metadata.ErrModelNotAPointer, "PrepareInsertModel")
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
			return fmt.Errorf(ErrNoFieldInMD, "PrepareInsertModel", fieldID)
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
		if fieldMd.DataType() == metadata.FieldTypeUndefined {
			b, err := json.Marshal(field.Interface())
			if err != nil {
				errorList.WriteString(err.Error())
				continue
			}
			dbInsert.AddField(fieldID, b)
		} else {
			dbInsert.AddField(fieldID, field.Interface())
		}

	}
	if errorList.Len() > 0 {
		return &ValidationError{ErrText: errorList.String()}
	}

	if insert, ok := dbInsert.(interface{ InsertFieldLen() int }); ok && insert.InsertFieldLen() == 0 {
		return fmt.Errorf(ErrNoInsertFields, "PrepareInsertModel")
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
		return fmt.Errorf(metadata.ErrModelNotAPointer, "ValidateModel")
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
			return fmt.Errorf(ErrNoFieldInMD, "ValidateModel", fieldID)
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
