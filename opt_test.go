// Copyright 2018 Huan Du. All rights reserved.
// Copyright 2022 OOO SuperJob. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package sphinxql

import (
	"testing"

	"github.com/huandu/go-assert"
)

func TestOption(t *testing.T) {
	a := assert.New(t)
	cases := map[string]func() string{
		"comment = $0": func() string { return newTestOption().Comment("kekw") },
		"field_weights = $0": func() string {
			return newTestOption().FieldWeights(NamedIntegerList{
				"first_field":  10,
				"second_field": 20,
				"third_field":  30,
			})
		},
		"max_matches = $0": func() string { return newTestOption().MaxMatches(5) },
		"ranker = $0":      func() string { return newTestOption().Ranker(RankerWordCount) },
	}

	for expected, f := range cases {
		actual := f()
		a.Equal(actual, expected)
	}
}

func newTestOption() *Opt {
	return &Opt{
		Args: &Args{},
	}
}
