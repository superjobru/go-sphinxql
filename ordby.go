// Copyright 2018 Huan Du. All rights reserved.
// Copyright 2022 OOO SuperJob. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package sphinxql

import (
	"fmt"
	"strings"
)

// OrdBy provides several helper methods to build ORDER BY clause.
type OrdBy struct{}

func (o *OrdBy) orderBy(value string, dir string) string {
	return fmt.Sprintf("%s%s", value, dir)
}

func (o *OrdBy) orderByMulti(values []string, dir string) string {
	if len(values) == 0 {
		return ""
	}

	return fmt.Sprintf("%s%s", strings.Join(values, fmt.Sprintf("%s, ", dir)), dir)
}

// NoDir builds an ORDER BY clause with no direction specified.
func (o *OrdBy) NoDir(value string) string {
	return o.orderBy(value, "")
}

// NoDirMulti builds an ORDER BY clause with multiple columns and no direction specified.
func (o *OrdBy) NoDirMulti(values []string) string {
	return o.orderByMulti(values, "")
}

// Asc builds an ORDER BY ASC clause.
func (o *OrdBy) Asc(value string) string {
	return o.orderBy(value, " ASC")
}

// AscMulti builds an ORDER BY ASC clause with multiple columns.
func (o *OrdBy) AscMulti(values []string) string {
	return o.orderByMulti(values, " ASC")
}

// Desc builds an ORDER BY DESC clause.
func (o *OrdBy) Desc(value string) string {
	return o.orderBy(value, " DESC")
}

// DescMulti builds an ORDER BY DESC clause with multiple columns.
func (o *OrdBy) DescMulti(values []string) string {
	return o.orderByMulti(values, " DESC")
}
