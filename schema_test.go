package buildsqlx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const TableToCreate = "big_tbl"

var (
	db = NewConnection("mysql").DB()
)

func TestDB_CreateTable(t *testing.T) {
	type args struct {
		tblName string
		fn      func(table *Table) error
	}
	tests := []struct {
		name    string
		r       *DB
		args    args
		wantSql []string
		wantErr bool
	}{
		{
			name: "全类型测试",
			r:    db,
			args: args{
				tblName: TableToCreate,
				fn: func(table *Table) error {
					table.Increments("increments")
					table.Boolean("boolean")
					// table.BigIncrements("bigincrements")
					table.SmallInt("smallint")
					table.MediumInt("mediumint")
					table.Integer("integer")
					table.BigInt("bigint")
					table.Float("float")
					table.Double("double")
					table.Decimal("decimal", 6, 2)
					table.Date("date")
					table.Time("time")
					table.Year("year")
					table.DateTime("datetime")
					table.Timestamp("timestamp", true)
					table.Timestamp("timestamp1", false)
					table.Char("char", 10)
					table.String("string", 20)
					table.Text("text")
					table.Blob("blob")
					table.LongText("longtext")
					table.LongBlob("longblob")
					table.Json("json")
					return nil
				},
			},
		},
		{
			name: "两个自增主键",
			r:    db,
			args: args{
				tblName: TableToCreate,
				fn: func(table *Table) error {
					table.Increments("increments")
					table.BigIncrements("bigincrements")
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "索引",
			r:    db,
			args: args{
				tblName: TableToCreate,
				fn: func(table *Table) error {
					table.BigIncrements("bigincrements")
					table.String("index", 6).Index("idx_aaa")
					return nil
				},
			},
		},
		{
			name: "外键",
			r:    db,
			args: args{
				tblName: TableToCreate,
				fn: func(table *Table) error {
					table.BigIncrements("bigincrements")
					table.BigInt("user_id").ForeignKey("fx_aaa", "user", "id", nil, nil)
					return nil
				},
			},
		},
		{
			name: "默认值",
			r:    db,
			args: args{
				tblName: TableToCreate,
				fn: func(table *Table) error {
					table.Increments("increments")
					table.Boolean("boolean").Default(false)
					table.SmallInt("smallint").Default(1)
					table.MediumInt("mediumint").Default(0)
					table.Integer("integer").Default(12)
					table.BigInt("bigint").Default(44)
					table.Float("float").Default(1.5)
					table.Double("double").Default(1.8)
					table.Decimal("decimal", 6, 2).Default(12.00)
					table.Date("date").Default("2012-01-01")
					table.Time("time").Default("10:10:01")
					table.Year("year").Default(2014)
					table.DateTime("datetime").Default("2012-01-01 10:10:01")
					table.Timestamp("timestamp", true)
					table.Timestamp("timestamp1", false)
					table.Char("char", 10).Default("aaa")
					table.String("string", 20).Default("哈哈哈")
					table.Text("text").Default("hhh哈哈哈")
					table.Blob("blob").Default("1212")
					table.LongText("longtext").Default("?///ddd")
					table.LongBlob("longblob").Default("sss...")
					table.Json("json").Default("{}")
					return nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSql, err := tt.r.CreateTable(tt.args.tblName, tt.args.fn)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				t.Logf("%v", gotSql)
				t.Logf("output %v", tt.wantSql)
			}
		})
	}
}

func TestDB_ModifyTable(t *testing.T) {
	type args struct {
		tblName string
		fn      func(table *Table) error
	}
	tests := []struct {
		name    string
		r       *DB
		args    args
		wantSql []string
		wantErr bool
	}{
		{
			name: "全类型测试",
			r:    db,
			args: args{
				tblName: TableToCreate,
				fn: func(table *Table) error {
					table.Increments("increments").Change()
					table.Boolean("boolean").Change()
					// table.BigIncrements("bigincrements")
					table.SmallInt("smallint").Change()
					table.MediumInt("mediumint").Change()
					table.Integer("integer").Change()
					table.BigInt("bigint").Change()
					table.Float("float").Change()
					table.Double("double").Change()
					table.Decimal("decimal", 6, 2).Change()
					table.Date("date").Change()
					table.Time("time").Change()
					table.Year("year").Change()
					table.DateTime("datetime").Change()
					table.Timestamp("timestamp", true).Change()
					table.Timestamp("timestamp1", false).Change()
					table.Char("char", 10).Change()
					table.String("string", 20).Change()
					table.Text("text").Change()
					table.Blob("blob").Change()
					table.LongText("longtext").Change()
					table.LongBlob("longblob").Change()
					table.Json("json").Change()
					return nil
				},
			},
		},
		{
			name: "添加/变更/删除测试",
			r:    db,
			args: args{
				tblName: TableToCreate,
				fn: func(table *Table) error {
					table.Increments("increments").Change()
					table.Boolean("boolean").Change()
					// table.BigIncrements("bigincrements")
					table.SmallInt("smallint").Index("idx_xxx")
					table.MediumInt("mediumint")
					table.DropColumn("integer")
					table.DropColumn("bigint")
					table.DropIndex("idx_xssa")
					table.Float("float").Change()
					table.Double("double").Change()
					table.Decimal("decimal", 6, 2).Change()
					table.Date("date").Change()
					table.Time("time").Change()
					table.Year("year").Change()
					table.DateTime("datetime").Change()
					table.Timestamp("timestamp", true).Change()
					table.Timestamp("timestamp1", false).Change()
					table.Char("char", 10).Change()
					table.String("string", 20).Change()
					table.Text("text").Change()
					table.Blob("blob").Change()
					table.LongText("longtext").Change()
					table.LongBlob("longblob").Change()
					table.Json("json").Change()
					return nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSql, err := tt.r.ModifyTable(tt.args.tblName, tt.args.fn)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				t.Logf("%v", gotSql)
				t.Logf("output %v", tt.wantSql)
			}
		})
	}
}
