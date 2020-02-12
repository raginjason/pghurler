/*
Copyright © 2019 Jason Walker <ragin.jason@me.com>
This file is part of pghurler.
*/
package main

import (
	"fmt"
	"path/filepath"
)

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
