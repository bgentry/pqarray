// Package pqarray provides helper types for parsing PostgreSQL arrays. It is
// compatible with pq (github.com/lib/pq).
//
// The split subpackage exports a function to use the array parser with an
// arbitrary array type.
package pqarray

import (
	"errors"
	"strconv"

	"github.com/bgentry/pqarray/split"
)

// Strings is a string slice that implements sql.Scanner
type Strings []string

// Scan implements sql.Scanner.
func (s *Strings) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return errors.New("scan source was not []bytes")
	}

	parts, err := split.Array(asBytes)
	if err != nil {
		return err
	}
	res := make(Strings, len(parts))
	for i := range parts {
		res[i] = string(parts[i])
	}
	(*s) = res

	return nil
}

// Ints is a int slice that implements sql.Scanner.
type Ints []int

// Scan implements sql.Scanner.
func (in *Ints) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return errors.New("scan source was not []bytes")
	}

	parts, err := split.Array(asBytes)
	if err != nil {
		return err
	}
	res := make(Ints, len(parts))
	for i := range parts {
		v, err := strconv.Atoi(string(parts[i]))
		if err != nil {
			return err
		}
		res[i] = v
	}
	(*in) = res

	return nil
}
