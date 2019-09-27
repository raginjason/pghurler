package main

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"io"
	"testing"
)

func TestDeriveDelimiter(t *testing.T) {
	tests := map[string]struct {
		filepath string
		want     rune
		err      error
	}{
		"csv":  {"foo.csv", ',', nil},
		"tsv":  {"foo.tsv", '\t', nil},
		"tab":  {"foo.tab", '\t', nil},
		"pipe": {"foo.pipe", '|', nil},

		"no ext":      {"foo", 0, errors.New("no extension to derive delimiter from")},
		"unknown ext": {"foo.foo", 0, errors.New("could not derive delimiter from '.foo' extension")},
	}

	equateErrorMessage := cmp.Comparer(func(x, y error) bool {
		if x == nil || y == nil {
			return x == nil && y == nil
		}
		return x.Error() == y.Error()
	})

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			d, err := deriveDelimiter(tc.filepath)

			if diff := cmp.Diff(tc.want, d); diff != "" {
				t.Errorf("deriveDelimiter(%q) mismatch (-want +got):\n%s", tc.filepath, diff)
			}

			if diff := cmp.Diff(tc.err, err, equateErrorMessage); diff != "" {
				t.Errorf("Error mismatch for deriveDelimiter(%q) (-want +got):\n%s", tc.filepath, diff)
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

func TestNewDelimitedFile(t *testing.T) {
	tests := map[string]struct {
		conf DelimitedFileConfig
		want *DelimitedFile
		err  error
	}{
		"defined delim":     {DelimitedFileConfig{Filepath: "foo.csv", Delimiter: '.'}, &DelimitedFile{Delimiter: '.', Filepath: "foo.csv"}, nil},
		"derivable delim":   {DelimitedFileConfig{Filepath: "foo.csv"}, &DelimitedFile{Delimiter: ',', Filepath: "foo.csv"}, nil},
		"underivable delim": {DelimitedFileConfig{Filepath: "foo.foo"}, nil, errors.New("could not derive delimiter from '.foo' extension")},
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
