// Package iniparser implements utility methods and functions to parse
// ini files and extract information from them.
package iniparser

import (
	"bufio"
	"errors"
	"fmt"
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
	ErrDataNotMatching         = errors.New("retrieved data is not matching test data")               // test data do not match retrieved data
	ErrParsedDataNotMatching   = errors.New("parsed data is not matching test data")                  // test parsed config data do not match retrieved config data
	ErrParsedStringNotMatching = errors.New("parsed string is not matching test string")              // test parsed config data do not match retrieved config data
	ErrParsedDataMatching      = errors.New("expected error, but got parsed data matching test data") // test parsed config data matching retrieved config data when data is invalid

	ErrInvalidSectionName = errors.New("invalid section name")      // input section name is invalid
	ErrKeyNotFound        = errors.New("key not found")             // input key not found in the section
	ErrSectionNotFound    = errors.New("section not found")         // input section not found in the file
	ErrSectionIsEmpty     = errors.New("section given is empty")    // input section is empty
	ErrKeyIsEmpty         = errors.New("key is empty")              // input key is empty
	ErrValueIsEmpty       = errors.New("value is empty")            // input value is empty
	ErrEmptyString        = errors.New("empty string")              // input is empty string
	ErrParsedDataEmpty    = errors.New("no parsed data to return")  // no parsed data to return
	ErrWritingToFile      = errors.New("error writing to the file") // failed to write to file
	ErrCommentOnNewLine   = errors.New("comment on new line")       // comment on new line
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
		return nil, fmt.Errorf("error: %w\n, given file path: %q", ErrOpeningFile, filePath)
	}

	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	parsedData, err := parseLines(fileLines)

	if err != nil {
		return nil, err
	}

	return parsedData, nil
}

// LoadFromString takes in a string data, parses it
// then store the parsed data in Parser parsedData field.
func (p *Parser) LoadFromString(data string) error {
	if data == "" {
		return ErrEmptyString
	}

	lines := strings.Split(data, "\n")

	parsedData, err := parseLines(lines)

	if err != nil {
		return err
	}

	p.parsedData = parsedData
	return nil
}

// Get retrieves the value of a key in a section.
func (p *Parser) Get(section_name, key string) (string, error) {
	property, sectionExists := p.parsedData[section_name]

	if !sectionExists {
		return "", fmt.Errorf("error: %w\n section: %q does not exist", ErrSectionNotFound, section_name)
	}

	val, keyExists := property[key]
	if !keyExists {
		return "", fmt.Errorf("error: %w\n key: %q does not exist in section: %q", ErrKeyNotFound, key, section_name)
	}

	return val, nil
}

// Set sets the value of a key in a section.
func (p *Parser) Set(section_name, key, value string) error {
	if section_name == "" {
		return ErrSectionIsEmpty
	}

	if key == "" {
		return fmt.Errorf("error: %w\n given section: %q", ErrKeyIsEmpty, section_name)
	}

	_, sectionExists := p.parsedData[section_name]

	if !sectionExists {
		p.parsedData[section_name] = make(map[string]string)
	}

	p.parsedData[section_name][key] = value

	return nil
}

// GetSectionNames returns a slice of section names.
func (p *Parser) GetSectionNames() []string {
	sectionNames := make([]string, 0, len(p.parsedData))

	for sectionName := range p.parsedData {
		sectionNames = append(sectionNames, sectionName)
	}

	return sectionNames
}

// GetSections returns a map of sections and their keys and values.
func (p *Parser) GetSections() (map[string]map[string]string, error) {

	if len(p.parsedData) == 0 {
		return nil, ErrParsedDataEmpty
	}
	return p.parsedData, nil
}

// String returns a string representation of the parsed data.
func (p *Parser) String() string {

	var str string
	for section, properties := range p.parsedData {
		str += "\n[" + section + "]\n"
		for key, value := range properties {
			str += key + "=" + value + "\n"
		}
	}

	return str
}

// SaveToFile saves the parsed data to a file.
func (p *Parser) SaveToFile(filePath string) error {
	if len(p.parsedData) == 0 {
		return ErrParsedDataEmpty
	}

	file, err := os.Create(filePath)
	if err != nil {
		return ErrOpeningFile
	}

	defer file.Close()

	for section, properties := range p.parsedData {
		_, err := file.WriteString("\n[" + section + "]\n")
		if err != nil {
			return fmt.Errorf("error: %w\n, given file path: %q", ErrWritingToFile, filePath)
		}
		for key, value := range properties {
			_, err := file.WriteString(key + "=" + value + "\n")

			if err != nil {
				return fmt.Errorf("error: %w\n, given file path: %q", ErrWritingToFile, filePath)
			}
		}
	}

	return nil
}

func parseLines(lines []string) (map[string]map[string]string, error) {
	parsedData := make(map[string]map[string]string)

	re := regexp.MustCompile(`\[.*?\]`)

	inSection := false

	for i := 0; i < len(lines); i++ {

		if !inSection && (strings.HasPrefix(lines[i], ";") || strings.HasPrefix(lines[i], "#")) {
			continue
		}

		sectionName := re.FindString(lines[i])

		if sectionName != "" {
			inSection = true
			sectionName = sectionName[1 : len(sectionName)-1]

			if sectionName == "" {
				return nil, ErrSectionIsEmpty
			}

			i++
			for i < len(lines) && lines[i] == "" {
				i++
			}

			if strings.HasPrefix(lines[i], ";") || strings.HasPrefix(lines[i], "#") {
				return nil, ErrCommentOnNewLine
			}

			parsedData[sectionName] = make(map[string]string)
			for ; i < len(lines) && lines[i] != ""; i++ {
				keyValuePair := strings.Split(lines[i], "=")
				if keyValuePair[0] == "" {
					return nil, fmt.Errorf("error: %w\n key for section: %q is empty", ErrKeyIsEmpty, sectionName)
				}
				if keyValuePair[1] == "" {
					return nil, fmt.Errorf("error: %w\n value of key: %q is empty", ErrValueIsEmpty, keyValuePair[0])
				}
				key := keyValuePair[0]
				value := keyValuePair[1]
				parsedData[sectionName][strings.TrimSpace(key)] = strings.TrimSpace(value)
			}
		} else {
			inSection = false
		}
	}

	return parsedData, nil
}
