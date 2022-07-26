# buildsqlx
Go Database query builder library [![Tweet](http://jpillora.com/github-twitter-button/img/tweet.png)](https://twitter.com/intent/tweet?text=Go%20database%20query%20builder%20library%20&url=https://github.com/go-kenka/buildsqlx&hashtags=go,golang,sql,builder,mysql,sql-builder,developers)

[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-kenka/buildsqlx)](https://goreportcard.com/report/github.com/go-kenka/buildsqlx)
[![Build and run](https://github.com/go-kenka/buildsqlx/workflows/Build%20and%20run/badge.svg)](https://github.com/go-kenka/buildsqlx/actions)
[![GoDoc](https://github.com/golang/gddo/blob/c782c79e0a3c3282dacdaaebeff9e6fd99cb2919/gddo-server/assets/status.svg)](https://godoc.org/github.com/go-kenka/buildsqlx)
[![codecov](https://codecov.io/gh/arthurkushman/buildsqlx/branch/master/graph/badge.svg)](https://codecov.io/gh/arthurkushman/buildsqlx)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

* [Installation](#user-content-installation)
* [Selects, Ordering, Limit & Offset](#user-content-selects-ordering-limit--offset)
* [GroupBy / Having](#user-content-groupby--having)
* [Where, AndWhere, OrWhere clauses](#user-content-where-andwhere-orwhere-clauses)
* [WhereIn / WhereNotIn](#user-content-wherein--wherenotin)
* [WhereNull / WhereNotNull](#user-content-wherenull--wherenotnull)
* [Left / Right / Cross / Inner / Left Outer Joins](#user-content-left--right--cross--inner--left-outer-joins)
* [Inserts](#user-content-inserts)
* [Updates](#user-content-updates)
* [Delete](#user-content-delete)
* [Drop, Truncate, Rename](#user-content-drop-truncate-rename)
* [Union / Union All](#user-content-union--union-all)
* [WhereExists / WhereNotExists](#user-content-whereexists--wherenotexists)
* [Determining If Records Exist](#user-content-determining-if-records-exist)
* [Aggregates](#user-content-aggregates)
* [Create table](#user-content-create-table)
* [Add / Modify / Drop columns](#user-content-add--modify--drop-columns)

## Installation
```bash
go get -u github.com/go-kenka/buildsqlx
```

## Selects, Ordering, Limit & Offset

You may not always want to select all columns from a database table. Using the select method, you can specify a custom select clause for the query:

```go
package yourpackage

import (
	"github.com/go-kenka/buildsqlx"
) 

var db = buildsqlx.NewDb(buildsqlx.NewConnection("mysql"))

func main() {
    qDb := db.Table("posts").Select("title", "body")

    // If you already have a query builder instance and you wish to add a column to its existing select clause, you may use the addSelect method:
    query, values := qDb.AddSelect("points").GroupBy("topic").OrderBy("points", "DESC").Limit(15).Offset(5).Query()
}
```

### InRandomOrder
```go
query, values := db.Table("users").Select("name", "post", "user_id").InRandomOrder().Query()
```

## GroupBy / Having
The GroupBy and Having methods may be used to group the query results. 
The having method's signature is similar to that of the where method:
```go
query, values := db.table("users").GroupBy("account_id").Having("account_id", ">", 100).Query()
```

## Where, AndWhere, OrWhere clauses
You may use the where method on a query builder instance to add where clauses to the query. 
The most basic call to where requires three arguments. 
The first argument is the name of the column. 
The second argument is an operator, which can be any of the database's supported operators. 
Finally, the third argument is the value to evaluate against the column.

```go
package yourpackage

import (
	"github.com/go-kenka/buildsqlx"
)

func main() {
    query, values := db.Table("table1").Select("foo", "bar", "baz").Where("foo", "=", cmp).AndWhere("bar", "!=", "foo").OrWhere("baz", "=", 123).Query()
}
```

You may chain where constraints together as well as add or clauses to the query. 
The orWhere method accepts the same arguments as the where method.

## WhereIn / WhereNotIn 
The whereIn method verifies that a given column's value is contained within the given slice:
```go
query, values := db.Table("table1").WhereIn("id", []int64{1, 2, 3}).OrWhereIn("name", []string{"John", "Paul"}).Query()
```

## WhereNull / WhereNotNull  
The whereNull method verifies that the value of the given column is NULL:
```go
query, values := db.Table("posts").WhereNull("points").OrWhereNotNull("title").Query()
```

## Left / Right / Cross / Inner / Left Outer Joins
The query builder may also be used to write join statements. 
To perform a basic "inner join", you may use the InnerJoin method on a query builder instance. 
The first argument passed to the join method is the name of the table you need to join to, 
while the remaining arguments specify the column constraints for the join. 
You can even join to multiple tables in a single query:
```go
query, values := db.Table("users").Select("name", "post", "user_id").LeftJoin("posts", "users.id", "=", "posts.user_id").Query()
```

## Inserts
The query builder also provides an insert method for inserting records into the database table. 
The insert method accepts a map of column names and values:

```go
package yourpackage

import (
	"github.com/go-kenka/buildsqlx"
)

func main() {
    // insert without getting id
    query, values := db.Table("table1").Insert(map[string]interface{}{"foo": "foo foo foo", "bar": "bar bar bar", "baz": int64(123)})

    // batch insert 
    query, values := db.Table("table1").InsertBatch([]map[string]interface{}{
                                    	0: {"foo": "foo foo foo", "bar": "bar bar bar", "baz": 123},
                                    	1: {"foo": "foo foo foo foo", "bar": "bar bar bar bar", "baz": 1234},
                                    	2: {"foo": "foo foo foo foo foo", "bar": "bar bar bar bar bar", "baz": 12345},
                                    })
}
```

## Updates
In addition to inserting records into the database, 
the query builder can also update existing records using the update method. 
The update method, like the insert method, accepts a slice of column and value pairs containing the columns to be updated. 
You may constrain the update query using where clauses:
```go
query, values := db.Table("posts").Where("points", ">", 3).Update(map[string]interface{}{"title": "awesome"})
```

## Delete
The query builder may also be used to delete records from the table via the delete method. 
You may constrain delete statements by adding where clauses before calling the delete method:
```go
query, values := db.Table("posts").Where("points", "=", 123).Delete()
```

## Drop, Truncate, Rename
```go
package yourpackage

import (
	"github.com/go-kenka/buildsqlx"
)

func main() {
    query := db.Drop("table_name")

    query := db.DropIfExists("table_name")

    query := db.Truncate("table_name")

    query := db.Rename("table_name1", "table_name2")
}
```

## Union / Union All
The query builder also provides a quick way to "union" two queries together. 
For example, you may create an initial query and use the union method to union it with a second query:
```go
union := db.Table("posts").Select("title", "likes").Union()
query, values := union.Table("users").Select("name", "points").Query()

// or if UNION ALL is of need
// union := db.Table("posts").Select("title", "likes").UnionAll()
```

## WhereBetween / WhereNotBetween
The whereBetween func verifies that a column's value is between two values:
```go
query, values := db.Table(UsersTable).Select("name").WhereBetween("points", 1233, 12345).Query()
```

The whereNotBetween func verifies that a column's value lies outside of two values:
```go
query, values := db.Table(UsersTable).Select("name").WhereNotBetween("points", 123, 123456).Query()
```

## Determining If Records Exist
Instead of using the count method to determine if any records exist that match your query's constraints, 
you may use the exists and doesntExist methods:
```go
query, values := db.Table(UsersTable).Select("name").Where("points", ">=", int64(12345)).Exists()
// use an inverse DoesntExists() if needed
```

## Aggregates
The query builder also provides a variety of aggregate methods such as Count, Max, Min, Avg, and Sum. 
You may call any of these methods after constructing your query:
```go
query, values := db.Table(UsersTable).WHere("points", ">=", 1234).Count()

query, values := db.Table(UsersTable).Avg("points")

query, values := db.Table(UsersTable).Max("points")

query, values := db.Table(UsersTable).Min("points")

query, values := db.Table(UsersTable).Sum("points")
```

## Create table
To create a new database table, use the CreateTable method. 
The Schema method accepts two arguments. 
The first is the name of the table, while the second is an anonymous function/closure which receives a Table struct that may be used to define the new table:
```go
query, values := db.Schema("big_tbl", func(table *Table) error {
    table.Increments("id")
    table.String("title", 128).Default("The quick brown fox jumped over the lazy dog").Unique("idx_ttl")
    table.SmallInt("cnt").Default(1)
    table.Integer("points").NotNull()
    table.BigInt("likes").Index("idx_likes")
    table.Text("comment").Comment("user comment").Collation("de_DE")
    table.DblPrecision("likes_to_points").Default(0.0)
    table.Char("tag", 10)
    table.DateTime("created_at", true)
    table.DateTimeTz("updated_at", true)		
    table.Decimal("tax", 2, 2)
    table.TsVector("body")
    table.TsQuery("body_query")		
    table.Jsonb("settings")
    table.Point("pt")
    table.Polygon("poly")		
    table.TableComment("big table for big data")	
	
	return nil
})

// to make a foreign key constraint from another table
query, values = db.Schema("tbl_to_ref", func(table *Table) error {
    table.Increments("id")
    table.Integer("big_tbl_id").ForeignKey("fk_idx_big_tbl_id", "big_tbl", "id")
    // to add index on existing column just repeat stmt + index e.g.:
    table.Char("tag", 10).Index("idx_tag")
    table.Rename("settings", "options")

    return nil
})	
```

## Add / Modify / Drop columns
The Table structure in the Schema's 2nd argument may be used to update existing tables. Just the way you've been created it.
The Change method allows you to modify some existing column types to a new type or modify the column's attributes.
```go
query, values := db.Schema("tbl_name", func(table *Table) error {
    table.String("title", 128).Change()

    return nil
})
```
Use DropColumn method to remove any column:
```go
query, values := db.Schema("tbl_name", func(table *Table) error {
    table.DropColumn("deleted_at")
    // To drop an index on the column    
    table.DropIndex("idx_title")

    return nil
})
```

PS Why use buildsqlx? Because it is simple and fast, yet versatile. 
The performance achieved because of structs conversion lack, as all that you need is just a columns - u can get it from an associated array/map while the conversion itself and it's processing eats more CPU/memory resources.

Supporters gratitude:

<img src="https://github.com/SoliDry/laravel-api/blob/master/tests/images/jetbrains-logo.png" alt="JetBrains logo" width="200" height="166" />