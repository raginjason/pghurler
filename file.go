package main

import (
	"fmt"
	"path/filepath"
)

// DelimitedFile struct
type DelimitedFile struct {
	Delimiter rune
	Filepath  string
}

// NewDelimitedFile func
func NewDelimitedFile(path string, delimiter rune, columns []string) (*DelimitedFile, error) {
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
