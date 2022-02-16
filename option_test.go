package sphinxql

import (
	"testing"

	"github.com/huandu/go-assert"
)

func TestOptionBuilder(t *testing.T) {
	a := assert.New(t)
	cases := map[string]func() string{
		"comment = comment": func() string {
			return newTestOptionBuilder().Comment("comment").Serialize()
		},

		"field_weights = (first_field=10, second_field=20, third_field=30)": func() string {
			return newTestOptionBuilder().
				FieldWeights(FieldWeightsOptionValues{
					"first_field":  10,
					"second_field": 20,
					"third_field":  30,
				}).
				Serialize()
		},

		"max_matches = 100": func() string {
			return newTestOptionBuilder().MaxMatches(100).Serialize()
		},

		"ranker = wordcount": func() string {
			return newTestOptionBuilder().Ranker(RankerWordCount).Serialize()
		},
	}

	for expected, f := range cases {
		actual := f()
		a.Equal(actual, expected)
	}
}

func newTestOptionBuilder() *OptionBuilder {
	return &OptionBuilder{}
}
