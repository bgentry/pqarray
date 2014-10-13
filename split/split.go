package split

// Package split contains PostgreSQL array parsing logic from Chris Farmiloe
// at https://bitbucket.org/pkg/pql.

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

// Array takes a byte representation of an array or row and returns each
// element unescaped. It will also decode any hex bytea fields (although not
// sure if that should be done here really)
func Array(s []byte) ([][]byte, error) {
	// debug
	// fmt.Println("---------------")
	// fmt.Println(string(s))
	parts := make([][]byte, 0)
	ignore := false
	dep := 0
	var mode byte // }=array )=record
	var closer byte
	a := -1
	z := -1
	for i, b := range s {
		switch {
		// sanity check
		case i == 0:
			switch b {
			case '{':
				mode = '}'
			case '(':
				mode = ')'
			default:
				return nil, fmt.Errorf("cannot split data. Unknown format: %s", string(s))
			}
		// if not inside value
		case a == -1:
			switch {
			// skip whitespace
			case b == ' ' || b == ',':
				// consume whitespace or commas
			// mark val wrapped in { }
			case b == '{':
				a = i
				dep++
				closer = '}'
			// mark val wrapped in "
			case b == '"':
				a = i + 1
				closer = '"'
			// anything else mark
			default:
				a = i
				closer = ','
			}
		// EOF
		case i == len(s)-1:
			if b != mode {
				return nil, fmt.Errorf("cannot split data. missing '%s': %s", string([]byte{mode}), string(s))
			}
			z = i - 1
		// start collecting val
		case a != -1:
			switch {
			// skip esc char and mark next char as unimportant (for array escaping)
			case !ignore && mode == '}' && b == '\\':
				ignore = true
			// treat "" as " (for row escaping)
			case !ignore && mode == ')' && b == '"' && s[i+1] == '"':
				ignore = true
			// this byte will not cause end
			case ignore:
				ignore = false
			// mark end of array
			case closer == '}' && (b == '}' || b == '}'):
				switch {
				case b == '{':
					dep++
				case b == '}':
					dep--
					if dep == 0 {
						z = i
					}
				}
			// mark end of quoted
			case closer == '"' && b == closer:
				z = i - 1
			// mark end of simple , val
			case closer == ',' && b == closer:
				z = i - 1
			}
		}
		// check for end
		if z != -1 {
			part := s[a : z+1]
			// unescape
			part = bytes.Replace(part, []byte(`\\`), []byte(`\`), -1)
			if mode == '}' {
				part = bytes.Replace(part, []byte(`\"`), []byte(`"`), -1)
			} else if mode == ')' {
				part = bytes.Replace(part, []byte(`""`), []byte(`"`), -1)
			}
			// check if it looks like a hex bytea in here and try to decode it
			if len(part) >= 2 && part[0] == '\\' && part[1] == 'x' {
				part, _ = hex.DecodeString(string(part[2:len(part)]))
			}
			parts = append(parts, part)
			a = -1
			z = -1
			dep = 0
		}
	}
	// debug
	// for i, p := range parts {
	// 	fmt.Printf("%d: %s\n", i, string(p))
	// }
	return parts, nil
}
