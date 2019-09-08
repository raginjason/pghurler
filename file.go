package main

import (
	"encoding/csv"
	"fmt"
	"path/filepath"
	"strings"
)

// DelimitedFile struct
type DelimitedFile struct {
	Delimiter rune
	Filepath  string
}

func parseHeader(delimiter rune, header string) ([]string, error) {
	r := csv.NewReader(strings.NewReader(header))

	record, err := r.Read()

	if err != nil {
		return nil, err
	}

	return record, nil
}

// NewDelimitedFile func
func NewDelimitedFile(path string, delimiter rune) (*DelimitedFile, error) {
	// If no delimiter is passed, derive it from the file extension
	if delimiter == '\x00' { // Rune zero value
		ext := filepath.Ext(path)
		switch ext {
		case ".csv":
			delimiter = ','
		case ".tsv":
			delimiter = '\t'
		case ".tab":
			delimiter = '\t'
		case ".pipe":
			delimiter = '|'
		case "":
			return nil, fmt.Errorf("no extension to derive delimiter from")
		default:
			return nil, fmt.Errorf("could not derive delimiter from '%s' extension", ext)
		}
	}
	f := &DelimitedFile{Delimiter: delimiter, Filepath: path}
	return f, nil
}
