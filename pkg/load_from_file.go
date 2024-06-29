package iniparser

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

// LoadFromFile opens designated file, read and parse its data
// then store the parsed data in IniParser parsedData field.
//
// A successful load would assign p.parsedData == data and err == nil.
//
// An unsuccessful load would return an error and leave p.parsedData as it is.
func (p *IniParser) LoadFromFile(filePath string) error {
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
