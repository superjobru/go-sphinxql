package sphinxql

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type OptionBuilder struct{}

func (ob *OptionBuilder) Comment(v CommentOptionValue) *CommentOption {
	return &CommentOption{value: v}
}

func (ob *OptionBuilder) FieldWeights(v FieldWeightsOptionValues) *FieldWeightsOption {
	return &FieldWeightsOption{values: v}
}

func (ob *OptionBuilder) MaxMatches(v MaxMatchesOptionValue) *MaxMatchesOption {
	return &MaxMatchesOption{value: v}
}

func (ob *OptionBuilder) Ranker(v RankerOptionValue) *RankerOption {
	return &RankerOption{value: v}
}

type Option interface {
	GetName() string
	GetValue() string
	Serialize() string
}

func serialize(o Option) string {
	return fmt.Sprintf("%s = %s", o.GetName(), o.GetValue())
}

func SerializeOptions(os []Option) string {
	var ss []string

	for _, v := range os {
		ss = append(ss, v.Serialize())
	}

	return strings.Join(ss, ", ")
}

type CommentOption struct {
	Option
	value CommentOptionValue
}

type CommentOptionValue string

func (co *CommentOption) GetName() string {
	return "comment"
}

func (co *CommentOption) GetValue() string {
	return string(co.value)
}

func (co *CommentOption) Serialize() string {
	return serialize(co)
}

type FieldWeightsOption struct {
	Option
	values FieldWeightsOptionValues
}

type FieldWeightsOptionValues map[string]int

func (fwo *FieldWeightsOption) GetName() string {
	return "field_weights"
}

func (fwo *FieldWeightsOption) GetValue() string {
	buf := &bytes.Buffer{}

	buf.WriteString("(")

	ks := make([]string, len(fwo.values))

	i := 0
	for k := range fwo.values {
		ks[i] = k
		i++
	}

	sort.Strings(ks)

	var nl []string
	for _, v := range ks {
		nl = append(nl, fmt.Sprintf("%s=%d", v, fwo.values[v]))
	}

	buf.WriteString(strings.Join(nl, ", "))

	buf.WriteString(")")

	return buf.String()
}

func (fwo *FieldWeightsOption) Serialize() string {
	return serialize(fwo)
}

type MaxMatchesOption struct {
	Option
	value MaxMatchesOptionValue
}

type MaxMatchesOptionValue int

func (mmo *MaxMatchesOption) GetName() string {
	return "max_matches"
}

func (mmo *MaxMatchesOption) GetValue() string {
	return strconv.Itoa(int(mmo.value))
}

func (mmo *MaxMatchesOption) Serialize() string {
	return serialize(mmo)
}

type RankerOption struct {
	Option
	value RankerOptionValue
}

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

func (ro *RankerOption) GetName() string {
	return "ranker"
}

func (ro *RankerOption) GetValue() string {
	return string(ro.value)
}

func (ro *RankerOption) Serialize() string {
	return serialize(ro)
}
