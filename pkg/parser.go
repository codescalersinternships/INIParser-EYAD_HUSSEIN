// Package iniparser implements utility methods and functions to parse
// ini files and extract information from them.
package iniparser

import (
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

	ErrKeyNotFound      = errors.New("key not found")             // input key not found in the section
	ErrSectionNotFound  = errors.New("section not found")         // input section not found in the file
	ErrSectionIsEmpty   = errors.New("section given is empty")    // input section is empty
	ErrKeyIsEmpty       = errors.New("key is empty")              // input key is empty
	ErrValueIsEmpty     = errors.New("value is empty")            // input value is empty
	ErrEmptyString      = errors.New("empty string")              // input is empty string
	ErrParsedDataEmpty  = errors.New("no parsed data to return")  // no parsed data to return
	ErrWritingToFile    = errors.New("error writing to the file") // failed to write to file
	ErrCommentOnNewLine = errors.New("comment on new line")       // comment on new line
)

// LoadFromFile opens designated file, read and parse its data
// then store the parsed data in Parser parsedData field.
//
// A successful load would assign p.parsedData == data and err == nil.
//
// An unsuccessful load would return an error and leave p.parsedData as it is.
func (p *Parser) LoadFromFile(filePath string) error {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("%w given file path %q", ErrOpeningFile, filePath)
	}
	return p.LoadFromString(string(fileData))
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
func (p *Parser) Get(sectionName, key string) (string, error) {
	property, sectionExists := p.parsedData[sectionName]

	if !sectionExists {
		return "", fmt.Errorf("%w section %q does not exist", ErrSectionNotFound, sectionName)
	}

	val, keyExists := property[key]
	if !keyExists {
		return "", fmt.Errorf("%w key %q does not exist in section %q", ErrKeyNotFound, key, sectionName)
	}

	return val, nil
}

// Set sets the value of a key in a section.
func (p *Parser) Set(sectionName, key, value string) error {
	if sectionName == "" {
		return ErrSectionIsEmpty
	}

	if key == "" {
		return fmt.Errorf("%w given section %q", ErrKeyIsEmpty, sectionName)
	}

	_, sectionExists := p.parsedData[sectionName]

	if !sectionExists {
		p.parsedData[sectionName] = make(map[string]string)
	}

	p.parsedData[sectionName][key] = value

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
		str += "[" + section + "]\n"
		for key, value := range properties {
			str += key + "=" + value + "\n"
		}
	}

	str = strings.TrimSuffix(str, "\n")

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

	_, err = file.WriteString(p.String())

	if err != nil {
		return ErrWritingToFile
	}
	return nil
}

func parseLines(lines []string) (map[string]map[string]string, error) {
	parsedData := make(map[string]map[string]string)
	re := regexp.MustCompile(`\[.*?\]`)

	inSection := false
	var sectionName string

	for i := 0; i < len(lines); i++ {
		if !inSection && (strings.HasPrefix(lines[i], ";") || strings.HasPrefix(lines[i], "#")) {
			continue
		}

		sectionNameMatch := re.FindString(lines[i])

		if len(sectionNameMatch) == 2 {
			return nil, ErrSectionIsEmpty
		}

		if sectionNameMatch != "" {
			sectionName = sectionNameMatch[1 : len(sectionNameMatch)-1]
			parsedData[sectionName] = make(map[string]string)
			inSection = true
			continue
		}

		if inSection {
			for i < len(lines) && lines[i] == "" {
				i++
			}

			if i >= len(lines) {
				break
			}

			if strings.HasPrefix(lines[i], ";") || strings.HasPrefix(lines[i], "#") {
				return nil, ErrCommentOnNewLine
			}

			for ; i < len(lines); i++ {
				if lines[i] == "" || re.MatchString(lines[i]) {
					inSection = false
					i--
					break
				}

				keyValuePair := strings.Split(lines[i], "=")
				if len(keyValuePair) != 2 {
					return nil, fmt.Errorf("invalid key-value pair %q", lines[i])
				}

				key := strings.TrimSpace(keyValuePair[0])
				value := strings.TrimSpace(keyValuePair[1])
				if key == "" {
					return nil, fmt.Errorf("%w key for section %q is empty", ErrKeyIsEmpty, sectionName)
				}
				if value == "" {
					return nil, fmt.Errorf("%w value of key %q is empty", ErrValueIsEmpty, key)
				}
				parsedData[sectionName][key] = value
			}
		}
	}

	return parsedData, nil
}
