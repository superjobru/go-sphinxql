// Copyright 2018 Huan Du. All rights reserved.
// Copyright 2022 OOO SuperJob. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package sphinxql

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	updateMarkerInit injectionMarker = iota
	updateMarkerAfterUpdate
	updateMarkerAfterSet
	updateMarkerAfterWhere
	updateMarkerAfterOption
)

// NewUpdateBuilder creates a new UPDATE builder.
func NewUpdateBuilder() *UpdateBuilder {
	return DefaultFlavor.NewUpdateBuilder()
}

func newUpdateBuilder() *UpdateBuilder {
	args := &Args{}
	return &UpdateBuilder{
		Cond: Cond{
			Args: args,
		},
		args:      args,
		injection: newInjection(),
	}
}

// UpdateBuilder is a builder to build UPDATE.
type UpdateBuilder struct {
	Cond
	OptionBuilder

	table       string
	assignments []string
	whereExprs  []string
	orderByCols []string
	options     []Option

	args *Args

	injection *injection
	marker    injectionMarker
}

var _ Builder = new(UpdateBuilder)

// Update sets table name in UPDATE.
func Update(table string) *UpdateBuilder {
	return DefaultFlavor.NewUpdateBuilder().Update(table)
}

// Update sets table name in UPDATE.
func (ub *UpdateBuilder) Update(table string) *UpdateBuilder {
	ub.table = Escape(table)
	ub.marker = updateMarkerAfterUpdate
	return ub
}

// Set sets the assignments in SET.
func (ub *UpdateBuilder) Set(assignment ...string) *UpdateBuilder {
	ub.assignments = assignment
	ub.marker = updateMarkerAfterSet
	return ub
}

// SetMore appends the assignments in SET.
func (ub *UpdateBuilder) SetMore(assignment ...string) *UpdateBuilder {
	ub.assignments = append(ub.assignments, assignment...)
	ub.marker = updateMarkerAfterSet
	return ub
}

// Where sets expressions of WHERE in UPDATE.
func (ub *UpdateBuilder) Where(andExpr ...string) *UpdateBuilder {
	ub.whereExprs = append(ub.whereExprs, andExpr...)
	ub.marker = updateMarkerAfterWhere
	return ub
}

// Option sets the OPTION in UPDATE.
func (ub *UpdateBuilder) Option(option ...Option) *UpdateBuilder {
	ub.options = option
	ub.marker = updateMarkerAfterOption
	return ub
}

// Assign represents SET "field = value" in UPDATE.
func (ub *UpdateBuilder) Assign(field string, value interface{}) string {
	return fmt.Sprintf("%s = %s", Escape(field), ub.args.Add(value))
}

// Incr represents SET "field = field + 1" in UPDATE.
func (ub *UpdateBuilder) Incr(field string) string {
	f := Escape(field)
	return fmt.Sprintf("%s = %s + 1", f, f)
}

// Decr represents SET "field = field - 1" in UPDATE.
func (ub *UpdateBuilder) Decr(field string) string {
	f := Escape(field)
	return fmt.Sprintf("%s = %s - 1", f, f)
}

// Add represents SET "field = field + value" in UPDATE.
func (ub *UpdateBuilder) Add(field string, value interface{}) string {
	f := Escape(field)
	return fmt.Sprintf("%s = %s + %s", f, f, ub.args.Add(value))
}

// Sub represents SET "field = field - value" in UPDATE.
func (ub *UpdateBuilder) Sub(field string, value interface{}) string {
	f := Escape(field)
	return fmt.Sprintf("%s = %s - %s", f, f, ub.args.Add(value))
}

// Mul represents SET "field = field * value" in UPDATE.
func (ub *UpdateBuilder) Mul(field string, value interface{}) string {
	f := Escape(field)
	return fmt.Sprintf("%s = %s * %s", f, f, ub.args.Add(value))
}

// Div represents SET "field = field / value" in UPDATE.
func (ub *UpdateBuilder) Div(field string, value interface{}) string {
	f := Escape(field)
	return fmt.Sprintf("%s = %s / %s", f, f, ub.args.Add(value))
}

// String returns the compiled UPDATE string.
func (ub *UpdateBuilder) String() string {
	s, _ := ub.Build()
	return s
}

// Build returns compiled UPDATE string and args.
// They can be used in `DB#Query` of package `database/sql` directly.
func (ub *UpdateBuilder) Build() (sql string, args []interface{}) {
	return ub.BuildWithFlavor(ub.args.Flavor)
}

// BuildWithFlavor returns compiled UPDATE string and args with flavor and initial args.
// They can be used in `DB#Query` of package `database/sql` directly.
func (ub *UpdateBuilder) BuildWithFlavor(flavor Flavor, initialArg ...interface{}) (sql string, args []interface{}) {
	buf := &bytes.Buffer{}
	ub.injection.WriteTo(buf, updateMarkerInit)
	buf.WriteString("UPDATE ")
	buf.WriteString(ub.table)
	ub.injection.WriteTo(buf, updateMarkerAfterUpdate)

	buf.WriteString(" SET ")
	buf.WriteString(strings.Join(ub.assignments, ", "))
	ub.injection.WriteTo(buf, updateMarkerAfterSet)

	if len(ub.whereExprs) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(ub.whereExprs, " AND "))
		ub.injection.WriteTo(buf, updateMarkerAfterWhere)
	}

	if len(ub.options) > 0 {
		buf.WriteString(" OPTION ")
		buf.WriteString(SerializeOptions(ub.options))

		ub.injection.WriteTo(buf, updateMarkerAfterOption)
	}

	return ub.args.CompileWithFlavor(buf.String(), flavor, initialArg...)
}

// SetFlavor sets the flavor of compiled sql.
func (ub *UpdateBuilder) SetFlavor(flavor Flavor) (old Flavor) {
	old = ub.args.Flavor
	ub.args.Flavor = flavor
	return
}

// SQL adds an arbitrary sql to current position.
func (ub *UpdateBuilder) SQL(sql string) *UpdateBuilder {
	ub.injection.SQL(ub.marker, sql)
	return ub
}
