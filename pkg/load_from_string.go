package iniparser

import (
	"regexp"
	"strings"
)

// LoadFromString takes in a string data, parses it
// then store the parsed data in IniParser parsedData field.
func (p *IniParser) LoadFromString(data string) {
	line := ""
	lines := []string{}

	for index, ch := range data {
		if ch == '\n' || index == len(data)-1 {
			if index == len(data)-1 {
				line += string(ch)
			}
			lines = append(lines, line)
			line = ""
		} else {
			line += string(ch)
		}
	}

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
