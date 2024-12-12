package types

type DbSorters interface {
	Add(fieldId string, direct SQLSortDirect)
	SQL() string
	Len() int
}
