package types

type DbLimit interface {
	From() int
	SetFrom(int)
	Count() int
	SetCount(int)
	SQL() string
}
