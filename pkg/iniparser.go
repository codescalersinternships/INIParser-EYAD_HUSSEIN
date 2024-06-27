// Package iniparser implements utility methods and functions to parse
// ini files and extract information from them.
package iniparser

import "errors"

// A IniParser loads and manipulates ini files as requested.
// The zero value for Parser is a parser ready to use.
type IniParser struct {
	filePath   string
	parsedData map[string]map[string]string
}

var (
	ErrOpeningFile           = errors.New("error opening the file")                   // failed to open file
	ErrReadingFile           = errors.New("error reading the file")                   // failed during reading file
	ErrDataNotMatching       = errors.New("retrieved data is not matching test data") // test data do not match retrieved data
	ErrParsedDataNotMatching = errors.New("parsed data is not matching test data")    // test parsed config data do not match retrieved config data
)
