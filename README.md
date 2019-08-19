# arsqlx
Go Active Record library for postgresql

## Selects, Ordering, Grouping, Limit & Offset

```go
package yourpackage

import (
	_ "github.com/lib/pq"
) 

qDb := db.Table(TestTable).Select("foo", "bar")

res, err := qDb.AddSelect("baz").GroupBy("foo").OrderBy("bar", "DESC").Limit(15).Offset(5).Get()
```

## Where, AndWhere, OrWhere clauses

```go
package yourpackage

import (
	_ "github.com/lib/pq"
)

res, err := db.Table(TestTable).Select("foo", "bar", "baz").Where("foo", "=", cmp).AndWhere("bar", "!=", "foo").OrWhere("baz", "=", 123).Get()
```

## Inserts

```go
package yourpackage

// insert without getting id
err := db.Table("test").Insert(map[string]interface{}{"foo": "foo foo foo", "bar": "bar bar bar", "baz": int64(123)})

// insert returning id
id, err := db.Table("test").InsertGetId(map[string]interface{}{"foo": "foo foo foo", "bar": "bar bar bar", "baz": int64(123)})

// batch insert 
err := db.Table("test").InsertBatch([]map[string]interface{}{
                                    	0: {"foo": "foo foo foo", "bar": "bar bar bar", "baz": 123},
                                    	1: {"foo": "foo foo foo foo", "bar": "bar bar bar bar", "baz": 1234},
                                    	2: {"foo": "foo foo foo foo foo", "bar": "bar bar bar bar bar", "baz": 12345},
                                    })
```