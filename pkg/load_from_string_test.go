package iniparser

import (
	"testing"
)

const validStringInput = `
[owner]
name=John Doe
organization=Acme Widgets Inc.

[database]
server=192.0.2.62
port=143
file=payroll.dat`

const inValidStringInput = `
[owner]
name=Eyad
organization=Acme Widgets Inc.

[database]
url=192.0.2.62
port=143`

func TestLoadFromString(t *testing.T) {
	t.Run("load from string with valid string input", func(t *testing.T) {
		parser := NewParser()

		parser.LoadFromString(validStringInput)

		assertValidParsedData(t, parser.parsedData, validParsedContent)
	})

	t.Run("load from string with invalid string input", func(t *testing.T) {
		parser := NewParser()

		parser.LoadFromString(inValidStringInput)

		assertInvalidParsedData(t, parser.parsedData, validParsedContent)
	})
}

func ExampleIniParser_LoadFromString() {
	parser := NewParser()

	parser.LoadFromString(validStringInput)
}
