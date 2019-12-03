/*
Copyright © 2019 Jason Walker <ragin.jason@me.com>
This file is part of pghurler.
*/
package main

import (
	"encoding/csv"
)

type Reader struct {
	columns      []string
	reader       *csv.Reader
	lineNumber   uint64
	recordNumber uint64
}

type ReaderOptions struct {
	Columns   []string
	SkipLines uint
}

type Record struct {
	RecordNumber uint64
	LineNumber   uint64
	Columns      []string
	Data         []string
}

func NewReader(csv *csv.Reader, opts *ReaderOptions) (*Reader, error) {

	hr := &Reader{reader: csv}

	// TODO process opts.SkipLines

	if opts != nil && opts.Columns != nil {
		hr.columns = make([]string, len(opts.Columns))
		copy(hr.columns, opts.Columns)
	} else {
		var err error
		hr.columns, err = csv.Read()
		if err != nil {
			return nil, err
		}
	}

	return hr, nil
}

func (r *Reader) Read() (*Record, error) {

	rec, err := r.reader.Read()
	if err != nil {
		return nil, err
	}

	r.lineNumber = r.lineNumber + 1
	r.recordNumber = r.recordNumber + 1
	hr := &Record{RecordNumber: r.recordNumber, LineNumber: r.lineNumber, Columns: r.columns, Data: rec}

	return hr, nil
}
