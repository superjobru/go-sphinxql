// Copyright 2018 Huan Du. All rights reserved.
// Copyright 2022 OOO SuperJob. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package sphinxql

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

type Opt struct {
	Args *Args
}

// NamedIntegerList represents named integer list for options.
type NamedIntegerList map[string]int

// RankerOptionValue is an alias of UnquotedString.
type RankerOptionValue = UnquotedString

// RankerOptionValue enum
const (
	RankerProximityBM25 RankerOptionValue = "proximity_bm25"
	RankerBM25                            = "bm25"
	RankerNone                            = "none"
	RankerWordCount                       = "wordcount"
	RankerProximity                       = "proximity"
	RankerMatchAny                        = "matchany"
	RankerFieldMask                       = "fieldmask"
	RankerSPH04                           = "sph04"
	RankerExpr                            = "expr"
	RankerExport                          = "export"
)

// Comment builds a comment OPTION.
func (o *Opt) Comment(value string) string {
	return fmt.Sprintf("comment = %s", o.Args.Add(value))
}

// FieldWeights builds a field_weights OPTION.
func (o *Opt) FieldWeights(values NamedIntegerList) string {
	buf := &bytes.Buffer{}

	buf.WriteString("(")

	ks := make([]string, len(values))

	i := 0
	for k := range values {
		ks[i] = k
		i++
	}

	sort.Strings(ks)

	var nl []string
	for _, v := range ks {
		nl = append(nl, fmt.Sprintf("%s=%d", v, values[v]))
	}

	buf.WriteString(strings.Join(nl, ", "))

	buf.WriteString(")")

	return fmt.Sprintf("field_weights = %s", o.Args.Add(buf.String()))
}

// MaxMatches builds a max_matches OPTION.
func (o *Opt) MaxMatches(value int) string {
	return fmt.Sprintf("max_matches = %s", o.Args.Add(value))
}

// Ranker builds a ranker OPTION.
func (o *Opt) Ranker(value RankerOptionValue) string {
	return fmt.Sprintf("ranker = %s", o.Args.Add(value))
}

func (o *Opt) exprRanker(value RankerOptionValue, expr string) string {
	return fmt.Sprintf("ranker = %s(%s)", value, o.Args.Add(expr))
}

// ExprRanker builds a ranker = expr(expr) OPTION.
func (o *Opt) ExprRanker(expr string) string {
	return o.exprRanker(RankerExpr, expr)
}

// ExportRanker builds a ranker = export(expr) OPTION.
func (o *Opt) ExportRanker(expr string) string {
	return o.exprRanker(RankerExport, expr)
}
