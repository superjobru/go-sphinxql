// Copyright 2018 Huan Du. All rights reserved.
// Copyright 2022 OOO SuperJob. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package sphinxql

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	selectMarkerInit injectionMarker = iota
	selectMarkerAfterSelect
	selectMarkerAfterFrom
	selectMarkerAfterWhere
	selectMarkerAfterGroupBy
	selectMarkerAfterWithinGroupOrderBy
	selectMarkerAfterOrderBy
	selectMarkerAfterLimit
	selectMarkerAfterOption
)

// NewSelectBuilder creates a new SELECT builder.
func NewSelectBuilder() *SelectBuilder {
	return DefaultFlavor.NewSelectBuilder()
}

func newSelectBuilder() *SelectBuilder {
	args := &Args{}
	return &SelectBuilder{
		Cond: Cond{
			Args: args,
		},
		Opt: Opt{
			Args: args,
		},
		limit:     -1,
		offset:    -1,
		args:      args,
		injection: newInjection(),
	}
}

// SelectBuilder is a builder to build SELECT.
type SelectBuilder struct {
	Cond
	OrdBy
	Opt

	tables                  []string
	selectCols              []string
	whereExprs              []string
	groupByCols             []string
	withinGroupOrderByExprs []string
	havingExprs             []string
	orderByExprs            []string
	limit                   int
	offset                  int
	optionExprs             []string

	args *Args

	injection *injection
	marker    injectionMarker
}

var _ Builder = new(SelectBuilder)

// Select sets columns in SELECT.
func Select(col ...string) *SelectBuilder {
	return DefaultFlavor.NewSelectBuilder().Select(col...)
}

// Select sets columns in SELECT.
func (sb *SelectBuilder) Select(col ...string) *SelectBuilder {
	sb.selectCols = col
	sb.marker = selectMarkerAfterSelect
	return sb
}

// From sets table names in SELECT.
func (sb *SelectBuilder) From(table ...string) *SelectBuilder {
	sb.tables = table
	sb.marker = selectMarkerAfterFrom
	return sb
}

// Where sets expressions of WHERE in SELECT.
func (sb *SelectBuilder) Where(andExpr ...string) *SelectBuilder {
	sb.whereExprs = append(sb.whereExprs, andExpr...)
	sb.marker = selectMarkerAfterWhere
	return sb
}

// Having sets expressions of HAVING in SELECT.
func (sb *SelectBuilder) Having(andExpr ...string) *SelectBuilder {
	sb.havingExprs = append(sb.havingExprs, andExpr...)
	sb.marker = selectMarkerAfterGroupBy
	return sb
}

// GroupBy sets columns of GROUP BY in SELECT.
func (sb *SelectBuilder) GroupBy(col ...string) *SelectBuilder {
	sb.groupByCols = col
	sb.marker = selectMarkerAfterGroupBy
	return sb
}

func filterEmpty(values *[]string) (r []string) {
	for _, v := range *values {
		if v != "" {
			r = append(r, v)
		}
	}
	return
}

// WithinGroupOrderBy sets expressions of WITHIN GROUP ORDER BY in SELECT.
func (sb *SelectBuilder) WithinGroupOrderBy(withinGroupOrderByExpr ...string) *SelectBuilder {
	sb.withinGroupOrderByExprs = filterEmpty(&withinGroupOrderByExpr)
	sb.marker = selectMarkerAfterWithinGroupOrderBy
	return sb
}

// OrderBy sets expressions of ORDER BY in SELECT.
func (sb *SelectBuilder) OrderBy(orderByExpr ...string) *SelectBuilder {
	sb.orderByExprs = filterEmpty(&orderByExpr)
	sb.marker = selectMarkerAfterOrderBy
	return sb
}

// Limit sets the LIMIT in SELECT.
func (sb *SelectBuilder) Limit(limit int) *SelectBuilder {
	sb.limit = limit
	sb.marker = selectMarkerAfterLimit
	return sb
}

// Offset sets the LIMIT offset in SELECT.
func (sb *SelectBuilder) Offset(offset int) *SelectBuilder {
	sb.offset = offset
	sb.marker = selectMarkerAfterLimit
	return sb
}

