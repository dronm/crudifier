package types

type DbInserter interface {
	Model() DbModel
	AddField(id string, val interface{})
	AddRetField(id string, val interface{})
}
