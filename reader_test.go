/*
Copyright © 2019 Jason Walker <ragin.jason@me.com>
This file is part of pghurler.
*/
package main

import (
	"encoding/csv"
	"errors"
	"github.com/google/go-cmp/cmp"
	"io"
	"strings"
	"testing"
)

func TestNewReader(t *testing.T) {

	tests := map[string]struct {
		reader io.Reader
		want   *Reader
		err    error
	}{
		"empty reader": {
			strings.NewReader(""),
			nil,
			errors.New("EOF"),
		},
		"header-only reader": {
			strings.NewReader("col1,col2"),
			&Reader{nil, 1, 0, []string{"col1", "col2"}},
			nil,
		},
		"header and data reader": {
			strings.NewReader("col1,col2\nval1,val2"),
			&Reader{nil, 1, 0, []string{"col1", "col2"}},
			nil,
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
			r, err := NewReader(csv.NewReader(tc.reader))

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
				if diff := cmp.Diff(tc.want.currentLine, r.currentLine); diff != "" {
					t.Errorf("NewReader() currentLine mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tc.want.currentRecord, r.currentRecord); diff != "" {
					t.Errorf("NewReader() currentRecord mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tc.want.columns, r.columns); diff != "" {
					t.Errorf("NewReader() columns mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
