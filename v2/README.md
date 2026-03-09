# crudifier

`crudifier` is a small Go library for preparing CRUD operations from Go structs and rendering PostgreSQL SQL builders.

It is designed around three ideas:

- describe database fields with struct tags;
- use typed nullable/settable field wrappers for patch-style writes;
- keep SQL rendering in dedicated PostgreSQL builders.

The package does **not** execute queries by itself. It prepares models, filters, sorters, limits, SQL strings, and query parameters so you can use them with `database/sql`, `pgx`, or another database layer.

## Packages

- `crudifier` — model preparation, validation, collection query parsing;
- `metadata` — struct tag parsing and field metadata cache;
- `fields` — nullable/settable wrappers such as `FieldInt`, `FieldText`, `FieldBool`, `FieldDateTime`;
- `pg` — PostgreSQL SQL builders;
- `types` — core interfaces.

## Installation

```bash
go get github.com/dronm/crudifier
```

## Core concepts

### 1. Model relation

Each model must implement:

```go
type DBModel interface {
	Relation() string
}
```

`Relation()` returns the SQL relation used in queries. In current form this is treated as **trusted SQL**.

### 2. Field annotations

By default the library uses the `json` tag as the field identifier.

Example:

```go
type User struct {
	ID		fields.FieldInt		`json:"id" primaryKey:"" srvCalc:""`
	Name	fields.FieldText	`json:"name" required:"" min:"1" max:"200"`
	Email	fields.FieldText	`json:"email" required:"" max:"255"`
	Age		fields.FieldInt		`json:"age" min:"0"`
}
```

Important tags in `metadata`:

- `required` — required at validation time;
- `dbRequired` — required for DB persistence, but may be server-generated;
- `primaryKey` — marks a primary key field;
- `srvCalc` — field is server-calculated on insert and should be returned;
- `agg` — aggregate expression for collection aggregate models;
- `min`, `max`, `fix`, `regExp`, `enum`, `valList` — validator constraints.

Global annotation names can be changed at startup:

```go
metadata.FieldAnnotationName = "json"
metadata.FieldFilterAnnotationName = "f"
metadata.ValListSeparator = "@@"
```

## Field wrappers

The `fields` package is intended for partial writes and nullable values.

Example:

```go
var name fields.FieldText
_ = json.Unmarshal([]byte(`"Alice"`), &name)

if name.IsSet() {
	// field was present in input
}

if name.IsNull() {
	// field was explicitly set to null
}
```

This distinction is critical for update/patch flows:

- not set → do not include in `UPDATE`;
- set to `null` → include the field with `NULL`;
- set to a concrete value → include the value.

## Minimal model example

```go
package example

import "github.com/dronm/crudifier/fields"

type User struct {
	ID		fields.FieldInt		`json:"id" primaryKey:"" srvCalc:""`
	Name	fields.FieldText	`json:"name" required:"" min:"1" max:"200"`
	Email	fields.FieldText	`json:"email" required:"" max:"255"`
}

func (User) Relation() string {
	return "users"
}

func (User) CollectionAgg() any {
	return &struct {
		Count int64 `json:"count" agg:"count(*)"`
	}{}
}
```

## Insert

```go
package example

import (
	"github.com/dronm/crudifier"
	"github.com/dronm/crudifier/fields"
	"github.com/dronm/crudifier/pg"
)

func buildInsert() (string, []any, error) {
	model := &User{}
	model.Name.SetValue("Alice")
	model.Email.SetValue("alice@example.com")

	ins := pg.NewPgInsert(model)
	if err := crudifier.PrepareInsertModel(ins); err != nil {
		return "", nil, err
	}

	params := make([]any, 0)
	sql := ins.SQL(&params)
	return sql, params, nil
}
```

Typical result:

```sql
INSERT INTO users (name,email) VALUES ($1,$2) RETURNING id
```

## Update

```go
package example

import (
	"github.com/dronm/crudifier"
	"github.com/dronm/crudifier/fields"
	"github.com/dronm/crudifier/pg"
)

func buildUpdate() (string, []any, error) {
	key := &struct {
		ID fields.FieldInt `json:"id"`
	}{}
	key.ID.SetValue(10)

	patch := &User{}
	patch.Name.SetValue("Alice Cooper")

	upd := pg.NewPgUpdate(patch)
	if err := crudifier.PrepareUpdateModel(key, upd); err != nil {
		return "", nil, err
	}

	params := make([]any, 0)
	sql := upd.SQL(&params)
	return sql, params, nil
}
```

Typical result:

```sql
UPDATE users SET name=$1 WHERE id = $2
```

As of the current patched version, `PrepareUpdateModel()` returns an error if there are no fields to update.

## Detail select

```go
package example

import (
	"github.com/dronm/crudifier"
	"github.com/dronm/crudifier/fields"
	"github.com/dronm/crudifier/pg"
)

func buildDetailSelect() (string, []any, error) {
	key := &struct {
		ID fields.FieldInt `json:"id"`
	}{}
	key.ID.SetValue(10)

	filter := &pg.PgFilters{}
	sel := pg.NewPgDetailSelect(&User{}, filter)

	if err := crudifier.PrepareFetchModel(key, sel); err != nil {
		return "", nil, err
	}

	params := make([]any, 0)
	sql := sel.SQL(&params)
	return sql, params, nil
}
```

