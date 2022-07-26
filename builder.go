package buildsqlx

import (
	"log"
	"os"

	"github.com/huandu/go-clone"
)

const (
	joinInner     = "INNER"
	joinLeft      = "LEFT"
	joinRight     = "RIGHT"
	joinFull      = "FULL"
	joinFullOuter = "FULL OUTER"
	where         = " WHERE "
	and           = " AND "
	or            = " OR "
)

// inner type to build qualified sql
type builder struct {
	sqlBuilder
	where         *sqlBuilder
	table         string
	from          string
	join          []string
	orderBy       map[string]string
	orderByRaw    *string
	groupBy       string
	having        *sqlBuilder
	columns       []string
	union         []string
	isUnionAll    bool
	offset        int64
	limit         int64
	lockForUpdate *string
}

func newBuilder() *builder {
	return &builder{
		columns: []string{"*"},
	}
}

func deepClone(b *builder) *builder {
	return clone.Slowly(b).(*builder)
}

// Target returns db driver
func (r *DB) Target() string {
	return r.Conn.driver
}

// Table appends table name to sql query
func (r *DB) Table(table string) *DB {
	// reset before constructing again
	r.reset()
	r.Builder.table = table
	return r
}

// resets all builder elements to prepare them for next round
func (r *DB) reset() {
	r.Builder.sqlBuilder = sqlBuilder{}
	r.Builder.table = ""
	r.Builder.columns = []string{"*"}
	r.Builder.where = &sqlBuilder{}
	r.Builder.groupBy = ""
	r.Builder.having = &sqlBuilder{}
	r.Builder.orderBy = make(map[string]string, 0)
	r.Builder.offset = 0
	r.Builder.limit = 0
	r.Builder.join = []string{}
	r.Builder.from = ""
	// union不初始化
	// r.Builder.union = []string{}
	r.Builder.isUnionAll = false
	r.Builder.lockForUpdate = nil
	r.Builder.orderByRaw = nil
}

// Select accepts columns to select from a table
func (r *DB) Select(args ...string) *DB {
	r.Builder.columns = []string{}
	r.Builder.columns = append(r.Builder.columns, args...)
	return r
}

// OrderBy adds ORDER BY expression to SQL stmt
func (r *DB) OrderBy(column string, direction string) *DB {
	r.Builder.orderBy[column] = direction
	return r
}

// OrderByRaw adds ORDER BY raw expression to SQL stmt
func (r *DB) OrderByRaw(exp string) *DB {
	r.Builder.orderByRaw = &exp
	return r
}

// InRandomOrder add ORDER BY random() - note be cautious on big data-tables it can lead to slowing down perf
func (r *DB) InRandomOrder() *DB {
	r.OrderByRaw("random()")
	return r
}

// GroupBy adds GROUP BY expression to SQL stmt
func (r *DB) GroupBy(expr string) *DB {
	r.Builder.groupBy = expr
	return r
}

// Having similar to Where but used with GroupBy to apply over the grouped results
func (r *DB) Having(col string, op Op, val interface{}) *DB {
	r.Builder.having.
		Ident(col).
		WriteOp(op).
		Arg(val)
	return r
}

// AddSelect accepts additional columns to select from a table
func (r *DB) AddSelect(args ...string) *DB {
	r.Builder.columns = append(r.Builder.columns, args...)
	return r
}

// SelectRaw accepts custom string to select from a table
func (r *DB) SelectRaw(raw string) *DB {
	r.Builder.columns = []string{raw}
	return r
}

// InnerJoin joins tables by getting elements if found in both
func (r *DB) InnerJoin(table, left, operator, right string) *DB {
	return r.buildJoin(joinInner, table, left+operator+right)
}

// LeftJoin joins tables by getting elements from left without those that null on the right
func (r *DB) LeftJoin(table, left, operator, right string) *DB {
	return r.buildJoin(joinLeft, table, left+operator+right)
}

// RightJoin joins tables by getting elements from right without those that null on the left
func (r *DB) RightJoin(table, left, operator, right string) *DB {
	return r.buildJoin(joinRight, table, left+operator+right)
}

