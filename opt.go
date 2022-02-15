// Copyright 2018 Huan Du. All rights reserved.
// Copyright 2022 OOO SuperJob. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package sphinxql

import (
	"fmt"
	"sort"
	"strings"
)

// Opt provides several helper methods to build options.
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
	RankerBM25          RankerOptionValue = "bm25"
	RankerNone          RankerOptionValue = "none"
	RankerWordCount     RankerOptionValue = "wordcount"
	RankerProximity     RankerOptionValue = "proximity"
	RankerMatchAny      RankerOptionValue = "matchany"
	RankerFieldMask     RankerOptionValue = "fieldmask"
	RankerSPH04         RankerOptionValue = "sph04"
	RankerExpr          RankerOptionValue = "expr"
	RankerExport        RankerOptionValue = "export"
)

// Comment builds a comment OPTION.
func (o *Opt) Comment(value string) string {
	return fmt.Sprintf("comment = %s", o.Args.Add(value))
}

// FieldWeights builds a field_weights OPTION.
func (o *Opt) FieldWeights(values NamedIntegerList) string {
	buf := &strings.Builder{}

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
