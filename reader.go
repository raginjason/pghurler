/*
Copyright © 2019 Jason Walker <ragin.jason@me.com>
This file is part of pghurler.
*/
package main

import (
	"encoding/csv"
	"io"
)

type Reader struct {
	reader        *csv.Reader
	currentLine   uint64
	currentRecord uint64
	columns       []string
}

type Record struct {
	LineNumber   uint64
	RecordNumber uint64
	Values       map[string]string
}

func NewReader(csv *csv.Reader) (*Reader, error) {

	r := &Reader{reader: csv}

	var err error
	r.columns, err = csv.Read()
	if err != nil {
		return nil, err
	}
	r.currentLine++

	return r, nil
}

func (r *Reader) Read() (*Record, error) {

	rec, err := r.reader.Read()
	if err != nil {
		return nil, err
	}

	values := make(map[string]string)
	for i, v := range r.columns {
		values[v] = rec[i]
	}

	r.currentRecord++
	r.currentLine++
	outRec := &Record{RecordNumber: r.currentRecord, LineNumber: r.currentLine, Values: values}
	return outRec, nil
}

func (r *Reader) ReadAll() (records []*Record, err error) {
	for {
		record, err := r.Read()
		if err == io.EOF {
			return records, nil
		}
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
}