// CrossJoin joins tables by getting intersection of sets
// todo: MySQL/PostgreSQL versions are different here impl their difference
//func (r *DB) CrossJoin(table string, left string, operator string, right string) *DB {
//	return r.buildJoin(JoinCross, table, left+operator+right)
//}

// FullJoin joins tables by getting all elements of both sets
func (r *DB) FullJoin(table, left, operator, right string) *DB {
	return r.buildJoin(joinFull, table, left+operator+right)
}

// FullOuterJoin joins tables by getting an outer sets
func (r *DB) FullOuterJoin(table, left, operator, right string) *DB {
	return r.buildJoin(joinFullOuter, table, left+operator+right)
}

// Union joins multiple queries omitting duplicate records
func (r *DB) Union() *DB {
	r.Builder.union = append(r.Builder.union, r.Builder.buildSelect())
	return r
}

// UnionAll joins multiple queries to select all rows from both tables with duplicate
func (r *DB) UnionAll() *DB {
	r.Union()
	r.Builder.isUnionAll = true
	return r
}

func (r *DB) buildJoin(joinType, table, on string) *DB {
	r.Builder.join = append(r.Builder.join, " "+joinType+" JOIN "+table+" ON "+on+" ")
	return r
}

// Where accepts left operand-operator-right operand to apply them to where clause
func (r *DB) WhereRaw(raw string, val ...interface{}) *DB {
	r.Builder.where.WriteString(" WHERE ").
		WriteString(raw).
		Args(val...)
	return r
}

// Where accepts left operand-operator-right operand to apply them to where clause
func (r *DB) Where(col string, op Op, val interface{}) *DB {
	r.Builder.where.WriteString(" WHERE ").
		Ident(col).
		WriteOp(op).
		Arg(val)
	return r
}

// AndWhere accepts left operand-operator-right operand to apply them to where clause
// with AND logical operator
func (r *DB) AndWhere(col string, op Op, val interface{}) *DB {
	r.Builder.where.WriteString(and).
		Pad().
		Ident(col).
		WriteOp(op).
		Arg(val)
	return r
}

// OrWhere accepts left operand-operator-right operand to apply them to where clause
// with OR logical operator
func (r *DB) OrWhere(col string, op Op, val interface{}) *DB {
	r.Builder.where.WriteString(or).
		Pad().
		Ident(col).
		WriteOp(op).
		Arg(val)
	return r
}

// WhereBetween sets the clause BETWEEN 2 values
func (r *DB) WhereBetween(col string, val1, val2 interface{}) *DB {
	r.Builder.where.WriteString(" WHERE ").
		Ident(col).
		WriteOp(OpBetween).
		Arg(val1).Pad().
		WriteString("AND").Pad().
		Arg(val2)
	return r
}

// OrWhereBetween sets the clause OR BETWEEN 2 values
func (r *DB) OrWhereBetween(col string, val1, val2 interface{}) *DB {
	r.Builder.where.WriteString(or).
		Pad().
		Ident(col).
		WriteOp(OpBetween).
		Arg(val1).Pad().
		WriteString("AND").Pad().
		Arg(val2)
	return r
}

// AndWhereBetween sets the clause AND BETWEEN 2 values
func (r *DB) AndWhereBetween(col string, val1, val2 interface{}) *DB {
	r.Builder.where.WriteString(and).
		Pad().
		Ident(col).
		WriteOp(OpBetween).
		Arg(val1).Pad().
		WriteString("AND").Pad().
		Arg(val2)
	return r
}

// WhereNotBetween sets the clause NOT BETWEEN 2 values
func (r *DB) WhereNotBetween(col string, val1, val2 interface{}) *DB {
	r.Builder.where.WriteString(" WHERE ").
		Ident(col).
		WriteOp(OpNotBetween).
		Arg(val1).Pad().
		WriteString("AND").Pad().
		Arg(val2)
	return r
}

// OrWhereNotBetween sets the clause OR BETWEEN 2 values
func (r *DB) OrWhereNotBetween(col string, val1, val2 interface{}) *DB {
	r.Builder.where.WriteString(or).
		Pad().
		Ident(col).
		WriteOp(OpNotBetween).
		Arg(val1).Pad().
		WriteString("AND").Pad().
		Arg(val2)
	return r
}

