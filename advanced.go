package buildsqlx

// Exists checks whether conditional rows are existing (returns true) or not (returns false)
func (r *DB) Exists() (query string, values []interface{}) {
	builder := r.Builder
	if builder.table == "" {
		panic(errTableCallBeforeOp)
	}

	builder.WriteString("SELECT EXISTS").
		Pad().
		Nested(func(s *sqlBuilder) {
			s.WriteString("SELECT 1 FROM")
			s.Pad()
			s.Ident(builder.table)
			s.Pad().WriteString(buildClauses(builder))
		}).Pad()

	query = builder.String()
	values = append(values, builder.where.args...)
	values = append(values, builder.having.args...)
	return
}

// Query builds all sql statements and return sql & values
func (r *DB) Query() (query string, values []interface{}) {
	builder := r.Builder
	if builder.table == "" {
		panic(errTableCallBeforeOp)
	}

	if len(builder.union) > 0 { // got union - need different logic to glue
		for _, uBuilder := range builder.union {
			builder.WriteString(uBuilder)
			builder.Pad().WriteString("UNION").Pad()

			if builder.isUnionAll {
				builder.WriteString("ALL").Pad()
			}
		}

		query += builder.buildSelect()
	} else { // std builder
		query = builder.buildSelect()
	}

	values = append(values, r.Builder.where.args...)
	values = append(values, r.Builder.having.args...)

	return

}
