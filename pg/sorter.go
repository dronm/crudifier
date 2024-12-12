package pg

import (
	"fmt"

	"github.com/dronm/crudifier/types"
)

type PgSorter struct {
	fieldID   string
	direct    types.SQLSortDirect
	fieldPref string
}

func (s PgSorter) FieldID() string {
	return s.fieldID
}

func (s PgSorter) Direct() types.SQLSortDirect {
	return s.direct
}

func (s PgSorter) FieldPref() string {
	return s.fieldPref
}

func (s PgSorter) SQL() string {
	if s.fieldPref != "" {
		return fmt.Sprintf("%s.%s %s", s.fieldPref, s.fieldID, s.direct)
	}

	return fmt.Sprintf("%s %s", s.fieldID, s.direct)
}
