// Copyright 2018 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package sphinxql

import (
	"fmt"
	"testing"

	"github.com/huandu/go-assert"
)

func TestArgs(t *testing.T) {
	a := assert.New(t)
	cases := map[string][]interface{}{
		"abc ? def\n[123]":                   {"abc $? def", 123},
		"abc ? def\n[456]":                   {"abc $0 def", 456},
		"abc  def\n[]":                       {"abc $1 def", 123},
		"abc ? def\n[789]":                   {"abc ${s} def", Named("s", 789)},
		"abc  def \n[]":                      {"abc ${unknown} def ", 123},
		"abc $ def\n[]":                      {"abc $$ def", 123},
		"abcdef$\n[]":                        {"abcdef$", 123},
		"abc ? ? ? ? def\n[123 456 123 456]": {"abc $? $? $0 $? def", 123, 456, 789},
		"abc ? raw ? raw def\n[123 123]":     {"abc $? $? $0 $? def", 123, Raw("raw"), 789},
	}

	for expected, c := range cases {
		args := new(Args)

		for i := 1; i < len(c); i++ {
			args.Add(c[i])
		}

		sql, values := args.Compile(c[0].(string))
		actual := fmt.Sprintf("%v\n%v", sql, values)

		a.Equal(actual, expected)
	}

	old := DefaultFlavor
	defer func() {
		DefaultFlavor = old
	}()
}

func TestArgsAdd(t *testing.T) {
	a := assert.New(t)
	args := &Args{}

	for i := 0; i < maxPredefinedArgs*2; i++ {
		actual := args.Add(i)
		a.Equal(actual, fmt.Sprintf("$%v", i))
	}
}
