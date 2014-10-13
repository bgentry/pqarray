/*
Package split contains PostgreSQL array parsing logic from Chris Farmiloe
(from package bitbucket.org/pkg/pql).

You can use the func Array to split a byte slice into its elements. They will
be escaped if necessary.

You can look at the main pqarray package for examples of how to use this
function. Here's an example of how to apply it to another type to make that
type scannable by the postgres driver:

	type ipnets []net.IPNet

	// Scan implements sql.Scanner for the cidrSlice type.
	func (ipns *ipnets) Scan(src interface{}) error {
		if src == nil {
			return nil
		}
		asBytes, ok := src.([]byte)
		if !ok {
			return errors.New("Scan source was not []bytes")
		}

		parts, err := split.Array(asBytes)
		if err != nil {
			return err
		}
		res := make(ipnets, len(parts))
		for i := range parts {
			ni := net.IPNet{}
			if err := ni.UnmarshalJSON(parts[i]); err != nil {
				return err
			}
			res[i] = *ni
		}
		(*ipns) = res

		return nil
	}
*/
package split
