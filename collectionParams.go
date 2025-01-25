package crudifier

import "encoding/json"

// This file provides types and functions for parsing
// collection parameters from client queries.

// FilterOperatorParam is a clietn query operator.
type FilterOperatorParam string

// client query fileter operator values
const (
	FILTER_OPER_PAR_E       FilterOperatorParam = "e"       //equal
	FILTER_OPER_PAR_L       FilterOperatorParam = "l"       //less
	FILTER_OPER_PAR_G       FilterOperatorParam = "g"       //greater
	FILTER_OPER_PAR_LE      FilterOperatorParam = "le"      //less and equal
	FILTER_OPER_PAR_GE      FilterOperatorParam = "ge"      //greater and equal
	FILTER_OPER_PAR_LK      FilterOperatorParam = "lk"      //like
	FILTER_OPER_PAR_ILK     FilterOperatorParam = "ilk"     //ilike
	FILTER_OPER_PAR_NE      FilterOperatorParam = "ne"      //not equal
	FILTER_OPER_PAR_I       FilterOperatorParam = "i"       // IS
	FILTER_OPER_PAR_IN      FilterOperatorParam = "in"      // in
	FILTER_OPER_PAR_INCL    FilterOperatorParam = "incl"    //include
	FILTER_OPER_PAR_ANY     FilterOperatorParam = "any"     //Any
	FILTER_OPER_PAR_OVERLAP FilterOperatorParam = "overlap" //overlap
)

type FilterJoinParam string

const (
	FILTER_PAR_JOIN_AND FilterJoinParam = "and"
	FILTER_PAR_JOIN_OR  FilterJoinParam = "or"
)

type SortParam string

// client query sorting values
const (
	SORT_PAR_ASC  SortParam = "a" // asc
	SORT_PAR_DESC SortParam = "d" // desc
)

type CollectionSorter struct {
	Field  string    `json:"f"`
	Direct SortParam `json:"d"`
}

type CollectionFilterField struct {
	Operator FilterOperatorParam `json:"o"`
	Value    interface{}         `json:"v"`
}

type CollectionFilter struct {
	Join   FilterJoinParam                  `json:"j"`
	Fields map[string]CollectionFilterField `json:"f"`
}

type CollectionFrom int
type CollectionCount int

// CollectionParams holds all unmarshaled
// client query parameters.
type CollectionParams struct {
	Filter []CollectionFilter `json:"filter"`
	Sorter []CollectionSorter `json:"sorter"`
	From   CollectionFrom     `json:"from"`
	Count  CollectionCount    `json:"count"`
}

func (p *CollectionParams) ParseQuery(queryParams string) error {
	if err := json.Unmarshal([]byte(queryParams), p); err != nil {
		return err
	}

	return nil
}
