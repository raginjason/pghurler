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

// NewDelimitedFile func
func NewDelimitedFile(conf DelimitedFileConfig) (*DelimitedFile, error) {
	// If no delimiter is passed, derive it from the file extension
	if conf.Delimiter == '\x00' { // Rune zero value
		ext := filepath.Ext(conf.Filepath)
		switch ext {
		case ".csv":
			conf.Delimiter = ','
		case ".tsv":
			conf.Delimiter = '\t'
		case ".tab":
			conf.Delimiter = '\t'
		case ".pipe":
			conf.Delimiter = '|'
		case "":
			return nil, fmt.Errorf("no extension to derive delimiter from")
		default:
			return nil, fmt.Errorf("could not derive delimiter from '%s' extension", ext)
		}
	}
	f := &DelimitedFile{Delimiter: conf.Delimiter, Filepath: conf.Filepath}
	return f, nil
}
