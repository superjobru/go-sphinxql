package sphinxql

import (
	"testing"
	"time"

	"github.com/huandu/go-assert"
)

func TestFlavorInterpolate(t *testing.T) {
	a := assert.New(t)
	dt := time.Date(2019, 4, 24, 12, 23, 34, 123456789, time.FixedZone("CST", 8*60*60)) // 2019-04-24 12:23:34.987654321 CST
	cases := []struct {
		flavor Flavor
		sql    string
		args   []interface{}
		query  string
		err    error
	}{
		{
			SphinxQL,
			"SELECT * FROM a WHERE name = ? AND state IN (?, ?, ?, ?, ?)", []interface{}{"I'm fine", 42, int8(8), int16(-16), int32(32), int64(64)},
			"SELECT * FROM a WHERE name = 'I\\'m fine' AND state IN (42, 8, -16, 32, 64)", nil,
		},
		{
			SphinxQL,
			"SELECT * FROM `a?` WHERE name = \"?\" AND state IN (?, '?', ?, ?, ?, ?, ?)", []interface{}{"\r\n\b\t\x1a\x00\\\"'", uint(42), uint8(8), uint16(16), uint32(32), uint64(64), "useless"},
			"SELECT * FROM `a?` WHERE name = \"?\" AND state IN ('\\r\\n\\b\\t\\Z\\0\\\\\\\"\\'', '?', 42, 8, 16, 32, 64)", nil,
		},
		{
			SphinxQL,
			"SELECT ?, ?, ?, ?, ?, ?, ?, ?, ?", []interface{}{true, false, float32(1.234567), float64(9.87654321), []byte(nil), []byte("I'm bytes"), dt, time.Time{}, nil},
			"SELECT TRUE, FALSE, 1.234567, 9.87654321, NULL, _binary'I\\'m bytes', '2019-04-24 12:23:34.123457', '0000-00-00', NULL", nil,
		},
		{
			SphinxQL,
			"SELECT '\\'?', \"\\\"?\", `\\`?`, \\?", []interface{}{SphinxQL},
			"SELECT '\\'?', \"\\\"?\", `\\`?`, \\'SphinxQL'", nil,
		},
		{
			SphinxQL,
			"SELECT ?", nil,
			"", ErrInterpolateMissingArgs,
		},
		{
			SphinxQL,
			"SELECT ?", []interface{}{complex(1, 2)},
			"", ErrInterpolateUnsupportedArgs,
		},
	}

	for idx, c := range cases {
		a.Use(&idx, &c)
		query, err := c.flavor.Interpolate(c.sql, c.args)

		a.Equal(query, c.query)
		a.Assert(err == c.err || err.Error() == c.err.Error())
	}
}
