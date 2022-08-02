package buildsqlx

import (
	"strconv"
	"strings"
)

const (
	// Errors
	errTableCallBeforeOp = "sql: there was no Table() call with table name set"
)

// buildSelect constructs a query for select statement
func (r *builder) buildSelect() string {

	// SELECT
	r.WriteString("SELECT").Pad()

	// field
	l := len(r.columns)
	for k, col := range r.columns {
		if col == "*" || strings.HasSuffix(col, ".*") {
			r.WriteString(col)
		} else {
			if strings.Contains(col, "`") {
				r.WriteString(col)
			} else {
				r.Ident(col)
			}
		}
		if k < l-1 {
			r.Comma()
		}
	}

	// from
	r.Pad().WriteString("FROM").Pad().Ident(r.table).Pad()

	// Clauses
	r.buildClauses()

	return r.String()
}

// builds query string clauses
func (r *builder) buildClauses() {
	for _, j := range r.join {
		r.WriteString(j)
	}

	// build where clause
	if r.where.Len() > 0 {
		r.WriteString(r.where.String())
	}

	if r.groupBy != "" {
		r.Pad().WriteString("GROUP BY").Pad()
		r.Ident(r.groupBy)
	}

	if r.having.Len() > 0 {
		r.Pad().WriteString("HAVING").Pad()
		r.WriteString(r.having.String())
	}

	r.composeOrderBy()

	if r.limit > 0 {
		r.Pad().WriteString("LIMIT").Pad()
		r.WriteString(strconv.FormatInt(r.limit, 10))
	}

	if r.offset > 0 {
		r.Comma()
		r.WriteString(strconv.FormatInt(r.offset, 10))
	}

	if r.lockForUpdate != nil {
		r.WriteString(*r.lockForUpdate)
	}
}

// builds query string clauses
func buildClauses(r *builder) string {

	b := deepClone(r)
	b.sb.Reset()

	for _, j := range b.join {
		b.WriteString(j)
	}

	// build where clause
	if b.where.Len() > 0 {
		b.WriteString(b.where.String())
	}

	if b.groupBy != "" {
		b.Pad().WriteString("GROUP BY").Pad()
		b.Ident(b.groupBy)
	}

	if b.having.Len() > 0 {
		b.Pad().WriteString("HAVING").Pad()
		b.WriteString(b.having.String())
	}

	b.composeOrderBy()

	if b.limit > 0 {
		b.Pad().WriteString("LIMIT").Pad()
		b.WriteString(strconv.FormatInt(b.limit, 10))
	}

	if b.offset > 0 {
		b.Comma()
		b.WriteString(strconv.FormatInt(b.offset, 10))
	}

	if b.lockForUpdate != nil {
		b.WriteString(*b.lockForUpdate)
	}

	return b.String()
}

// composers ORDER BY clause string for particular query stmt
func (r *builder) composeOrderBy() {
	if len(r.orderBy) > 0 {
		fist := true
		for f, d := range r.orderBy {
			if fist {
				fist = false
				r.Pad().WriteString("ORDER BY").Pad().Ident(f).Pad().WriteString(d)
			} else {
				r.Pad().Comma().Ident(f).Pad().WriteString(d)
			}
		}
		return
	} else if r.orderByRaw != nil {
		r.Pad().WriteString("ORDER BY").Pad().WriteString(*r.orderByRaw)
	}
}

// Insert inserts one row with param bindings
func (r *DB) Insert(data map[string]interface{}) (query string, values []interface{}) {
	builder := r.Builder
	if builder.table == "" {
		panic(errTableCallBeforeOp)
	}

	columns, values, bindings := prepareBindings(data)

	builder.WriteString("INSERT INTO").
		Pad().Ident(builder.table).Pad().
		Nested(func(s *sqlBuilder) {
			l := len(columns)
			for k, col := range columns {
				s.Ident(col)
				if k < l-1 {
					s.Comma()
				}
			}
		}).
		Pad().WriteString("VALUES").
		Nested(func(s *sqlBuilder) {
			s.WriteString(strings.Join(bindings, `, `))
		})

	query = builder.String()

	return
}

// prepareBindings prepares slices to split in favor of INSERT sql statement
func prepareBindings(data map[string]interface{}) (columns []string, values []interface{}, bindings []string) {
	i := 1
	for column, value := range data {
		columns = append(columns, column)
		values = append(values, value)
		bindings = append(bindings, "?")
		i++
	}

	return
}