// AndWhereNotBetween sets the clause AND BETWEEN 2 values
func (r *DB) AndWhereNotBetween(col string, val1, val2 interface{}) *DB {
	r.Builder.where.WriteString(and).
		Pad().
		Ident(col).
		WriteOp(OpNotBetween).
		Arg(val1).Pad().
		WriteString("AND").Pad().
		Arg(val2)
	return r
}

// Offset accepts offset to start slicing results from
func (r *DB) Offset(off int64) *DB {
	r.Builder.offset = off
	return r
}

// Limit accepts limit to end slicing results to
func (r *DB) Limit(lim int64) *DB {
	r.Builder.limit = lim
	return r
}

// WhereIn appends IN (val1, val2, val3...) stmt to WHERE clause
func (r *DB) WhereIn(col string, in ...interface{}) *DB {
	r.Builder.where.WriteString(" WHERE ").
		Ident(col).
		WriteOp(OpIn).
		Nested(func(b *sqlBuilder) {
			b.Args(in...)
		})
	return r
}

// WhereNotIn appends NOT IN (val1, val2, val3...) stmt to WHERE clause
func (r *DB) WhereNotIn(col string, in ...interface{}) *DB {
	r.Builder.where.WriteString(" WHERE ").
		Ident(col).
		WriteOp(OpNotIn).
		Nested(func(b *sqlBuilder) {
			b.Args(in...)
		})
	return r
}

// OrWhereIn appends OR IN (val1, val2, val3...) stmt to WHERE clause
func (r *DB) OrWhereIn(col string, in ...interface{}) *DB {
	r.Builder.where.WriteString(or).
		Pad().
		Ident(col).
		WriteOp(OpIn).
		Nested(func(b *sqlBuilder) {
			b.Args(in...)
		})
	return r
}

// OrWhereNotIn appends OR NOT IN (val1, val2, val3...) stmt to WHERE clause
func (r *DB) OrWhereNotIn(col string, in ...interface{}) *DB {
	r.Builder.where.WriteString(or).
		Pad().
		Ident(col).
		WriteOp(OpNotIn).
		Nested(func(b *sqlBuilder) {
			b.Args(in...)
		})
	return r
}

// AndWhereIn appends OR IN (val1, val2, val3...) stmt to WHERE clause
func (r *DB) AndWhereIn(col string, in ...interface{}) *DB {
	r.Builder.where.WriteString(and).
		Pad().
		Ident(col).
		WriteOp(OpIn).
		Nested(func(b *sqlBuilder) {
			b.Args(in...)
		})
	return r
}

// AndWhereNotIn appends OR NOT IN (val1, val2, val3...) stmt to WHERE clause
func (r *DB) AndWhereNotIn(col string, in ...interface{}) *DB {
	r.Builder.where.WriteString(and).
		Pad().
		Ident(col).
		WriteOp(OpNotIn).
		Nested(func(b *sqlBuilder) {
			b.Args(in...)
		})
	return r
}

// WhereNull appends col IS NULL stmt to WHERE clause
func (r *DB) WhereNull(col string) *DB {
	r.Builder.where.WriteString(" WHERE ").
		Ident(col).
		WriteOp(OpIsNull)
	return r
}

// WhereNotNull appends col IS NOT NULL stmt to WHERE clause
func (r *DB) WhereNotNull(col string) *DB {
	r.Builder.where.WriteString(" WHERE ").
		Ident(col).
		WriteOp(OpNotNull)
	return r
}

// OrWhereNull appends col IS NULL stmt to WHERE clause
func (r *DB) OrWhereNull(col string) *DB {
	r.Builder.where.WriteString(or).
		Pad().
		Ident(col).
		WriteOp(OpIsNull)
	return r
}

// OrWhereNotNull appends col IS NOT NULL stmt to WHERE clause
func (r *DB) OrWhereNotNull(col string) *DB {
	r.Builder.where.WriteString(or).
		Pad().
		Ident(col).
		WriteOp(OpNotNull)
	return r
}

