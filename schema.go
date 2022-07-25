package buildsqlx

import (
	"errors"
	"strconv"
	"strings"

	"github.com/spf13/cast"
)

var (
	errTableOnlySupportOneAutoIncrements = errors.New("sql: the table only support one increments column")
)

// mysql column types
const (
	TypeTinyInt      = "TINYINT"
	TypeSmallInt     = "SMALLINT"
	TypeMediumInt    = "MEDIUMINT"
	TypeInt          = "INTEGER"
	TypeBigInt       = "BIGINT"
	TypeFloat        = "FLOAT"
	TypeDouble       = "DOUBLE"
	TypeDecimal      = "DECIMAL"
	TypeDate         = "DATE"
	TypeTime         = "TIME"
	TypeYear         = "YEAR"
	TypeDateTime     = "DATETIME"
	TypeTimestamp    = "TIMESTAMP"
	CurrentTimestamp = "CURRENT_TIMESTAMP"
	TypeChar         = "CHAR"
	TypeVarchar      = "VARCHAR"
	TypeBlob         = "BLOB"
	TypeText         = "TEXT"
	TypeLongBlob     = "LONGBLOB"
	TypeLongText     = "LONGTEXT"
	TypeJson         = "JSON"
)

// specific for Mysql driver
const (
	SemiColon  = ";"
	AlterTable = "ALTER TABLE "
	Add        = " ADD "
	Modify     = " ALTER "
	Drop       = " DROP "
	Rename     = " RENAME "
)

type colType string

// Table is the type for operations on table schema
type Table struct {
	columns []*column
	tblName string
	comment *string
	sb      *schemaBuilder
}

// collection of properties for the column
type column struct {
	Name          string
	RenameTo      *string
	IsNotNull     *bool
	AutoIncrement bool
	IsPrimaryKey  bool
	ColumnType    colType
	Default       *string
	IsIndex       bool
	IsUnique      bool
	ForeignKey    *string
	IdxName       string
	Comment       *string
	IsDrop        bool
	IsModify      bool
	After         *string
	ChartSet      *string
	Collation     *string
	Op            string
}

// CreateTable creates and/or manipulates table structure with an appropriate types/indices/comments/defaults/nulls etc
func (r *DB) CreateTable(tblName string, fn func(table *Table) error) (sql []string, err error) {
	tbl := &Table{tblName: tblName, sb: newSchemaBuilder()}
	err = fn(tbl) // run fn with Table struct passed to collect columns to []*column slice
	if err != nil {
		return nil, err
	}

	l := len(tbl.columns)
	if l > 0 {
		// create table with relative columns/indices
		return r.createTable(tbl)
	}

	return
}

// ModifyTable creates and/or manipulates table structure with an appropriate types/indices/comments/defaults/nulls etc
func (r *DB) ModifyTable(tblName string, fn func(table *Table) error) (sql []string, err error) {
	tbl := &Table{tblName: tblName, sb: newSchemaBuilder()}
	err = fn(tbl) // run fn with Table struct passed to collect columns to []*column slice
	if err != nil {
		return nil, err
	}

	l := len(tbl.columns)
	if l > 0 {
		return r.modifyTable(tbl)
	}

	return
}

// Increments creates auto incremented primary key integer column
func (t *Table) Increments(colNm string) *Table {
	t.columns = append(t.columns, &column{Name: colNm, ColumnType: TypeInt, IsPrimaryKey: true, AutoIncrement: true})
	return t
}

// Boolean creates boolean type column
func (t *Table) Boolean(colNm string) *Table {
	t.columns = append(t.columns, &column{Name: colNm, ColumnType: TypeTinyInt})
	return t
}

// BigIncrements creates auto incremented primary key big integer column
func (t *Table) BigIncrements(colNm string) *Table {
	t.columns = append(t.columns, &column{Name: colNm, ColumnType: TypeBigInt, IsPrimaryKey: true, AutoIncrement: true})
	return t
}

// SmallInt creates small integer column
func (t *Table) SmallInt(colNm string) *Table {
	t.columns = append(t.columns, &column{Name: colNm, ColumnType: TypeSmallInt})
	return t
}

// MediumInt creates medium integer column
func (t *Table) MediumInt(colNm string) *Table {
	t.columns = append(t.columns, &column{Name: colNm, ColumnType: TypeMediumInt})
	return t
}