// InsertBatch inserts multiple rows based on transaction
func (r *DB) InsertBatch(data []map[string]interface{}) (query string, values [][]interface{}) {
	builder := r.Builder
	if builder.table == "" {
		panic(errTableCallBeforeOp)
	}

	columns, values, bindings := prepareInsertBatch(data)

	builder.WriteString("INSERT INTO").
		Pad().Ident(builder.table).Pad().
		Nested(func(s *sqlBuilder) {
			l := len(columns)
			for k, col := range columns {
				s.Ident(col)
				if k < l-1 {
					s.Comma()
				}
			}
		}).
		Pad().WriteString("VALUES").
		Nested(func(s *sqlBuilder) {
			s.WriteString(strings.Join(bindings, `, `))
		})

	query = builder.String()

	return
}

// prepareInsertBatch prepares slices to split in favor of INSERT sql statement
func prepareInsertBatch(data []map[string]interface{}) (columns []string, values [][]interface{}, bindings []string) {
	values = make([][]interface{}, len(data))
	colToIdx := make(map[string]int)

	i := 0
	for k, v := range data {
		values[k] = make([]interface{}, len(v))

		for column, value := range v {
			if k == 0 {
				columns = append(columns, column)
				bindings = append(bindings, "?")
				// todo: don't know yet how to match them explicitly (it is bad idea, but it works well now)
				colToIdx[column] = i
				i++
			}

			values[k][colToIdx[column]] = value
		}
	}

	return
}

// Update builds an UPDATE sql stmt with corresponding where/from clauses if stated
// returning affected rows
func (r *DB) Update(data map[string]interface{}) (query string, values []interface{}) {
	builder := r.Builder
	if builder.table == "" {
		panic(errTableCallBeforeOp)
	}

	columns, values, bindings := prepareBindings(data)

	builder.WriteString("UPDATE").
		Pad().Ident(builder.table).Pad().
		WriteString("SET")

	l := len(columns)
	for k, col := range columns {
		builder.Pad().Ident(col).WriteOp(OpEQ).WriteString(bindings[k])
		if k < l-1 {
			builder.Comma()
		}
	}

	builder.buildClauses()

	query += builder.String()
	values = append(values, r.Builder.where.args...)

	return
}

// Delete builds a DELETE stmt with corresponding where clause if stated
// returning affected rows
func (r *DB) Delete() (query string, values []interface{}) {
	builder := r.Builder
	if builder.table == "" {
		panic(errTableCallBeforeOp)
	}

	builder.WriteString("DELETE FROM").
		Pad().Ident(builder.table).Pad()

	builder.buildClauses()

	query = builder.String()
	values = r.Builder.where.args

	return
}

// Replace inserts data if conflicting row hasn't been found, else it will update an existing one
func (r *DB) Replace(data map[string]interface{}, conflict string) (query string, values []interface{}) {
	builder := r.Builder
	if builder.table == "" {
		panic(errTableCallBeforeOp)
	}

	columns, values, bindings := prepareBindings(data)

	builder.WriteString("INSERT INTO").
		Pad().Ident(builder.table).Pad().
		Nested(func(s *sqlBuilder) {
			s.WriteString(strings.Join(columns, `, `))
		}).
		Pad().WriteString("VALUES").
		Nested(func(s *sqlBuilder) {
			s.WriteString(strings.Join(bindings, `, `))
		}).
		WriteString("ON DUPLICATE KEY UPDATE").Pad()

	l := len(columns)
	for i, v := range columns {
		builder.Ident(v).WriteOp(OpEQ).Pad().WriteString("excluded.").WriteString(v)
		if i < l-1 {
			builder.Comma()
		}
	}

	query = builder.String()

	return
}

// Drop drops >=1 tables
func (r *DB) Drop(tables string) string {
	return "DROP TABLE " + r.Builder.Quote(tables)
}

// Truncate clears >=1 tables
func (r *DB) Truncate(tables string) string {
	return "TRUNCATE " + r.Builder.Quote(tables)
}

// DropIfExists drops >=1 tables if they are existent
func (r *DB) DropIfExists(tables string) string {
	return "DROP TABLE IF EXISTS " + r.Builder.Quote(tables)
}

// Rename renames from - to new table name
func (r *DB) Rename(from, to string) string {
	return "ALTER TABLE " + r.Builder.Quote(from) + " RENAME TO " + r.Builder.Quote(to)
}
