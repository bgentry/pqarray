package pqarray

import "testing"

var stringsScanTests = []struct {
	input []byte
	want  Strings
	err   bool
}{
	{
		input: []byte("{foo,bar}"),
		want:  Strings{"foo", "bar"},
	},
	{
		input: []byte("{foo, bar}"),
		want:  Strings{"foo", "bar"},
	},
	{
		input: []byte(`{foo,"bar"," baz"}`),
		want:  Strings{"foo", "bar", " baz"},
	},
	{
		input: []byte(`{"\\ ","\"\"\"",bar," baz"}`),
		want:  Strings{`\ `, `"""`, "bar", " baz"},
	},
	// TODO: multi-dimensional arrays
	{
		input: []byte("{199.27.128.0/21, 173.245.48.0/20, 2400:cb00::/32}"),
		want:  Strings{"199.27.128.0/21", "173.245.48.0/20", "2400:cb00::/32"},
	},
	{
		input: []byte("{}"),
		want:  Strings{},
	},
	{
		input: []byte(""),
		err:   true,
	},
	{
		input: []byte("{"),
		err:   true,
	},
	{
		input: []byte("}}"),
		err:   true,
	},
}

func TestStringsScan(t *testing.T) {
	for _, test := range stringsScanTests {
		var got Strings
		if err := got.Scan(test.input); err != nil {
			switch {
			case test.err:
			default:
				t.Error(err)
			}
			continue
		}
		if wantLen, gotLen := len(test.want), len(got); wantLen != gotLen {
			t.Errorf("input=%q want len=%d, got %d", test.input, wantLen, gotLen)
			continue
		}
		for i, want := range test.want {
			if got[i] != want {
				t.Errorf("input=%q element %d want %q, got %q", test.input, i, want, got[i])
			}
		}
	}
}

var intsScanTests = []struct {
	input []byte
	want  Ints
	err   bool
}{
	{
		input: []byte("{1,2}"),
		want:  Ints{1, 2},
	},
	{
		input: []byte(`{1, "2"}`),
		want:  Ints{1, 2},
	},
	{
		input: []byte("{ 1, 2 }"),
		err:   true,
	},
	{
		input: []byte(`{1, " 2"}`),
		err:   true,
	},
	{
		input: []byte(`{1, a}`),
		err:   true,
	},
	// TODO: multi-dimensional arrays
}

func TestIntsScan(t *testing.T) {
	for _, test := range intsScanTests {
		var got Ints
		if err := got.Scan(test.input); err != nil {
			switch {
			case test.err:
			default:
				t.Errorf("input=%q %v", test.input, err)
			}
			continue
		}
		if wantLen, gotLen := len(test.want), len(got); wantLen != gotLen {
			t.Errorf("input=%q want len=%d, got %d", test.input, wantLen, gotLen)
			continue
		}
		for i, want := range test.want {
			if got[i] != want {
				t.Errorf("input=%q element %d want %d, got %d", test.input, i, want, got[i])
			}
		}
	}
}
