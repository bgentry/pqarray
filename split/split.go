// Copyright (C) 2013 Chris Farmiloe
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package split

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
			case closer == '}' && b == '}':
				// mark end of array
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
