// Copyright 2018 Huan Du. All rights reserved.
// Copyright 2022 OOO SuperJob. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package sphinxql

import (
	"testing"

	"github.com/huandu/go-assert"
)

func TestOrdBy(t *testing.T) {
	a := assert.New(t)
	cases := map[string]func() string{
		"kekw":                 func() string { return newTestOrdBy().NoDir("kekw") },
		"memes ASC":            func() string { return newTestOrdBy().Asc("memes") },
		"oof DESC":             func() string { return newTestOrdBy().Desc("oof") },
		"kekw, wkek":           func() string { return newTestOrdBy().NoDirMulti([]string{"kekw", "wkek"}) },
		"memes ASC, jokes ASC": func() string { return newTestOrdBy().AscMulti([]string{"memes", "jokes"}) },
		"oof DESC, foo DESC":   func() string { return newTestOrdBy().DescMulti([]string{"oof", "foo"}) },
	}

	for expected, f := range cases {
		actual := f()
		a.Equal(actual, expected)
	}
}

func newTestOrdBy() *OrdBy {
	return &OrdBy{}
}
