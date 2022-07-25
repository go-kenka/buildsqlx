package buildsqlx

// Count counts
func (r *DB) Count() (query string, args []interface{}) {
	builder := r.Builder
	builder.columns = []string{"COUNT(*)"}
	query = builder.buildSelect()
	args = append(args, builder.where.args...)
	args = append(args, builder.having.args...)
	return
}

// Avg calculates average for specified column
func (r *DB) Avg(column string) (query string, args []interface{}) {
	builder := r.Builder
	builder.columns = []string{"AVG(" + column + ")"}
	query = builder.buildSelect()
	args = append(args, builder.where.args...)
	args = append(args, builder.having.args...)
	return
}

// Min calculates minimum for specified column
func (r *DB) Min(column string) (query string, args []interface{}) {
	builder := r.Builder
	builder.columns = []string{"MIN(" + column + ")"}
	query = builder.buildSelect()
	args = append(args, builder.where.args...)
	args = append(args, builder.having.args...)
	return
}

// Max calculates maximum for specified column
func (r *DB) Max(column string) (query string, args []interface{}) {
	builder := r.Builder
	builder.columns = []string{"MAX(" + column + ")"}
	query = builder.buildSelect()
	args = append(args, builder.where.args...)
	args = append(args, builder.having.args...)
	return
}

// Sum calculates sum for specified column
func (r *DB) Sum(column string) (query string, args []interface{}) {
	builder := r.Builder
	builder.columns = []string{"SUM(" + column + ")"}
	query = builder.buildSelect()
	args = append(args, builder.where.args...)
	args = append(args, builder.having.args...)
	return
}