// Integer creates an integer column
func (t *Table) Integer(colNm string) *Table {
	t.columns = append(t.columns, &column{Name: colNm, ColumnType: TypeInt})
	return t
}

// BigInt creates big integer column
func (t *Table) BigInt(colNm string) *Table {
	t.columns = append(t.columns, &column{Name: colNm, ColumnType: TypeBigInt})
	return t
}

// Float creates float column
func (t *Table) Float(colNm string) *Table {
	t.columns = append(t.columns, &column{Name: colNm, ColumnType: TypeFloat})
	return t
}

// Double creates double column
func (t *Table) Double(colNm string) *Table {
	t.columns = append(t.columns, &column{Name: colNm, ColumnType: TypeDouble})
	return t
}

// Decimal alias for Numeric as for PostgreSQL they are the same
func (t *Table) Decimal(colNm string, precision, scale uint64) *Table {
	t.columns = append(t.columns, &column{Name: colNm, ColumnType: colType(TypeDecimal + "(" + strconv.FormatUint(precision, 10) + ", " + strconv.FormatUint(scale, 10) + ")")})
	return t
}

// Date	creates date column with an ability to set current_date as default value
func (t *Table) Date(colNm string) *Table {
	t.columns = append(t.columns, &column{Name: colNm, ColumnType: TypeDate})
	return t
}

// Time creates time column with an ability to set current_time as default value
func (t *Table) Time(colNm string) *Table {
	t.columns = append(t.columns, &column{Name: colNm, ColumnType: TypeTime})
	return t
}

// Year creates year column
func (t *Table) Year(colNm string) *Table {
	t.columns = append(t.columns, &column{Name: colNm, ColumnType: TypeYear})
	return t
}

// DateTime creates datetime column with an ability to set NOW() as default value
func (t *Table) DateTime(colNm string) *Table {
	t.columns = append(t.columns, &column{Name: colNm, ColumnType: TypeDateTime})
	return t
}

// Timestamp creates timestamp column with an ability to set NOW() as default value
func (t *Table) Timestamp(colNm string, isDefault bool) *Table {
	t.columns = append(t.columns, buildDateTIme(colNm, TypeTimestamp, CurrentTimestamp, isDefault))
	return t
}

// Char creates char(len) column
func (t *Table) Char(colNm string, len uint64) *Table {
	t.columns = append(t.columns, &column{Name: colNm, ColumnType: colType(TypeChar + "(" + strconv.FormatUint(len, 10) + ")")})
	return t
}

// String creates varchar(len) column
func (t *Table) String(colNm string, len uint64) *Table {
	t.columns = append(t.columns, &column{Name: colNm, ColumnType: colType(TypeVarchar + "(" + strconv.FormatUint(len, 10) + ")")})
	return t
}

// Text	creates text type column
func (t *Table) Text(colNm string) *Table {
	t.columns = append(t.columns, &column{Name: colNm, ColumnType: TypeText})
	return t
}

// Blob	creates blob type column
func (t *Table) Blob(colNm string) *Table {
	t.columns = append(t.columns, &column{Name: colNm, ColumnType: TypeBlob})
	return t
}

// LongText	creates long text type column
func (t *Table) LongText(colNm string) *Table {
	t.columns = append(t.columns, &column{Name: colNm, ColumnType: TypeLongText})
	return t
}

// LongBlob	creates long blob type column
func (t *Table) LongBlob(colNm string) *Table {
	t.columns = append(t.columns, &column{Name: colNm, ColumnType: TypeLongBlob})
	return t
}

// Json creates json text typed column
func (t *Table) Json(colNm string) *Table {
	t.columns = append(t.columns, &column{Name: colNm, ColumnType: TypeJson})
	return t
}

// NotNull sets the last column to not null
func (t *Table) NotNull() *Table {
	isNotNull := true
	t.columns[len(t.columns)-1].IsNotNull = &isNotNull
	return t
}

// Collation sets the last column to specified collation
func (t *Table) Collation(coll string) *Table {
	t.columns[len(t.columns)-1].Collation = &coll
	return t
}

// Collation sets the last column to specified collation
func (t *Table) After(coll string) *Table {
	t.columns[len(t.columns)-1].After = &coll
	return t
}

