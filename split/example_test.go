package split

import (
	"net"

	"github.com/bgentry/pqarray/split"
)

func ExampleArray() {
	asBytes := []byte("{199.27.128.0/21,173.245.48.0/20,2400:cb00::/32}")
	parts, err := split.Array(asBytes)
	if err != nil {
		return err
	}
	res := make([]net.IPNet, len(parts))
	for i := range parts {
		_, ipn, err := net.ParseCIDR(string(parts[i]))
		if err != nil {
			return err
		}
		res[i] = *ni
	}
}
