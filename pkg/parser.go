// Package iniparser implements utility methods and functions to parse
// ini files and extract information from them.
package iniparser

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strings"
)

// A Parser loads and manipulates ini files as requested.
// The zero value for Parser is a parser ready to use.
type Parser struct {
	parsedData map[string]map[string]string
}

// NewParser returns a new Parser.
func NewParser() *Parser {
	return &Parser{}
}

var (
	ErrOpeningFile             = errors.New("error opening the file")                                 // failed to open file
	ErrReadingFile             = errors.New("error reading the file")                                 // failed during reading file
	ErrDataNotMatching         = errors.New("retrieved data is not matching test data")               // test data do not match retrieved data
	ErrParsedDataNotMatching   = errors.New("parsed data is not matching test data")                  // test parsed config data do not match retrieved config data
	ErrParsedStringNotMatching = errors.New("parsed string is not matching test string")              // test parsed config data do not match retrieved config data
	ErrParsedDataMatching      = errors.New("expected error, but got parsed data matching test data") // test parsed config data matching retrieved config data when data is invalid

	ErrKeyNotFound     = errors.New("key not found")     // key not found in the section
	ErrSectionNotFound = errors.New("section not found") // section not found in the file
	ErrSectionIsEmtpy  = errors.New("section is empty")  // section is empty
	ErrKeyIsEmtpy      = errors.New("key is empty")      // key is empty
)

// LoadFromFile opens designated file, read and parse its data
// then store the parsed data in Parser parsedData field.
//
// A successful load would assign p.parsedData == data and err == nil.
//
// An unsuccessful load would return an error and leave p.parsedData as it is.
func (p *Parser) LoadFromFile(filePath string) error {
	data, err := parseFileData(filePath)

	if err != nil {
		return err
	}

	p.parsedData = data
	return nil
}

func parseFileData(filePath string) (map[string]map[string]string, error) {

	readFile, err := os.Open(filePath)

	if err != nil {
		return nil, ErrOpeningFile
	}

	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	parsedData := make(map[string]map[string]string)

	re := regexp.MustCompile(`\[.*?\]`)

	for i := 0; i < len(fileLines); i++ {

		sectionName := re.FindString(fileLines[i])

		if sectionName != "" {
			i++
			for fileLines[i] == "" {
				i++
			}

			sectionName = sectionName[1 : len(sectionName)-1]

			parsedData[sectionName] = make(map[string]string)
			for ; i < len(fileLines) && fileLines[i] != ""; i++ {
				keyValuePair := strings.Split(fileLines[i], "=")
				key := keyValuePair[0]
				value := keyValuePair[1]
				parsedData[sectionName][strings.TrimSpace(key)] = strings.TrimSpace(value)
			}
		}
	}

	return parsedData, nil
}

// LoadFromString takes in a string data, parses it
// then store the parsed data in Parser parsedData field.
func (p *Parser) LoadFromString(data string) {

	lines := strings.Split(data, "\n")

	re := regexp.MustCompile(`\[.*?\]`)
	parsedData := make(map[string]map[string]string)

	for i := 0; i < len(lines); i++ {

		sectionName := re.FindString(lines[i])

		if sectionName != "" {
			i++
			for lines[i] == "" {
				i++
			}

			sectionName = sectionName[1 : len(sectionName)-1]

			parsedData[sectionName] = make(map[string]string)
			for ; i < len(lines) && lines[i] != ""; i++ {
				keyValuePair := strings.Split(lines[i], "=")
				key := keyValuePair[0]
				value := keyValuePair[1]
				parsedData[sectionName][strings.TrimSpace(key)] = strings.TrimSpace(value)
			}
		}
	}

	p.parsedData = parsedData
}

// Get retrieves the value of a key in a section.
func (p *Parser) Get(section_name, key string) (string, error) {
	property, sectionExists := p.parsedData[section_name]

	if !sectionExists {
		return "", ErrSectionNotFound
	}

	val, keyExists := property[key]
	if !keyExists {
		return "", ErrKeyNotFound
	}

	return val, nil
}

// Set sets the value of a key in a section.
func (p *Parser) Set(section_name, key, value string) error {
	if section_name == "" {
		return ErrSectionIsEmtpy
	}

	if key == "" {
		return ErrKeyIsEmtpy
	}

	_, sectionExists := p.parsedData[section_name]

	if !sectionExists {
		p.parsedData[section_name] = make(map[string]string)
	}

	p.parsedData[section_name][key] = value

	return nil
}
