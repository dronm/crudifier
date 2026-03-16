package crudifier

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/dronm/crudifier/v2/metadata"
	"github.com/dronm/crudifier/v2/pg"
	"github.com/dronm/crudifier/v2/types"
)

func ParseLimitParams(dbLimit types.DBLimit, params CollectionParams) error {
	if dbLimit == nil {
		return fmt.Errorf(ErrLimitNotInit, "ParseLimitParams")
	}
	dbLimit.SetFrom(int(params.From))
	dbLimit.SetCount(int(params.Count))
	return nil
}

func ParseSorterParams(model types.DBModel, dbSorter types.DBSorters, params CollectionParams) error {
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
			return fmt.Errorf(ErrNoFieldInMD, "ParseSorterParams", fieldID)
		}
		if _, err := pg.SanitizeSQLFieldRef(sorter.Field); err != nil {
			return fmt.Errorf(ErrInvalidFieldExpr, "ParseSorterParams", sorter.Field)
		}

		var sortDirect types.SQLSortDirect
		switch sorter.Direct {
		case SortParDesc:
			sortDirect = types.SQLSortDesc
		default:
			sortDirect = types.SQLSortAsc
		}

		if dbSorter == nil {
			return fmt.Errorf(ErrSorterNotInit, "ParseSorterParams")
		}
		dbSorter.Add(sorter.Field, sortDirect)
	}

	return nil
}

func ParseFilterParams(model types.DBModel, dbFilter types.DBFilters, params CollectionParams) error {
	modelMd, err := metadata.NewModelMetadata(model)
	if err != nil {
		return err
	}

	var errorList strings.Builder

	for _, filter := range params.Filter {
		var join types.FilterJoin
		switch filter.Join {
		case FilterParJoinOr:
			join = types.SQLFilterJoinOr
		default:
			join = types.SQLFilterJoinAnd
		}

		for filterFieldID, filterField := range filter.Fields {
			filterModelFieldID := filterFieldID

			structFieldInd := strings.Index(filterModelFieldID, "->")
			if structFieldInd >= 0 {
				filterModelFieldID = filterModelFieldID[:structFieldInd]
			}

			if _, err := pg.SanitizeSQLFieldRef(filterFieldID); err != nil {
				return fmt.Errorf(ErrInvalidFieldExpr, "ParseFilterParams", filterFieldID)
			}

			fieldMd, ok := modelMd.Fields[filterModelFieldID]
			if !ok {
				return fmt.Errorf(ErrNoFieldInMD, "ParseFilterParams", filterModelFieldID)
			}

			if filterField.Value != nil {
				if _, err := fieldMd.Validate(reflect.ValueOf(filterField.Value)); err != nil {
					if errorList.Len() > 0 {
						errorList.WriteString(" ")
					}
					errorList.WriteString(err.Error())
					continue
				}
			}

			var operator types.SQLFilterOperator
			switch filterField.Operator {
			case FilterOperParEq:
				operator = types.SQLFilterOperatorEq
			case FilterOperParLess:
				operator = types.SQLFilterOperatorLess
			case FilterOperParGr:
				operator = types.SQLFilterOperatorGr
			case FilterOperParLessEq:
				operator = types.SQLFilterOperatorLessEq
			case FilterOperParGrEq:
				operator = types.SQLFilterOperatorGrEq
			case FilterOperParLk:
				operator = types.SQLFilterOperatorLk
			case FilterOperParILk:
				operator = types.SQLFilterOperatorILk
			case FilterOperParNotEq:
				operator = types.SQLFilterOperatorNotEq
			case FilterOperParIs:
				operator = types.SQLFilterOperatorIs
			case FilterOperParIn:
				operator = types.SQLFilterOperatorIn
			case FilterOperParIncl:
				operator = types.SQLFilterOperatorIncl
			case FilterOperParAny:
				operator = types.SQLFilterOperatorAny
			case FilterOperParHas:
				operator = types.SQLFilterOperatorHas
			case FilterOperParOverlap:
				operator = types.SQLFilterOperatorOverlap
			case FilterOperParContains:
				operator = types.SQLFilterOperatorContains
			case FilterOperParTS:
				operator = types.SQLFilterOperatorTS
			}

			if dbFilter == nil {
				return fmt.Errorf(ErrFilterNotInit, "ParseFilterParams")
			}
			switch operator {
			case types.SQLFilterOperatorTS:
				dbFilter.AddFullTextSearch(filterFieldID, filterField.Value, join)
			case types.SQLFilterOperatorIncl:
				dbFilter.AddArrayInclude(filterFieldID, filterField.Value, join)
			case types.SQLFilterOperatorHas:
				dbFilter.AddColumnArrayInclude(filterFieldID, filterField.Value, join)
			default:
				dbFilter.Add(filterFieldID, filterField.Value, operator, join)
			}
		}
	}

	if errorList.Len() > 0 {
		return errors.New(errorList.String())
	}

	return nil
}
