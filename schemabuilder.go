package buildsqlx

import (
	"strings"
)

type schemaBuilder struct {
	sb    *strings.Builder
	child *schemaBuilder
}

func newSchemaBuilder() *schemaBuilder {
	return &schemaBuilder{sb: &strings.Builder{}, child: &schemaBuilder{sb: &strings.Builder{}}}
}

func (b *schemaBuilder) WriteString(s string) *schemaBuilder {
	if b.sb == nil {
		b.sb = &strings.Builder{}
	}

	b.sb.WriteString(s)
	return b
}
func (b *schemaBuilder) WriteByte(s byte) *schemaBuilder {
	if b.sb == nil {
		b.sb = &strings.Builder{}
	}

	b.sb.WriteByte(s)
	return b
}
func (b *schemaBuilder) String() string {
	if b.sb == nil {
		return ""
	}

	return b.sb.String()
}

// Ident adds a comma to the query.
func (b *schemaBuilder) Ident(str string) *schemaBuilder {
	b.WriteString(b.Quote(str))
	return b
}

// Ident adds a point to the query.
func (b *schemaBuilder) IdentPoint(str string) *schemaBuilder {
	b.WriteString(b.Quote(str)).WriteByte('.')
	return b
}

// Ident adds a comma to the query.
func (b *schemaBuilder) Quote(ident string) string {
	quote := "`"
	return quote + ident + quote
}

// Comma adds a comma to the query.
func (b *schemaBuilder) Comma() *schemaBuilder {
	return b.WriteString(", ")
}

// Pad adds a space to the query.
func (b *schemaBuilder) Pad() *schemaBuilder {
	return b.WriteByte(' ')
}

// SemiColon adds a ; to the query.
func (b *schemaBuilder) SemiColon() *schemaBuilder {
	return b.WriteByte(';')
}

// Nested gets a callback, and wraps its result with parentheses.
func (b *schemaBuilder) Nested(f func(*schemaBuilder)) *schemaBuilder {
	nb := newSchemaBuilder()
	nb.WriteByte('(')
	f(nb)
	nb.WriteByte(')')
	b.WriteString(nb.String())
	return b
}
