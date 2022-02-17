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

type NamedIntegerList map[string]int

type RankerOptionValue string

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

func (o *Opt) Comment(value string) string {
	return fmt.Sprintf("comment = '%s'", o.Args.Add(value))
}

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

func (o *Opt) MaxMatches(value int) string {
	return fmt.Sprintf("max_matches = %s", o.Args.Add(value))
}

func (o *Opt) Ranker(value RankerOptionValue) string {
	return fmt.Sprintf("ranker = %s", o.Args.Add(value))
}
