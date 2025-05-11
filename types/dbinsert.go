package types

type DbInserter interface {
	Model() DbModel
	AddField(id string, val any)
	AddRetField(id string, val any)
}
