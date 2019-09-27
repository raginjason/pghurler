package main

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"io"
	"testing"
)

func TestNewDelimitedFile(t *testing.T) {
	tests := map[string]struct {
		conf DelimitedFileConfig
		want *DelimitedFile
		err  error
	}{
		"csv":  {DelimitedFileConfig{Filepath: "foo.csv", Delimiter: '\x00'}, &DelimitedFile{Delimiter: ',', Filepath: "foo.csv"}, nil},
		"tsv":  {DelimitedFileConfig{Filepath: "foo.tsv", Delimiter: '\x00'}, &DelimitedFile{Delimiter: '\t', Filepath: "foo.tsv"}, nil},
		"tab":  {DelimitedFileConfig{Filepath: "foo.tab", Delimiter: '\x00'}, &DelimitedFile{Delimiter: '\t', Filepath: "foo.tab"}, nil},
		"pipe": {DelimitedFileConfig{Filepath: "foo.pipe", Delimiter: '\x00'}, &DelimitedFile{Delimiter: '|', Filepath: "foo.pipe"}, nil},

		"no ext":      {DelimitedFileConfig{Filepath: "foo", Delimiter: '\x00'}, nil, errors.New("no extension to derive delimiter from")},
		"unknown ext": {DelimitedFileConfig{Filepath: "foo.foo", Delimiter: '\x00'}, nil, errors.New("could not derive delimiter from '.foo' extension")},
	}

	equateErrorMessage := cmp.Comparer(func(x, y error) bool {
		if x == nil || y == nil {
			return x == nil && y == nil
		}
		return x.Error() == y.Error()
	})

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			f, err := NewDelimitedFile(tc.conf)

			if diff := cmp.Diff(tc.want, f); diff != "" {
				t.Errorf("NewDelimitedFile(%q) mismatch (-want +got):\n%s", tc.conf, diff)
			}

			if diff := cmp.Diff(tc.err, err, equateErrorMessage); diff != "" {
				t.Errorf("Error mismatch for NewDelimitedFile(%q) (-want +got):\n%s", tc.conf, diff)
			}
		})
	}
}

func TestParseHeader(t *testing.T) {
	tests := map[string]struct {
		delim  rune
		header string
		want   []string
		err    error
	}{
		"csv":        {',', "foo,bar,baz", []string{"foo", "bar", "baz"}, nil},
		"csv quotes": {',', `foo,"bar",baz`, []string{"foo", "bar", "baz"}, nil},
		"empty":      {',', "", []string(nil), io.EOF},
	}

	equateErrorMessage := cmp.Comparer(func(x, y error) bool {
		if x == nil || y == nil {
			return x == nil && y == nil
		}
		return x.Error() == y.Error()
	})

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			c, err := parseHeader(tc.delim, tc.header)

			if diff := cmp.Diff(tc.want, c); diff != "" {
				t.Errorf("parseHeader(%q, %s) mismatch (-want +got):\n%s", tc.delim, tc.header, diff)
			}

			if diff := cmp.Diff(tc.err, err, equateErrorMessage); diff != "" {
				t.Errorf("Error mismatch for parseHeader(%q, %s) (-want +got):\n%s", tc.delim, tc.header, diff)
			}
		})
	}
}
