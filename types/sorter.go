package types

type SQLSortDirect string

const (
	SQL_SORT_ASC  SQLSortDirect = "ASC"
	SQL_SORT_DESC SQLSortDirect = "DESC"
)

type SQLSorter interface {
	FieldID() string
	Direct() SQLSortDirect
}