// Default sets the default column value
func (t *Table) Default(val interface{}) *Table {
	v := cast.ToString(val)
	t.columns[len(t.columns)-1].Default = &v
	return t
}

// Comment sets the column comment
func (t *Table) Comment(cmt string) *Table {
	t.columns[len(t.columns)-1].Comment = &cmt
	return t
}

// TableComment sets the comment for table
func (t *Table) TableComment(cmt string) {
	t.comment = &cmt
}

// Index sets the last column to btree index
func (t *Table) Index(idxName string) *Table {
	t.columns[len(t.columns)-1].IdxName = idxName
	t.columns[len(t.columns)-1].IsIndex = true
	return t
}

// Unique sets the last column to unique index
func (t *Table) Unique(idxName string) *Table {
	t.columns[len(t.columns)-1].IdxName = idxName
	t.columns[len(t.columns)-1].IsUnique = true
	return t
}

// ForeignKey sets the last column to reference rfcTbl on onCol with idxName foreign key index
func (t *Table) ForeignKey(idxName, rfcTbl, onCol string, update, delete *string) *Table {

	query := strings.Builder{}
	query.WriteString(" CONSTRAINT ")
	query.WriteByte('`')
	query.WriteString(idxName)
	query.WriteByte('`')
	query.WriteString(" FOREIGN KEY ( ")
	query.WriteByte('`')
	query.WriteString(t.columns[len(t.columns)-1].Name)
	query.WriteByte('`')
	query.WriteString(") REFERENCES ")
	query.WriteByte('`')
	query.WriteString(rfcTbl)
	query.WriteByte('`')
	query.WriteString(" (")
	query.WriteByte('`')
	query.WriteString(onCol)
	query.WriteByte('`')
	query.WriteString(")")

	if update != nil {
		query.WriteString(" ON UPDATE ")
		query.WriteString(*update)
	} else {
		query.WriteString(" ON UPDATE NO ACTION ")
	}

	if delete != nil {
		query.WriteString(" ON DELETE ")
		query.WriteString(*delete)
	} else {
		query.WriteString(" ON DELETE NO ACTION ")
	}

	key := query.String()
	t.columns[len(t.columns)-1].ForeignKey = &key
	return t
}

// build any date/time type with defaults preset
func buildDateTIme(colNm, t, defType string, isDefault bool) *column {
	isNotNull := false
	if isDefault {
		isNotNull = true
	}
	col := &column{Name: colNm, ColumnType: colType(t), IsNotNull: &isNotNull}
	if isDefault {
		col.Default = &defType
	}
	return col
}

// Change the column type/length/nullable etc options
func (t *Table) Change() {
	t.columns[len(t.columns)-1].IsModify = true
}

// Rename the column "from" to the "to"
func (t *Table) Rename(from, to string) *Table {
	t.columns = append(t.columns, &column{Name: from, RenameTo: &to, IsModify: true})
	return t
}

// DropColumn the column named colNm in this table context
func (t *Table) DropColumn(colNm string) {
	t.columns = append(t.columns, &column{Name: colNm, IsDrop: true})
}

// DropIndex the column named idxNm in this table context
func (t *Table) DropIndex(idxNm string) {
	t.columns = append(t.columns, &column{IdxName: idxNm, IsDrop: true, IsIndex: true})
}

