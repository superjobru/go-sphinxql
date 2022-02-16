// Copyright 2018 Huan Du. All rights reserved.
// Copyright 2022 OOO SuperJob. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package sphinxql

import (
	"fmt"
	"strconv"
	"time"
	"unicode/utf8"
	"unsafe"
)

// sphinxSearchInterpolate parses query and replace all "?" with encoded args.
// If there are more "?" than len(args), returns ErrMissingArgs.
// Otherwise, if there are less "?" than len(args), the redundant args are omitted.
func sphinxSearchInterpolate(query string, args ...interface{}) (string, error) {
	return sphinxSearchLikeInterpolate(SphinxSearch, query, args...)
}

func sphinxSearchLikeInterpolate(flavor Flavor, query string, args ...interface{}) (string, error) {
	// Roughly estimate the size to avoid useless memory allocation and copy.
	buf := make([]byte, 0, len(query)+len(args)*20)

	var quote rune
	var err error
	cnt := 0
	max := len(args)
	escaping := false
	offset := 0
	target := query
	r, sz := utf8.DecodeRuneInString(target)

	for ; sz != 0; r, sz = utf8.DecodeRuneInString(target) {
		offset += sz
		target = query[offset:]

		if escaping {
			escaping = false
			continue
		}

		switch r {
		case '?':
			if quote != 0 {
				continue
			}

			if cnt >= max {
				return "", ErrInterpolateMissingArgs
			}

			buf = append(buf, query[:offset-sz]...)
			buf, err = encodeValue(buf, args[cnt], flavor)

			if err != nil {
				return "", err
			}

			query = target
			offset = 0
			cnt++

		case '\'':
			if quote == '\'' {
				quote = 0
				continue
			}

			if quote == 0 {
				quote = '\''
			}

		case '"':
			if quote == '"' {
				quote = 0
				continue
			}

			if quote == 0 {
				quote = '"'
			}

		case '`':
			if quote == '`' {
				quote = 0
				continue
			}

			if quote == 0 {
				quote = '`'
			}

		case '\\':
			if quote != 0 {
				escaping = true
			}
		}
	}

	buf = append(buf, query...)
	return *(*string)(unsafe.Pointer(&buf)), nil
}

func encodeValue(buf []byte, arg interface{}, flavor Flavor) ([]byte, error) {
	switch v := arg.(type) {
	case nil:
		buf = append(buf, "NULL"...)

	case bool:
		if v {
			buf = append(buf, "TRUE"...)
		} else {
			buf = append(buf, "FALSE"...)
		}

	case int:
		buf = strconv.AppendInt(buf, int64(v), 10)

	case int8:
		buf = strconv.AppendInt(buf, int64(v), 10)

	case int16:
		buf = strconv.AppendInt(buf, int64(v), 10)

	case int32:
		buf = strconv.AppendInt(buf, int64(v), 10)

	case int64:
		buf = strconv.AppendInt(buf, v, 10)

	case uint:
		buf = strconv.AppendUint(buf, uint64(v), 10)

	case uint8:
		buf = strconv.AppendUint(buf, uint64(v), 10)

	case uint16:
		buf = strconv.AppendUint(buf, uint64(v), 10)

	case uint32:
		buf = strconv.AppendUint(buf, uint64(v), 10)

	case uint64:
		buf = strconv.AppendUint(buf, v, 10)

	case float32:
		buf = strconv.AppendFloat(buf, float64(v), 'g', -1, 32)

	case float64:
		buf = strconv.AppendFloat(buf, v, 'g', -1, 64)

	case []byte:
		if v == nil {
			buf = append(buf, "NULL"...)
			break
		}

		buf = append(buf, "_binary"...)
		buf = quoteStringValue(buf, *(*string)(unsafe.Pointer(&v)), flavor)

	case string:
		buf = quoteStringValue(buf, v, flavor)

	case time.Time:
		if v.IsZero() {
			buf = append(buf, "'0000-00-00'"...)
			break
		}

		// In SQL standard, the precision of fractional seconds in time literal is up to 6 digits.
		// Round up v.
		v = v.Add(500 * time.Nanosecond)
		buf = append(buf, '\'')

		buf = append(buf, v.Format("2006-01-02 15:04:05.999999")...)

		buf = append(buf, '\'')

	case fmt.Stringer:
		buf = quoteStringValue(buf, v.String(), flavor)

	default:
		return nil, ErrInterpolateUnsupportedArgs
	}

	return buf, nil
}

func quoteStringValue(buf []byte, s string, _ Flavor) []byte {
	buf = append(buf, '\'')
	r, sz := utf8.DecodeRuneInString(s)

	for ; sz != 0; r, sz = utf8.DecodeRuneInString(s) {
		switch r {
		case '\x00':
			buf = append(buf, "\\0"...)

		case '\b':
			buf = append(buf, "\\b"...)

		case '\n':
			buf = append(buf, "\\n"...)

		case '\r':
			buf = append(buf, "\\r"...)

		case '\t':
			buf = append(buf, "\\t"...)

		case '\x1a':
			buf = append(buf, "\\Z"...)

		case '\'':
			buf = append(buf, "\\'"...)

		case '"':
			buf = append(buf, "\\\""...)

		case '\\':
			buf = append(buf, "\\\\"...)

		default:
			buf = append(buf, s[:sz]...)
		}

		s = s[sz:]
	}

	buf = append(buf, '\'')
	return buf
}
