package crudifier

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/dronm/crudifier/metadata"
	"github.com/dronm/crudifier/types"
)

func ParseLimitParams(dbLimit types.DbLimit, params CollectionParams) error {
	if dbLimit == nil {
		return fmt.Errorf(ER_LIMIT_NOT_INIT, "ParseLimitParams")
	}
	dbLimit.SetFrom(int(params.From))
	dbLimit.SetCount(int(params.Count))
	return nil
}

func ParseSorterParams(model types.DbModel, dbSorter types.DbSorters, params CollectionParams) error {
	modelMd, err := metadata.NewModelMetadata(model)
	if err != nil {
		return err
	}

	for _, sorter := range params.Sorter {
		fieldID := sorter.Field
		structFieldInd := strings.Index(fieldID, "->")
		if structFieldInd >= 0 {
			fieldID = fieldID[:structFieldInd]
		}
		if _, ok := modelMd.Fields[fieldID]; !ok {
			return fmt.Errorf(ER_NO_FIELD_IN_MD, "ParseSorterParams", fieldID)
		}

		var sortDirect types.SQLSortDirect
		switch sorter.Direct {
		case SORT_PAR_DESC:
			sortDirect = types.SQL_SORT_DESC
		default:
			sortDirect = types.SQL_SORT_ASC
		}

		if dbSorter == nil {
			return fmt.Errorf(ER_SORTER_NOT_INIT, "ParseSorterParams")
		}
		dbSorter.Add(sorter.Field, sortDirect)
	}

	return nil
}

func ParseFilterParams(model types.DbModel, dbFilter types.DbFilters, params CollectionParams) error {
	modelMd, err := metadata.NewModelMetadata(model)
	if err != nil {
		return err
	}

	var errorList strings.Builder // validation errors

	// check if every filter field exists in dbSelect.Model
	for _, filter := range params.Filter {
		// join is common for all fields in this filter
		var join types.FilterJoin
		switch filter.Join {
		case FILTER_PAR_JOIN_OR:
			join = types.SQL_FILTER_JOIN_OR
		default:
			join = types.SQL_FILTER_JOIN_AND
		}

		// and filter value can be assigned to model field.
		for filterFieldID, filterField := range filter.Fields {
			// if filterField.Value == nil {
			// 	fmt.Println("ParseFilterParams skip nil value for field:",filterFieldID)
			// 	continue
			// }

			filterModelFieldID := filterFieldID

			//json field search
			structFieldInd := strings.Index(filterModelFieldID, "->")
			if structFieldInd >= 0 {
				filterModelFieldID = filterModelFieldID[:structFieldInd]
			}

			fieldMd, ok := modelMd.Fields[filterModelFieldID]
			if !ok {
				return fmt.Errorf(ER_NO_FIELD_IN_MD, "ParseFilterParams", filterModelFieldID)
			}

			//nil values must be included!
			if filterField.Value != nil {
				if _, err := fieldMd.Validate(reflect.ValueOf(filterField.Value)); err != nil {
					if errorList.Len() > 0 {
						errorList.WriteString(" ")
					}
					errorList.WriteString(err.Error())
					continue
				}
			}

			// resolve operator
			var operator types.SQLFilterOperator
			switch filterField.Operator {
			case FILTER_OPER_PAR_E:
				operator = types.SQL_FILTER_OPERATOR_E
			case FILTER_OPER_PAR_L:
				operator = types.SQL_FILTER_OPERATOR_L
			case FILTER_OPER_PAR_G:
				operator = types.SQL_FILTER_OPERATOR_G
			case FILTER_OPER_PAR_LE:
				operator = types.SQL_FILTER_OPERATOR_LE
			case FILTER_OPER_PAR_GE:
				operator = types.SQL_FILTER_OPERATOR_GE
			case FILTER_OPER_PAR_LK:
				operator = types.SQL_FILTER_OPERATOR_LK
			case FILTER_OPER_PAR_ILK:
				operator = types.SQL_FILTER_OPERATOR_ILK
			case FILTER_OPER_PAR_NE:
				operator = types.SQL_FILTER_OPERATOR_NE
			case FILTER_OPER_PAR_I:
				operator = types.SQL_FILTER_OPERATOR_I
			case FILTER_OPER_PAR_IN:
				operator = types.SQL_FILTER_OPERATOR_IN
			case FILTER_OPER_PAR_INCL:
				operator = types.SQL_FILTER_OPERATOR_INCL
			case FILTER_OPER_PAR_ANY:
				operator = types.SQL_FILTER_OPERATOR_ANY
			case FILTER_OPER_PAR_OVERLAP:
				operator = types.SQL_FILTER_OPERATOR_OVERLAP
			case FILTER_OPER_PAR_CONTAINS:
				operator = types.SQL_FILTER_OPERATOR_CONTAINS
			case FILTER_OPER_PAR_TS:
				operator = types.SQL_FILTER_OPERATOR_TS
			}

			if dbFilter == nil {
				return fmt.Errorf(ER_FILTER_NOT_INIT, "ParseFilterParams")
			}
			switch operator {
			case types.SQL_FILTER_OPERATOR_TS:
				dbFilter.AddFullTextSearch(filterFieldID, filterField.Value, join)
			case types.SQL_FILTER_OPERATOR_INCL:
				dbFilter.AddArrayInclude(filterFieldID, filterField.Value, join)
			default:
				dbFilter.Add(filterFieldID, filterField.Value, operator, join)
			}
			// if operator == types.SQL_FILTER_OPERATOR_TS {
			// 	dbFilter.AddFullTextSearch(filterFieldID, filterField.Value, join)
			// } else if operator == types.SQL_FILTER_OPERATOR_INCL {
			// 	dbFilter.AddArrayInclude(filterFieldID, filterField.Value, join)
			// } else {
			// 	dbFilter.Add(filterFieldID, filterField.Value, operator, join)
			// }
		}
	}

	if errorList.Len() > 0 {
		return errors.New(errorList.String())
	}

	return nil
}
