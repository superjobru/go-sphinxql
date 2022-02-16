// Copyright 2018 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package sphinxql

import (
	"database/sql"
	"fmt"
)

func ExampleSelect() {
	// Build a SQL to create a HIVE table.
	sql := CreateTable("users").
		SQL("PARTITION BY (year)").
		SQL("AS").
		SQL(
			Select("columns[0] id", "columns[1] name", "columns[2] year").
				From("`all-users.csv`").
				Limit(100).
				String(),
		).
		String()

	fmt.Println(sql)

	// Output:
	// CREATE TABLE users PARTITION BY (year) AS SELECT columns[0] id, columns[1] name, columns[2] year FROM `all-users.csv` LIMIT 100
}

func ExampleSelectBuilder() {
	sb := NewSelectBuilder()
	sb.Select("id", "name", sb.As("COUNT(*)", "t"))
	sb.From("demo.user")
	sb.Where(
		sb.GreaterThan("id", 1234),
		sb.Like("name", "%Du"),
		sb.Or(
			sb.IsNull("id_card"),
			sb.In("status", 1, 2, 5),
		),
		sb.NotIn(
			"id",
			NewSelectBuilder().Select("id").From("banned"),
		), // Nested SELECT.
		"modified_at > created_at + "+sb.Var(86400), // It's allowed to write arbitrary SQL.
	)
	sb.GroupBy("status").Having(sb.NotIn("status", 4, 5))
	sb.OrderBy("modified_at").Asc()
	sb.Limit(10).Offset(5)
	sb.Option(
		sb.Comment("kekw"),
		sb.Ranker(RankerWordCount),
	)

	sql, args := sb.Build()
	fmt.Println(sql)
	fmt.Println(args)

	// Output:
	// SELECT id, name, COUNT(*) AS t FROM demo.user WHERE id > ? AND name LIKE ? AND (id_card IS NULL OR status IN (?, ?, ?)) AND id NOT IN (SELECT id FROM banned) AND modified_at > created_at + ? GROUP BY status HAVING status NOT IN (?, ?) ORDER BY modified_at ASC LIMIT 10 OFFSET 5 OPTION comment = kekw, ranker = wordcount
	// [1234 %Du 1 2 5 86400 4 5]
}

func ExampleSelectBuilder_advancedUsage() {
	sb := NewSelectBuilder()
	innerSb := NewSelectBuilder()

	sb.Select("id", "name")
	sb.From(
		sb.BuilderAs(innerSb, "user"),
	)
	sb.Where(
		sb.In("status", Flatten([]int{1, 2, 3})...),
		sb.Between("created_at", sql.Named("start", 1234567890), sql.Named("end", 1234599999)),
	)
	sb.OrderBy("modified_at").Desc()

	innerSb.Select("*")
	innerSb.From("banned")
	innerSb.Where(
		innerSb.NotIn("name", Flatten([]string{"Huan Du", "Charmy Liu"})...),
	)

	sql, args := sb.Build()
	fmt.Println(sql)
	fmt.Println(args)

	// Output:
	// SELECT id, name FROM (SELECT * FROM banned WHERE name NOT IN (?, ?)) AS user WHERE status IN (?, ?, ?) AND created_at BETWEEN @start AND @end ORDER BY modified_at DESC
	// [Huan Du Charmy Liu 1 2 3 {{} start 1234567890} {{} end 1234599999}]
}

func ExampleSelectBuilder_limit_offset() {
	flavors := []Flavor{SphinxSearch}
	results := make([][]string, len(flavors))
	sb := NewSelectBuilder()
	saveResults := func() {
		for i, f := range flavors {
			sql, _ := sb.BuildWithFlavor(f)
			results[i] = append(results[i], sql)
		}
	}

	sb.Select("*")
	sb.From("user")

	// Case #1: limit < 0 and offset < 0
	//
	// All: No limit or offset in query.
	sb.Limit(-1)
	sb.Offset(-1)
	saveResults()

	// Case #2: limit < 0 and offset >= 0
	//
	// SphinxSearch: Ignore offset if the limit is not set.
	sb.Limit(-1)
	sb.Offset(0)
	saveResults()

	// Case #3: limit >= 0 and offset >= 0
	//
	// All: Set both limit and offset.
	sb.Limit(1)
	sb.Offset(0)
	saveResults()

	// Case #4: limit >= 0 and offset < 0
	//
	// All: Set limit in query.
	sb.Limit(1)
	sb.Offset(-1)
	saveResults()

	for i, result := range results {
		fmt.Println()
		fmt.Println(flavors[i])

		for n, sql := range result {
			fmt.Printf("#%d: %s\n", n+1, sql)
		}
	}

	// Output:
	//
	// SphinxSearch
	// #1: SELECT * FROM user
	// #2: SELECT * FROM user
	// #3: SELECT * FROM user LIMIT 1 OFFSET 0
	// #4: SELECT * FROM user LIMIT 1
}

func ExampleSelectBuilder_varInCols() {
	// Column name may contain some characters, e.g. the $ sign, which have special meanings in builders.
	// It's recommended to call Escape() or EscapeAll() to escape the name.

	sb := NewSelectBuilder()
	v := sb.Var("foo")
	sb.Select(Escape("colHasA$Sign"), v)
	sb.From("table")

	sql, args := sb.Build()
	fmt.Println(sql)
	fmt.Println(args)

	// Output:
	// SELECT colHasA$Sign, ? FROM table
	// [foo]
}

func ExampleSelectBuilder_SQL() {
	sb := NewSelectBuilder()
	sb.SQL("/* before */")
	sb.Select("u.id", "u.name", "c.type", "p.nickname")
	sb.SQL("/* after select */")
	sb.From("user u")
	sb.SQL("/* after from */")
	sb.Where(
		"u.modified_at > u.created_at",
	)
	sb.SQL("/* after where */")
	sb.OrderBy("id")
	sb.SQL("/* after order by */")
	sb.Limit(10)
	sb.SQL("/* after limit */")

	sql := sb.String()
	fmt.Println(sql)

	// Output:
	// /* before */ SELECT u.id, u.name, c.type, p.nickname /* after select */ FROM user u /* after from */ WHERE u.modified_at > u.created_at /* after where */ ORDER BY id /* after order by */ LIMIT 10 /* after limit */
}
