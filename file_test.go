package main

/*
import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"io"
	"io/ioutil"
	"os"
	"syscall"
	"testing"
)

func TestDeriveDelimiter(t *testing.T) {
	tests := map[string]struct {
		filepath string
		want     rune
		err      error
	}{
		"csv":  {"foo.csv", ',', nil},
		"tsv":  {"foo.tsv", '\t', nil},
		"tab":  {"foo.tab", '\t', nil},
		"pipe": {"foo.pipe", '|', nil},

		"no ext":      {"foo", 0, errors.New("no extension to derive delimiter from")},
		"unknown ext": {"foo.foo", 0, errors.New("could not derive delimiter from '.foo' extension")},
	}

	equateErrorMessage := cmp.Comparer(func(x, y error) bool {
		if x == nil || y == nil {
			return x == nil && y == nil
		}
		return x.Error() == y.Error()
	})

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			d, err := deriveDelimiter(tc.filepath)

			if diff := cmp.Diff(tc.want, d); diff != "" {
				t.Errorf("deriveDelimiter(%q) mismatch (-want +got):\n%s", tc.filepath, diff)
			}

			if diff := cmp.Diff(tc.err, err, equateErrorMessage); diff != "" {
				t.Errorf("Error mismatch for deriveDelimiter(%q) (-want +got):\n%s", tc.filepath, diff)
			}
		})
	}
}

func TestNewDelimitedFile(t *testing.T) {
	var err error

	var nonExistentCSVFile *os.File
	nonExistentCSVFile, err = ioutil.TempFile("", "*.empty.csv")
	if err != nil {
		t.Errorf("failed to create non-existent csv temp file: %s", err)
	}
	nonExistentCSVFile.Close()
	os.Remove(nonExistentCSVFile.Name()) // clean up

	var emptyCSVFile *os.File
	emptyCSVFile, err = ioutil.TempFile("", "*.empty.csv")
	if err != nil {
		t.Errorf("failed to create empty csv temp file: %s", err)
	}
	emptyCSVFile.Close()
	defer os.Remove(emptyCSVFile.Name()) // clean up

	var headerCSVFile *os.File
	headerCSVFile, err = ioutil.TempFile("", "*.header.csv")
	if _, err := headerCSVFile.Write([]byte("fcol1,fcol2,fcol3")); err != nil {
		headerCSVFile.Close()
	}
	headerCSVFile.Close()
	defer os.Remove(headerCSVFile.Name()) // clean up

	tests := map[string]struct {
		conf DelimitedFileConfig
		want *DelimitedFile
		err  error
	}{
		"all empty":              {DelimitedFileConfig{}, nil, errors.New("no extension to derive delimiter from")},
		"non-existent filepath":  {DelimitedFileConfig{Filepath: nonExistentCSVFile.Name()}, nil, &os.PathError{Op: "open", Path: nonExistentCSVFile.Name(), Err: syscall.Errno(0x02)}},
		"derivable filepath":     {DelimitedFileConfig{Filepath: emptyCSVFile.Name()}, nil, errors.New("EOF")},
		"underivable filepath":   {DelimitedFileConfig{Filepath: "*.foo.foo"}, nil, errors.New("could not derive delimiter from '.foo' extension")},
		"arg header only":        {DelimitedFileConfig{Header: "col1,col2,col3"}, nil, errors.New("no extension to derive delimiter from")},
		"delim only":             {DelimitedFileConfig{Delimiter: ','}, nil, errors.New("EOF")},
		"header+delim":           {DelimitedFileConfig{Delimiter: ',', Header: "col1,col2,col3"}, &DelimitedFile{Delimiter: ',', Columns: []string{"col1", "col2", "col3"}}, nil},
		"file header":            {DelimitedFileConfig{Filepath: headerCSVFile.Name(), Delimiter: ','}, &DelimitedFile{Filepath: headerCSVFile.Name(), Delimiter: ',', Columns: []string{"fcol1", "fcol2", "fcol3"}}, nil},
		"file header+arg header": {DelimitedFileConfig{Filepath: headerCSVFile.Name(), Delimiter: ',', Header: "hcol1,hcol2,hcol3"}, &DelimitedFile{Filepath: headerCSVFile.Name(), Delimiter: ',', Columns: []string{"hcol1", "hcol2", "hcol3"}}, nil},
	}

	equateErrorMessage := cmp.Comparer(func(x, y error) bool {
		if x == nil || y == nil {
			return x == nil && y == nil
		}
		return x.Error() == y.Error()
	})

	equateDelimitedFile := cmpopts.IgnoreUnexported(DelimitedFile{})

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			f, err := NewDelimitedFile(tc.conf)

			if diff := cmp.Diff(tc.want, f, equateDelimitedFile); diff != "" {
				t.Errorf("NewDelimitedFile(%q) mismatch (-want +got):\n%s", tc.conf, diff)
			}

			if diff := cmp.Diff(tc.err, err, equateErrorMessage); diff != "" {
				t.Errorf("Error mismatch for NewDelimitedFile(%q) (-want +got):\n%s", tc.conf, diff)
			}
		})
	}
}

func TestDelimitedFileRead(t *testing.T) {
	var err error

	var emptyCSVFile *os.File
	emptyCSVFile, err = ioutil.TempFile("", "*.empty.csv")
	if err != nil {
		t.Errorf("failed to create empty csv temp file: %s", err)
	}
	emptyCSVFile.Close()
	defer os.Remove(emptyCSVFile.Name()) // clean up

	var headerCSVFile *os.File
	headerCSVFile, err = ioutil.TempFile("", "*.header.csv")
	if _, err := headerCSVFile.Write([]byte("fcol1,fcol2,fcol3")); err != nil {
		headerCSVFile.Close()
	}
	headerCSVFile.Close()
	defer os.Remove(headerCSVFile.Name()) // clean up

	var contentCSVFile *os.File
	contentCSVFile, err = ioutil.TempFile("", "*.content.csv")
	content := `fcol1,fcol2,fcol3
row1col1,row1col2,row1col3
row2col1,row2col2,row2col3
row3col1,row3col2,row3col3
`
	if _, err := contentCSVFile.Write([]byte(content)); err != nil {
		contentCSVFile.Close()
	}
	contentCSVFile.Close()
	defer os.Remove(contentCSVFile.Name()) // clean up

	tests := map[string]struct {
		conf DelimitedFileConfig
		want []*DelimitedRecord
		err  error
	}{
		//"derivable filepath":     {DelimitedFileConfig{Filepath: emptyCSVFile.Name()}, nil, nil},
		"file header":            {DelimitedFileConfig{Filepath: headerCSVFile.Name(), Delimiter: ','}, nil, io.EOF},
		"file header+arg header": {DelimitedFileConfig{Filepath: headerCSVFile.Name(), Delimiter: ',', Header: "hcol1,hcol2,hcol3"}, nil, errors.New("file already closed")},
		"file header+content": {DelimitedFileConfig{Filepath: contentCSVFile.Name(), Delimiter: ','},
			[]*DelimitedRecord{
				{RecordNumber: 1, Record: []string{"row1col1", "row1col2", "row1col3"}},
				{RecordNumber: 2, Record: []string{"row2col1", "row2col2", "row2col3"}},
				{RecordNumber: 3, Record: []string{"row3col1", "row3col2", "row3col3"}},
			},
			nil},
		"file header+arg header+content": {DelimitedFileConfig{Filepath: contentCSVFile.Name(), Delimiter: ',', Header: "hcol1,hcol2,hcol3"},
			[]*DelimitedRecord{
				{RecordNumber: 1, Record: []string{"row1col1", "row1col2", "row1col3"}},
				{RecordNumber: 2, Record: []string{"row2col1", "row2col2", "row2col3"}},
				{RecordNumber: 3, Record: []string{"row3col1", "row3col2", "row3col3"}},
			},
			errors.New("file already closed")},
	}

		equateErrorMessage := cmp.Comparer(func(x, y error) bool {
			if x == nil || y == nil {
				return x == nil && y == nil
			}
			return x.Error() == y.Error()
		})

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			f, err := NewDelimitedFile(tc.conf)
			var allRecs []*DelimitedRecord
			for {
				var rec *DelimitedRecord
				rec, err = f.Read()
				if err == io.EOF {
					//		t.Errorf("Error reading record: %s", err)
					break
				}
				if err != nil {
					//t.Errorf("Error reading record: %s", err.Error())
					break
				}
				allRecs = append(allRecs, rec)
			}

			if diff := cmp.Diff(tc.want, allRecs); diff != "" {
				t.Errorf("NewDelimitedFile(%q) mismatch (-want +got):\n%s", tc.conf, diff)
			}

			//			if diff := cmp.Diff(tc.err, err, equateErrorMessage); diff != "" {
			//				t.Errorf("Error mismatch for NewDelimitedFile(%q) (-want +got):\n%s", tc.conf, diff)
			//			}
		})
	}
}
*/