Typical result:

```sql
SELECT id,name,email FROM users WHERE id = $1
```

## Collection select

`CollectionParams` is the transport structure for filters, sorters, and pagination.

Example JSON:

```json
{
	"filter": [
		{
			"j": "and",
			"f": {
				"name": { "o": "ilk", "v": "%alice%" },
				"age": { "o": "ge", "v": 18 }
			}
		}
	],
	"sorter": [
		{ "f": "name", "d": "a" }
	],
	"from": 0,
	"count": 50
}
```

Usage:

```go
package example

import (
	"github.com/dronm/crudifier"
	"github.com/dronm/crudifier/pg"
)

func buildCollectionSelect(params crudifier.CollectionParams) (string, string, []any, error) {
	filter := &pg.PgFilters{}
	sorter := &pg.PgSorters{}
	limit := &pg.PgLimit{}
	sel := pg.NewPgSelect(&User{}, filter, sorter, limit)

	if err := crudifier.PrepareFetchModelCollection(sel, params); err != nil {
		return "", "", nil, err
	}

	queryParams := make([]any, 0)
	rowsSQL, aggSQL := sel.CollectionSQL(&queryParams)
	return rowsSQL, aggSQL, queryParams, nil
}
```

## Supported filter operators

Client-side operators from `CollectionParams`:

| Param | Meaning | SQL shape |
|---|---|---|
| `e` | equal | `=` |
| `ne` | not equal | `<>` |
| `l` | less than | `<` |
| `le` | less than or equal | `<=` |
| `g` | greater than | `>` |
| `ge` | greater than or equal | `>=` |
| `lk` | like | `LIKE` |
| `ilk` | ilike | `ILIKE` |
| `i` | is | `IS` |
| `in` | is not | `IS NOT` |
| `incl` | scalar in parameter array | `field = ANY($n)` |
| `any` | scalar equals any parameter array | `field = ANY($n)` |
| `has` | parameter exists in array column | `$n = ANY(field)` |
| `overlap` | array overlap | `&&` |
| `contains` | array contains | `@>` |
| `fts` | full text search | `@@ to_tsquery(...)` |

Notes:

- `le` was fixed in the patched version to generate `<=` correctly.
- `any` is rendered as `field = ANY($n)` in the PostgreSQL builder.
- full text search in `pg` currently uses PostgreSQL `to_tsquery('russian', ...)`.

## Sorting and field expressions

The patched version validates field references before they reach SQL generation.

Allowed forms are intentionally narrow:

- plain identifiers: `name`
- dotted identifiers: `users.name`
- PostgreSQL JSON extraction chains: `payload->'a'->>'b'`

Rejected input includes arbitrary SQL fragments, function calls, commas, comments, and injected clauses.

## Aggregate models

Collection queries may expose aggregate scans through `CollectionAgg()`.

Example:

```go
type UserAgg struct {
	Count int64 `json:"count" agg:"count(*)"`
}

type User struct {}

func (User) Relation() string {
	return "users"
}

func (User) CollectionAgg() any {
	return &UserAgg{}
}
```

`PrepareFetchModelCollection()` adds aggregate scan targets and `PgSelect.CollectionSQL()` returns two queries:

- collection query;
- aggregate query.

## Safety notes

### Relation trust boundary

`Relation()` is inserted into SQL as-is. That means it is suitable for trusted table names, views, or prebuilt relation fragments under your control.

Do not feed user input into `Relation()`.

### Field expression validation

The patched version adds strict field-reference sanitization in the PostgreSQL package and in collection parameter parsing. This closes the most obvious identifier/path injection cases for filters, sorters, insert fields, and update fields.

### Empty write guards

The patched version rejects:

- insert preparation with no actual insert fields;
- update preparation with no actual update fields.

This prevents invalid SQL such as:

```sql
INSERT INTO users () VALUES ()
```

or

```sql
UPDATE users SET  WHERE id = $1
```

## Error handling

Validation errors are returned as `*crudifier.ValidationError` with a combined message string.

Examples of preparation-time failures:

- unknown field in metadata;
- invalid field expression;
- invalid field value by validator;
- empty insert payload;
- empty update payload;
- missing key fields for detail fetch.

## Current limitations

- PostgreSQL is the only SQL renderer included.
- `Relation()` is still a trusted raw SQL fragment.
- validation errors are string-aggregated, not structured.
- the package builds SQL but does not execute queries or scan rows automatically.
- full text search behavior is PostgreSQL-specific.

## Development

Run tests:

```bash
go test ./...
```

Format:

```bash
gofmt -w .
```

## Patch notes for this version

This README reflects the patched variant that includes:

- correct `<=` mapping for collection filters;
- working PostgreSQL `ANY` rendering;
- field-expression sanitization for SQL builders and parameter parsing;
- explicit empty insert/update guards;
- regression tests for the fixes above.

## License

No license file is included in this package snapshot. Add one explicitly before publishing or redistributing.
