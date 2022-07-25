package buildsqlx

import (
	"testing"
)

func TestDB_Query(t *testing.T) {
	// Selects, Ordering, Limit & Offset
	qDb := db.Table("posts").Select("title", "body")
	// If you already have a query builder instance and you wish to add a column to its existing select clause, you may use the addSelect method:
	query, values := qDb.AddSelect("points").GroupBy("topic").OrderBy("points", "DESC").Limit(15).Offset(5).Query()

	t.Logf("Selects, Ordering, Limit & Offset query: %v, values:%+v", query, values)

	// InRandomOrder
	query, values = db.Table("users").Select("name", "post", "user_id").InRandomOrder().Query()
	t.Logf("InRandomOrder query: %v, values:%+v", query, values)

	// GroupBy / Having
	query, values = db.Table("users").GroupBy("account_id").Having("account_id", OpGT, 100).Query()
	t.Logf("GroupBy / Having query: %v, values:%+v", query, values)
	// Where, AndWhere, OrWhere clauses
	query, values = db.Table("table1").Select("foo", "bar", "baz").Where("foo", OpEQ, "cmp").AndWhere("bar", OpNEQ, "foo").OrWhere("baz", OpEQ, 123).Query()
	t.Logf("Where, AndWhere, OrWhere clauses query: %v, values:%+v", query, values)
	// WhereIn / WhereNotIn
	query, values = db.Table("table1").WhereIn("id", 1, 2, 3).OrWhereIn("name", "John", "Paul").Query()
	t.Logf("WhereIn / WhereNotIn query: %v, values:%+v", query, values)
	// 	WhereNull / WhereNotNull
	query, values = db.Table("posts").WhereNull("points").OrWhereNotNull("title").Query()
	t.Logf("WhereNull / WhereNotNull query: %v, values:%+v", query, values)
	// 	Left / Right / Cross / Inner / Left Outer Joins
	query, values = db.Table("users").Select("name", "post", "user_id").LeftJoin("posts", "users.id", "=", "posts.user_id").Query()
	t.Logf("WhereNull / WhereNotNull query: %v, values:%+v", query, values)
	// 	WhereBetween
	query, values = db.Table("users").Select("name").WhereBetween("points", 1233, 12345).Query()
	t.Logf("WhereBetween query: %v, values:%+v", query, values)
	// 	Union / Union All
	union := db.Table("posts").Select("title", "likes").Union()
	query, values = union.Table("users").Select("name", "points").Query()
	t.Logf("Union / Union All query: %v, values:%+v", query, values)
	// 	Determining If Records Exist
	query, values = db.Table("user").Select("name").Where("points", OpGTE, int64(12345)).Exists()
	t.Logf("Determining If Records Exist query: %v, values:%+v", query, values)

}
