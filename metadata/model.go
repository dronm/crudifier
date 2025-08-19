// Package metadata
package metadata

import (
	"fmt"
	"reflect"
	"sync"
)

var ControledTags []string //defined at startup, tags to cache

// ModelMetadata holds the mapping of json tags to field types.
type ModelMetadata struct {
	ID           string
	Fields       map[string]FieldValidator    // FieldMetadata, key is an sql field
	FieldList    []string                     // structure field Names in original order, 
	FieldTagList []string                     // sql fields in original order for retrieving metadata from Firlds
	Tags         map[string]map[string]string // controled tag values. List of controled tags is defined in ControledTags
}

// cache to store metadata for different types.
var (
	metadataCache = make(map[reflect.Type]ModelMetadata)
	cacheMutex    sync.RWMutex
)

// NewModelMetadata returns the Metadata structure for a given type.
// Internally it stores a race-safe cache.
func NewModelMetadata(model any) (*ModelMetadata, error) {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem() // Dereference pointer types
	}

	if modelType.Kind() != reflect.Struct {
		return nil, fmt.Errorf(ER_MODEL_NOT_A_POINTER_OR_STRUCT, "NewModelMetadata")
	}

	// Check the cache for existing metadata
	cacheMutex.RLock()
	if meta, found := metadataCache[modelType]; found {
		cacheMutex.RUnlock()
		return &meta, nil
	}
	cacheMutex.RUnlock()

	// Build metadata if not found in cache
	meta := ModelMetadata{ID: modelType.String(), Fields: make(map[string]FieldValidator)}
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fieldID := field.Name
		fieldTagVal := field.Tag.Get(FieldAnnotationName)

		// Skip fields without a FieldFilterAnnotationName tag
		if fieldTagVal == "-" || fieldTagVal == "" {
			// check filter tag
			fieldTagVal = field.Tag.Get(FieldFilterAnnotationName)
			if fieldTagVal == "-" || fieldTagVal == "" {
				continue
			}
		}

		for _, tag := range ControledTags {
			tagVal := field.Tag.Get(tag)
			if tagVal != "" {
				if meta.Tags == nil {
					meta.Tags = make(map[string]map[string]string)
				}
				if meta.Tags[fieldID] == nil {
					meta.Tags[fieldID] = make(map[string]string)
				}
				meta.Tags[fieldID][tag] = tagVal
			}
		}

		meta.FieldList = append(meta.FieldList, fieldID)
		meta.FieldTagList = append(meta.FieldTagList, fieldTagVal)

		fieldType := ""
		if field.Type.Kind() == reflect.Ptr {
			fieldType = field.Type.Elem().Name()
		} else {
			fieldType = field.Type.Name()
		}
		switch ParseFieldType(fieldType) {
		case FIELD_TYPE_BOOL:
			// no contraints
			meta.Fields[fieldTagVal] = NewFieldBoolMedata(fieldID, fieldTagVal)

		case FIELD_TYPE_TEXT:
			validator := NewFieldTextMedata(fieldID, fieldTagVal)
			setTextValidatorConstraints(field, validator)
			meta.Fields[fieldTagVal] = validator

		case FIELD_TYPE_INT:
			validator := NewFieldIntMedata(fieldID, fieldTagVal)
			setIntValidatorConstraints(field, validator)
			meta.Fields[fieldTagVal] = validator

		case FIELD_TYPE_FLOAT:
			validator := NewFieldFloatMedata(fieldID, fieldTagVal)
			setFloatValidatorConstraints(field, validator)
			meta.Fields[fieldTagVal] = validator

		case FIELD_TYPE_DATE:
			validator := NewFieldDateMedata(fieldID, fieldTagVal, FIELD_TYPE_DATE)
			meta.Fields[fieldTagVal] = validator

		case FIELD_TYPE_DATETIME:
			validator := NewFieldDateMedata(fieldID, fieldTagVal, FIELD_TYPE_DATETIME)
			meta.Fields[fieldTagVal] = validator

		case FIELD_TYPE_DATETIMETZ:
			validator := NewFieldDateMedata(fieldID, fieldTagVal, FIELD_TYPE_DATETIMETZ)
			meta.Fields[fieldTagVal] = validator
		default:
			meta.Fields[fieldTagVal] = &FieldMetadata{modelId: fieldID, id: fieldTagVal}
		}

		// common tags
		meta.Fields[fieldTagVal].SetAlias(annotationTagStringVal(field, ANNOT_TAG_ALIAS))
		meta.Fields[fieldTagVal].SetRequired(annotationTagBoolVal(field, ANNOT_TAG_REQUIRED))
		meta.Fields[fieldTagVal].SetPrimaryKey(annotationTagBoolVal(field, ANNOT_TAG_PRIM_KEY))
		meta.Fields[fieldTagVal].SetSrvCalc(annotationTagBoolVal(field, ANNOT_TAG_SRV_CALC))
	}

	// Save to cache
	cacheMutex.Lock()
	metadataCache[modelType] = meta
	cacheMutex.Unlock()

	return &meta, nil
}
