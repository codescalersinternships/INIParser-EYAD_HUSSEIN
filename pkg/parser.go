// Package iniparser implements utility methods and functions to parse
// ini files and extract information from them.
package iniparser

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// A Parser loads and manipulates ini files as requested.
// The zero value for Parser is a parser ready to use.
type Parser struct {
	parsedData map[string]map[string]string
}

// NewParser returns a new Parser.
func NewParser() *Parser {
	return &Parser{parsedData: make(map[string]map[string]string)}
}

var (
	ErrOpeningFile      = errors.New("error opening the file")    // failed to open file
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
func (p *Parser) GetSections() map[string]map[string]string {
	return p.parsedData
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
	err := os.WriteFile(filePath, []byte(p.String()), 0644)

	if err != nil {
		return ErrWritingToFile
	}
	return nil
}

func parseLines(lines []string) (map[string]map[string]string, error) {
	parsedData := make(map[string]map[string]string)
	var currentSection string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			sectionName := strings.TrimSpace(line[1 : len(line)-1])
			if sectionName == "" {
				return nil, ErrSectionIsEmpty
			}
			currentSection = sectionName
			if _, exists := parsedData[currentSection]; !exists {
				parsedData[currentSection] = make(map[string]string)
			}
			continue
		}

		if strings.Contains(line, "=") {
			if strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
				return nil, ErrCommentOnNewLine
			}
			keyValuePair := strings.Split(line, "=")
			key := strings.TrimSpace(keyValuePair[0])
			value := strings.TrimSpace(keyValuePair[1])
			if key == "" {
				return nil, ErrKeyIsEmpty
			}
			if value == "" {
				return nil, ErrValueIsEmpty
			}
			parsedData[currentSection][key] = value
		}
	}

	return parsedData, nil
}
