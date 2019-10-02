package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// DelimitedFile struct
type DelimitedFile struct {
	Delimiter rune
	Filepath  string
	Columns   []string
	reader    *csv.Reader
}

// DelimitedFileConfig struct
type DelimitedFileConfig struct {
	Delimiter rune
	Filepath  string
	Header    string
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

	// If header is passed, wrap it with an io.Reader and squirrel it away for MultiReader
	var readers []io.Reader
	if conf.Header != "" {
		readers = append(readers, strings.NewReader(conf.Header+"\n"))
	}

	if conf.Filepath != "" {
		file, err := os.Open(conf.Filepath)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		readers = append(readers, file)
	}

	// Create a csv.Reader and set it up
	r := csv.NewReader(io.MultiReader(readers...))
	r.Comma = conf.Delimiter

	columns, err := r.Read()
	if err != nil {
		return nil, err
	}

	f := &DelimitedFile{Delimiter: conf.Delimiter, Filepath: conf.Filepath, Columns: columns, reader: r}
	return f, nil
}
