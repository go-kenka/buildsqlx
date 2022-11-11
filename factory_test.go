package buildsqlx

import (
	"testing"
)

func TestDB_Insert(t *testing.T) {
	// 	Insert
	query, values := db.Table("table1").Insert(map[string]interface{}{"foo": "foo foo foo", "bar": "bar bar bar", "baz": int64(123)})
	t.Logf("Insert query: %v, values:%+v", query, values)
}

func TestDB_InsertBatch(t *testing.T) {
	// 	Insert
	query, values := db.Table("table1").InsertBatch([]map[string]interface{}{
		0: {"foo": "foo foo foo", "bar": "bar bar bar", "baz": 123},
		1: {"foo": "foo foo foo foo", "bar": "bar bar bar bar", "baz": 1234},
		2: {"foo": "foo foo foo foo foo", "bar": "bar bar bar bar bar", "baz": 12345},
	})
	t.Logf("Insert query: %v, values:%+v", query, values)
}
func TestDB_Updates(t *testing.T) {
	// 	Insert
	query, values := db.Table("posts").Where("points", OpGT, 3).Update(map[string]interface{}{"title": "awesome"})
	t.Logf("Insert query: %v, values:%+v", query, values)
}
func TestDB_Delete(t *testing.T) {
	// 	Insert
	query, values := db.Table("posts").Where("points", OpGT, 3).Delete()
	t.Logf("Insert query: %v, values:%+v", query, values)
}

func TestDB_UpdateBatch(t *testing.T) {
	// where
	where := make(map[string][]int)
	where["id"] = []int{1, 2, 3, 4}
	// where["did"] = []int{1, 2, 3, 4}
	// update
	update := make(map[string][]interface{})
	update["name"] = []interface{}{"a1", "a2", "a3", "a4"}
	update["org"] = []interface{}{"b1", "b2", "b3", "b4"}

	// 	update
	query, values := db.Table("table1").UpdateBatch(where, update)
	t.Logf("Update query: %v, values:%+v", query, values)
}
