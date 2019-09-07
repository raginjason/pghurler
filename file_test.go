package main

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestNewDelimitedFile(t *testing.T) {
	tests := map[string]struct {
		path      string
		delimiter rune
		want      *DelimitedFile
		err       error
	}{
		"csv":  {"foo.csv", '\x00', &DelimitedFile{',', "foo.csv"}, nil},
		"tsv":  {"foo.tsv", '\x00', &DelimitedFile{'\t', "foo.tsv"}, nil},
		"tab":  {"foo.tab", '\x00', &DelimitedFile{'\t', "foo.tab"}, nil},
		"pipe": {"foo.pipe", '\x00', &DelimitedFile{'|', "foo.pipe"}, nil},

		"no ext":      {"foo", '\x00', nil, errors.New("no extension to derive delimiter from")},
		"unknown ext": {"foo.foo", '\x00', nil, errors.New("could not derive delimiter from '.foo' extension")},
	}

	equateErrorMessage := cmp.Comparer(func(x, y error) bool {
		if x == nil || y == nil {
			return x == nil && y == nil
		}
		return x.Error() == y.Error()
	})

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			f, err := NewDelimitedFile(tc.path, tc.delimiter)

			if diff := cmp.Diff(tc.want, f); diff != "" {
				t.Errorf("NewDelimitedFile(%s, %q) mismatch (-want +got):\n%s", tc.path, tc.delimiter, diff)
			}

			if diff := cmp.Diff(tc.err, err, equateErrorMessage); diff != "" {
				t.Errorf("Error mismatch for NewDelimitedFile(%s, %q) (-want +got):\n%s", tc.path, tc.delimiter, diff)
			}
		})
	}
}
