// Copyright 2018 Huan Du. All rights reserved.
// Copyright 2022 OOO SuperJob. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package sphinxql

import (
	"fmt"
)

// OrdBy provides several helper methods to build ORDER BY clause.
type OrdBy struct{}

func (o *OrdBy) orderBy(value string, dir string) string {
	return fmt.Sprintf("%s%s", value, dir)
}

// NoDir builds an ORDER BY clause with no direction specified.
func (o *OrdBy) NoDir(value string) string {
	return o.orderBy(value, "")
}

// Asc builds an ORDER BY ASC clause.
func (o *OrdBy) Asc(value string) string {
	return o.orderBy(value, " "+"ASC")
}

// Desc builds an ORDER BY DESC clause.
func (o *OrdBy) Desc(value string) string {
	return o.orderBy(value, " "+"DESC")
}
