package buildsqlx

import (
	"fmt"
	"strings"
)

// An Op represents an operator.
type Op int

const (
	// Predicate operators.
	OpEQ         Op = iota // =
	OpNEQ                  // <>
	OpGT                   // >
	OpGTE                  // >=
	OpLT                   // <
	OpLTE                  // <=
	OpIn                   // IN
	OpNotIn                // NOT IN
	OpLike                 // LIKE
	OpNotLike              // NOT LIKE
	OpBetween              // BETWEEN
	OpNotBetween           // NOT BETWEEN
	OpIsNull               // IS NULL
	OpNotNull              // IS NOT NULL
)

var ops = [...]string{
	OpEQ:         "=",
	OpNEQ:        "<>",
	OpGT:         ">",
	OpGTE:        ">=",
	OpLT:         "<",
	OpLTE:        "<=",
	OpIn:         "IN",
	OpNotIn:      "NOT IN",
	OpLike:       "LIKE",
	OpNotLike:    "NOT LIKE",
	OpIsNull:     "IS NULL",
	OpNotNull:    "IS NOT NULL",
	OpBetween:    "BETWEEN",
	OpNotBetween: "NOT BETWEEN",
}

type sqlBuilder struct {
	sb   *strings.Builder
	args []interface{}
}

// Query returns query representation of a predicate.
func (b *sqlBuilder) Query() (string, []interface{}) {
	return b.String(), b.args
}

func (b *sqlBuilder) WriteString(s string) *sqlBuilder {
	if b.sb == nil {
		b.sb = &strings.Builder{}
	}

	b.sb.WriteString(s)
	return b
}
func (b *sqlBuilder) WriteByte(s byte) *sqlBuilder {
	if b.sb == nil {
		b.sb = &strings.Builder{}
	}

	b.sb.WriteByte(s)
	return b
}
func (b *sqlBuilder) String() string {
	if b.sb == nil {
		return ""
	}

	return b.sb.String()
}
func (b *sqlBuilder) Len() int {
	if b.sb == nil {
		return 0
	}

	return b.sb.Len()
}

// WriteOp writes an operator to the builder.
func (b *sqlBuilder) WriteOp(op Op) *sqlBuilder {
	switch {
	case op >= OpEQ && op <= OpNotBetween:
		b.Pad().WriteString(ops[op]).Pad()
	case op == OpIsNull || op == OpNotNull:
		b.Pad().WriteString(ops[op])
	default:
		panic(fmt.Sprintf("invalid op %d", op))
	}
	return b
}

// Arg appends an input argument to the builder.
func (b *sqlBuilder) Arg(a interface{}) *sqlBuilder {
	b.args = append(b.args, a)
	// Default placeholder param (MySQL and SQLite).
	param := "?"
	b.WriteString(param)
	return b
}

// Args appends a list of arguments to the builder.
func (b *sqlBuilder) Args(a ...interface{}) *sqlBuilder {
	for i := range a {
		if i > 0 {
			b.Comma()
		}
		b.Arg(a[i])
	}
	return b
}

// Args appends a list of arguments to the builder.
func (b *sqlBuilder) Params(a ...interface{}) *sqlBuilder {
	for i := range a {
		b.args = append(b.args, a[i])
	}
	return b
}

// Comma adds a comma to the query.
func (b *sqlBuilder) Comma() *sqlBuilder {
	return b.WriteString(", ")
}

// Pad adds a space to the query.
func (b *sqlBuilder) Pad() *sqlBuilder {
	return b.WriteByte(' ')
}

// Ident adds a comma to the query.
func (b *sqlBuilder) Ident(str string) *sqlBuilder {
	b.WriteString(b.Quote(str))
	return b
}

// Ident adds a point to the query.
func (b *sqlBuilder) IdentPoint(str string) *sqlBuilder {
	b.WriteString(b.Quote(str)).WriteByte('.')
	return b
}

// Ident adds a comma to the query.
func (b *sqlBuilder) Quote(ident string) string {
	quote := "`"
	return quote + ident + quote
}

// Nested gets a callback, and wraps its result with parentheses.
func (b *sqlBuilder) Nested(f func(*sqlBuilder)) *sqlBuilder {
	nb := &sqlBuilder{sb: &strings.Builder{}}
	nb.WriteByte('(')
	f(nb)
	nb.WriteByte(')')
	b.WriteString(nb.String())
	b.args = append(b.args, nb.args...)
	return b
}