// createTable create table with relative columns/indices
func (r *DB) createTable(t *Table) (sql []string, err error) {
	l := len(t.columns)
	autoIncr := 0

	t.sb.WriteString("CREATE TABLE")
	t.sb.Pad().Ident(t.tblName)
	t.sb.Nested(func(sb *schemaBuilder) {

		for k, col := range t.columns {
			sb.Ident(col.Name).Pad()
			sb.WriteString(string(col.ColumnType))

			// 自增
			if col.AutoIncrement {
				autoIncr++
				sb.Pad().WriteString("AUTO_INCREMENT")
			}
			// 字符集
			if col.ChartSet != nil {
				sb.Pad().WriteString("CHARACTER SET").Pad().WriteByte('\'').WriteString(*col.ChartSet).WriteByte('\'')
			}
			// Collation
			if col.Collation != nil {
				sb.Pad().WriteString("COLLATE").Pad().WriteByte('\'').WriteString(*col.Collation).WriteByte('\'')
			}
			// 不为空
			if col.IsNotNull != nil {
				if *col.IsNotNull {
					sb.Pad().WriteString("NOT NULL")
				} else {
					sb.Pad().WriteString("NULL")
				}
			}
			// 默认值
			if col.Default != nil {

				if strings.HasPrefix(string(col.ColumnType), TypeChar) {
					sb.Pad().WriteString("DEFAULT").Pad().WriteByte('\'').WriteString(*col.Default).WriteByte('\'')

				} else if strings.HasPrefix(string(col.ColumnType), TypeVarchar) {
					sb.Pad().WriteString("DEFAULT").Pad().WriteByte('\'').WriteString(*col.Default).WriteByte('\'')
				} else {
					switch colType(col.ColumnType) {
					case TypeDate, TypeTime, TypeDateTime:
						sb.Pad().WriteString("DEFAULT").Pad().WriteByte('\'').WriteString(*col.Default).WriteByte('\'')
					case TypeBlob, TypeLongBlob, TypeText, TypeLongText, TypeJson:
						// do nothing
					default:
						sb.Pad().WriteString("DEFAULT").Pad().WriteString(*col.Default)
					}
				}
			}
			// 备注
			if col.Comment != nil {
				sb.Pad().WriteString("COMMENT").Pad().WriteByte('\'').WriteString(*col.Comment).WriteByte('\'')
			}

			if k < l-1 {
				sb.Comma()
			}

			// 主键
			if col.IsPrimaryKey {
				sb.child.Comma().WriteString("PRIMARY KEY").Nested(func(csb *schemaBuilder) {
					csb.Ident(col.Name)
				})
			}

			// 索引
			if col.IsIndex {
				sb.child.Comma().Pad().WriteString("INDEX").Pad().Ident(col.IdxName).Pad().Nested(func(csb *schemaBuilder) {
					csb.Ident(col.Name)
					csb.Pad()
					csb.WriteString("ASC")
				})
			}
			// 唯一索引
			if col.IsUnique {
				sb.child.Comma().Pad().WriteString("UNIQUE INDEX").Pad().Ident(col.IdxName).Pad().Nested(func(csb *schemaBuilder) {
					csb.Ident(col.Name)
					csb.Pad()
					csb.WriteString("ASC")
				})
			}
			// 外键
			if col.ForeignKey != nil {
				sb.child.Comma().WriteString(*col.ForeignKey)
			}
		}

		sb.WriteString(sb.child.String())
	})

	if autoIncr > 1 {
		return nil, errTableOnlySupportOneAutoIncrements
	}

	if t.comment != nil {
		t.sb.WriteString("COMMENT").Pad().WriteByte('\'').WriteString(*t.comment).WriteByte('\'')
	}

	sql = append(sql, t.sb.String())
	return
}