// Option sets expressions of OPTION in SELECT.
func (sb *SelectBuilder) Option(optionExpr ...string) *SelectBuilder {
	sb.optionExprs = optionExpr
	sb.marker = selectMarkerAfterOption
	return sb
}

// As returns an AS expression.
func (sb *SelectBuilder) As(name, alias string) string {
	return fmt.Sprintf("%s AS %s", name, alias)
}

// BuilderAs returns an AS expression wrapping a complex SQL.
// According to SQL syntax, SQL built by builder is surrounded by parens.
func (sb *SelectBuilder) BuilderAs(builder Builder, alias string) string {
	return fmt.Sprintf("(%s) AS %s", sb.Var(builder), alias)
}

// String returns the compiled SELECT string.
func (sb *SelectBuilder) String() string {
	s, _ := sb.Build()
	return s
}

// Build returns compiled SELECT string and args.
// They can be used in `DB#Query` of package `database/sql` directly.
func (sb *SelectBuilder) Build() (sql string, args []interface{}) {
	return sb.BuildWithFlavor(sb.args.Flavor)
}

// BuildWithFlavor returns compiled SELECT string and args with flavor and initial args.
// They can be used in `DB#Query` of package `database/sql` directly.
func (sb *SelectBuilder) BuildWithFlavor(flavor Flavor, initialArg ...interface{}) (sql string, args []interface{}) {
	buf := &strings.Builder{}
	sb.injection.WriteTo(buf, selectMarkerInit)
	buf.WriteString("SELECT ")

	buf.WriteString(strings.Join(sb.selectCols, ", "))
	sb.injection.WriteTo(buf, selectMarkerAfterSelect)

	buf.WriteString(" FROM ")
	buf.WriteString(strings.Join(sb.tables, ", "))
	sb.injection.WriteTo(buf, selectMarkerAfterFrom)

	if len(sb.whereExprs) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(sb.whereExprs, " AND "))

		sb.injection.WriteTo(buf, selectMarkerAfterWhere)
	}

	if len(sb.groupByCols) > 0 {
		buf.WriteString(" GROUP BY ")
		buf.WriteString(strings.Join(sb.groupByCols, ", "))

		if len(sb.havingExprs) > 0 {
			buf.WriteString(" HAVING ")
			buf.WriteString(strings.Join(sb.havingExprs, " AND "))
		}

		sb.injection.WriteTo(buf, selectMarkerAfterGroupBy)
	}

	if len(sb.withinGroupOrderByExprs) > 0 {
		buf.WriteString(" WITHIN GROUP ORDER BY ")
		buf.WriteString(strings.Join(sb.withinGroupOrderByExprs, ", "))

		sb.injection.WriteTo(buf, selectMarkerAfterWithinGroupOrderBy)
	}

	if len(sb.orderByExprs) > 0 {
		buf.WriteString(" ORDER BY ")
		buf.WriteString(strings.Join(sb.orderByExprs, ", "))

		sb.injection.WriteTo(buf, selectMarkerAfterOrderBy)
	}

	if sb.limit >= 0 {
		buf.WriteString(" LIMIT ")
		buf.WriteString(strconv.Itoa(sb.limit))

		if sb.offset >= 0 {
			buf.WriteString(" OFFSET ")
			buf.WriteString(strconv.Itoa(sb.offset))
		}

		sb.injection.WriteTo(buf, selectMarkerAfterLimit)
	}

	if len(sb.optionExprs) > 0 {
		buf.WriteString(" OPTION ")
		buf.WriteString(strings.Join(sb.optionExprs, ", "))

		sb.injection.WriteTo(buf, selectMarkerAfterOption)
	}

	return sb.args.CompileWithFlavor(buf.String(), flavor, initialArg...)
}

// SetFlavor sets the flavor of compiled sql.
func (sb *SelectBuilder) SetFlavor(flavor Flavor) (old Flavor) {
	old = sb.args.Flavor
	sb.args.Flavor = flavor
	return
}

// SQL adds an arbitrary sql to current position.
func (sb *SelectBuilder) SQL(sql string) *SelectBuilder {
	sb.injection.SQL(sb.marker, sql)
	return sb
}