// AndWhereNull appends col IS NULL stmt to WHERE clause
func (r *DB) AndWhereNull(col string) *DB {
	r.Builder.where.WriteString(and).
		Pad().
		Ident(col).
		WriteOp(OpIsNull)
	return r
}

// AndWhereNotNull appends col IS NOT NULL stmt to WHERE clause
func (r *DB) AndWhereNotNull(col string) *DB {
	r.Builder.where.WriteString(and).
		Pad().
		Ident(col).
		WriteOp(OpNotNull)
	return r
}

// WhereLike appends col is LIKE pattern stmt to WHERE clause
func (r *DB) WhereLike(col string, pattern string) *DB {
	r.Builder.where.WriteString(" WHERE ").
		Ident(col).
		WriteOp(OpLike).
		Args(pattern)
	return r
}

// OrWhereLike appends col is LIKE pattern stmt to WHERE clause
func (r *DB) OrWhereLike(col string, pattern string) *DB {
	r.Builder.where.WriteString(or).
		Pad().
		Ident(col).
		WriteOp(OpLike).
		Args(pattern)
	return r
}

// AndWhereLike appends col is LIKE pattern stmt to WHERE clause
func (r *DB) AndWhereLike(col string, pattern string) *DB {
	r.Builder.where.WriteString(and).
		Pad().
		Ident(col).
		WriteOp(OpLike).
		Args(pattern)
	return r
}

// WhereNotLike appends col is NOT LIKE pattern stmt to WHERE clause
func (r *DB) WhereNotLike(col string, pattern string) *DB {
	r.Builder.where.WriteString(" WHERE ").
		Ident(col).
		WriteOp(OpNotLike).
		Args(pattern)
	return r
}

// OrWhereNotLike appends col is NOT LIKE pattern stmt to WHERE clause
func (r *DB) OrWhereNotLike(col string, pattern string) *DB {
	r.Builder.where.WriteString(or).
		Pad().
		Ident(col).
		WriteOp(OpNotLike).
		Args(pattern)
	return r
}

// AndWhereNotLike appends col is NOT LIKE pattern stmt to WHERE clause
func (r *DB) AndWhereNotLike(col string, pattern string) *DB {
	r.Builder.where.WriteString(and).
		Pad().
		Ident(col).
		WriteOp(OpNotLike).
		Args(pattern)
	return r
}

// WhereEmpty appends col IS NOT NULL and IS EMPRTY str stmt to WHERE clause
func (r *DB) WhereEmpty(col string) *DB {
	r.Builder.where.WriteString(" WHERE ").
		Ident(col).
		WriteOp(OpEQ).
		Arg("").
		WriteString(and).
		Ident(col).
		WriteOp(OpIsNull)
	return r
}

// OrWhereEmpty appends col IS NOT NULL and IS EMPRTY str stmt to WHERE clause
func (r *DB) OrWhereEmpty(col string) *DB {
	r.Builder.where.WriteString(or).
		Pad().
		Ident(col).
		WriteOp(OpEQ).
		Arg("").
		WriteString(and).
		Ident(col).
		WriteOp(OpIsNull)
	return r
}

// AndWhereEmpty appends col IS NOT NULL and IS EMPRTY str stmt to WHERE clause
func (r *DB) AndWhereEmpty(col string) *DB {
	r.Builder.where.WriteString(and).
		Pad().
		Ident(col).
		WriteOp(OpEQ).
		Arg("").
		WriteString(and).
		Ident(col).
		WriteOp(OpIsNull)
	return r
}

// From prepares sql stmt to set data from another table, ex.:
// UPDATE employees SET sales_count = sales_count + 1 FROM accounts
func (r *DB) From(fromTbl string) *DB {
	r.Builder.from = fromTbl
	return r
}

// LockForUpdate locks table/row
func (r *DB) LockForUpdate() *DB {
	str := " FOR UPDATE"
	r.Builder.lockForUpdate = &str
	return r
}

// Dump prints raw sql to stdout
func (r *DB) Dump() {
	log.SetOutput(os.Stdout)
	log.Println(r.Builder.buildSelect())
}

// Dd prints raw sql to stdout and exit
func (r *DB) Dd() {
	r.Dump()
	os.Exit(0)
}