// adds, modifies or deletes column
func (r *DB) modifyTable(t *Table) (sql []string, err error) {
	l := len(t.columns)

	t.sb.WriteString("ALTER TABLE")
	t.sb.Pad().Ident(t.tblName).Pad()
	for k, col := range t.columns {

		if col.IsDrop {
			if col.IsIndex {
				// 删除索引
				t.sb.WriteString("DROP INDEX").Pad().Ident(col.IdxName)
			} else {
				// 字段删除
				t.sb.WriteString("DROP COLUMN").Pad().Ident(col.Name)
			}
		} else if col.IsModify {
			t.sb.WriteString("CHANGE COLUMN")
			t.sb.Pad().Ident(col.Name)
			// 改名
			if col.RenameTo != nil {
				t.sb.Pad().Ident(*col.RenameTo)
			} else {
				t.sb.Pad().Ident(col.Name)
			}
			t.sb.Pad().WriteString(string(col.ColumnType))
			// 自增
			if col.AutoIncrement {
				t.sb.Pad().WriteString("AUTO_INCREMENT")
			}
			// 字符集
			if col.ChartSet != nil {
				t.sb.Pad().WriteString("CHARACTER SET").Pad().WriteByte('\'').WriteString(*col.ChartSet).WriteByte('\'')
			}
			// Collation
			if col.Collation != nil {
				t.sb.Pad().WriteString("COLLATE").Pad().WriteByte('\'').WriteString(*col.Collation).WriteByte('\'')
			}
			// 不为空
			if col.IsNotNull != nil {
				if *col.IsNotNull {
					t.sb.Pad().WriteString("NOT NULL")
				} else {
					t.sb.Pad().WriteString("NULL")
				}
			}
			// 默认值
			if col.Default != nil {
				switch colType(col.ColumnType) {
				case TypeChar, TypeVarchar, TypeDate, TypeTime, TypeDateTime:
					t.sb.Pad().WriteString("DEFAULT").Pad().WriteByte('\'').WriteString(*col.Default).WriteByte('\'')
				case TypeBlob, TypeLongBlob, TypeText, TypeLongText, TypeJson:
					// do nothing
				default:
					t.sb.Pad().WriteString("DEFAULT").Pad().WriteString(*col.Default)
				}
			}
			// 备注
			if col.Comment != nil {
				t.sb.Pad().WriteString("COMMENT").Pad().WriteByte('\'').WriteString(*col.Comment).WriteByte('\'')
			}
		} else {
			// 添加字段
			t.sb.WriteString("ADD COLUMN")
			t.sb.Pad().Ident(col.Name)
			t.sb.Pad().WriteString(string(col.ColumnType))
			// 自增
			if col.AutoIncrement {
				t.sb.Pad().WriteString("AUTO_INCREMENT")
			}
			// 字符集
			if col.ChartSet != nil {
				t.sb.Pad().WriteString("CHARACTER SET").Pad().WriteByte('\'').WriteString(*col.ChartSet).WriteByte('\'')
			}
			// Collation
			if col.Collation != nil {
				t.sb.Pad().WriteString("COLLATE").Pad().WriteByte('\'').WriteString(*col.Collation).WriteByte('\'')
			}
			// 不为空
			if col.IsNotNull != nil {
				if *col.IsNotNull {
					t.sb.Pad().WriteString("NOT NULL")
				} else {
					t.sb.Pad().WriteString("NULL")
				}
			}
			// 默认值
			if col.Default != nil {
				switch colType(col.ColumnType) {
				case TypeChar, TypeVarchar, TypeDate, TypeTime, TypeDateTime:
					t.sb.Pad().WriteString("DEFAULT").Pad().WriteByte('\'').WriteString(*col.Default).WriteByte('\'')
				case TypeBlob, TypeLongBlob, TypeText, TypeLongText, TypeJson:
					// do nothing
				default:
					t.sb.Pad().WriteString("DEFAULT").Pad().WriteString(*col.Default)
				}
			}
			// 备注
			if col.Comment != nil {
				t.sb.Pad().WriteString("COMMENT").Pad().WriteByte('\'').WriteString(*col.Comment).WriteByte('\'')
			}
			// After,默认添加到after之后
			if col.After != nil {
				t.sb.Pad().WriteString("AFTER").Pad().Ident(*col.After)
			} else {
				t.sb.Pad().WriteString("AFTER").Pad().Ident("id")
			}
		}

		if k < l-1 {
			t.sb.Comma()
		}

		if !col.IsDrop {
			// 索引
			if col.IsIndex {
				t.sb.child.Comma().
					Pad().WriteString("ADD INDEX").
					Pad().Ident(col.IdxName).Pad().
					Nested(func(csb *schemaBuilder) {
						csb.Ident(col.Name)
						csb.Pad()
						csb.WriteString("ASC")
					})
			}
			// 唯一索引
			if col.IsUnique {
				t.sb.child.Comma().
					Pad().WriteString("ADD UNIQUE INDEX").
					Pad().Ident(col.IdxName).Pad().
					Nested(func(csb *schemaBuilder) {
						csb.Ident(col.Name)
						csb.Pad()
						csb.WriteString("ASC")
					})
			}
			// 外键
			if col.ForeignKey != nil {
				t.sb.child.Comma().WriteString("ADD").Pad().WriteString(*col.ForeignKey)
			}
		}
	}

	t.sb.WriteString(t.sb.child.String())

	if t.comment != nil {
		t.sb.WriteString("COMMENT").Pad().WriteByte('\'').WriteString(*t.comment).WriteByte('\'')
	}

	sql = append(sql, t.sb.String())
	return
}
