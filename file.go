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

// DelimitedFileConfig struct
type DelimitedFileConfig struct {
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

func deriveDelimiter(path string) (rune, error) {
	var delimiter rune
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
		return 0, fmt.Errorf("no extension to derive delimiter from")
	default:
		return 0, fmt.Errorf("could not derive delimiter from '%s' extension", ext)
	}

	return delimiter, nil
}

// NewDelimitedFile func
func NewDelimitedFile(conf DelimitedFileConfig) (*DelimitedFile, error) {
	// If no delimiter is passed, attempt to derive one from the file extension
	if conf.Delimiter == 0 { // Rune zero value
		var err error
		conf.Delimiter, err = deriveDelimiter(conf.Filepath)
		if err != nil {
			return nil, err
		}
	}
	f := &DelimitedFile{Delimiter: conf.Delimiter, Filepath: conf.Filepath}
	return f, nil
}
