/*
Copyright © 2019 Jason Walker <ragin.jason@me.com>
This file is part of pghurler.
*/
package main

import (
	"encoding/csv"
	"errors"
	"github.com/google/go-cmp/cmp"
	"strings"
	"testing"
)

func TestNewReader(t *testing.T) {

	tests := map[string]struct {
		reader *csv.Reader
		opts   *ReaderOptions
		want   *Reader
		err    error
	}{
		"nil opts": {
			csv.NewReader(strings.NewReader("col1,col2")),
			nil,
			&Reader{[]string{"col1", "col2"}, nil, 0, 0},
			nil,
		},
		"opts with col": {
			csv.NewReader(strings.NewReader("val1,val2")),
			&ReaderOptions{[]string{"col1", "col2"}, 0},
			&Reader{[]string{"col1", "col2"}, nil, 0, 0},
			nil,
		},
		"empty csv opts with col": { // Negative test
			csv.NewReader(strings.NewReader("")),
			&ReaderOptions{[]string{"col1", "col2"}, 0},
			&Reader{[]string{"col1", "col2"}, nil, 0, 0},
			nil,
		},
		"empty csv": {
			csv.NewReader(strings.NewReader("")),
			nil,
			nil,
			errors.New("EOF"),
		},
	}

	equateErrorMessage := cmp.Comparer(func(x, y error) bool {
		if x == nil || y == nil {
			return x == nil && y == nil
		}
		return x.Error() == y.Error()
	})

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			r, err := NewReader(tc.reader, tc.opts)

			if diff := cmp.Diff(tc.err, err, equateErrorMessage); diff != "" {
				t.Fatalf("Error mismatch for NewReader() (-want +got):\n%s", diff)
			}

			if r == nil && tc.want != nil {
				t.Fatalf("NewReader() returned nil when it should not be")
			}

			if r != nil && tc.want == nil {
				t.Fatalf("NewReader() return something when it should be nil")
			}

			if tc.want != nil && r != nil {
				if diff := cmp.Diff(tc.want.lineNumber, r.lineNumber); diff != "" {
					t.Errorf("NewReader() lineNumber mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tc.want.recordNumber, r.recordNumber); diff != "" {
					t.Errorf("NewReader() recordNumber mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tc.want.columns, r.columns); diff != "" {
					t.Errorf("NewReader() columns mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
