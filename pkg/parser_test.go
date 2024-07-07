package iniparser

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"testing"
)

var validParsedContent = map[string]map[string]string{
	"owner":    {"name": "John Doe", "organization": "Acme Widgets Inc."},
	"database": {"server": "192.0.2.62", "port": "143", "file": "payroll.dat"},
}

const validInput = `;owner section
[owner]
name=John Doe
organization=Acme Widgets Inc.

;database section
[database]
server=192.0.2.62
port=143
file=payroll.dat`

const invalidEmptySectionNameInput = `[owner]
name=John Doe
organization=Acme Widgets Inc.

[]
server=192.0.2.62
port=143
file=payroll.dat`

const invalidEmptyKeyNameInput = `[owner]
name=John Doe
=Acme Widgets Inc.

[database]
server=192.0.2.62
port=143
file=payroll.dat`

const invalidEmptyValueInput = `[owner]
name=John Doe
organization=Acme Widgets Inc.

[database]
server=
port=143
file=payroll.dat`

const inValidCommentOnNewLineInput = `[owner]
name=John Doe
organization=Acme Widgets Inc.

[database]
#server=192.0.2.62
port=143
file=payroll.dat`

func TestParser_LoadFromFile(t *testing.T) {
	parser := NewParser()

	t.Run("load from file with valid data", func(t *testing.T) {
		filePath := "./testdata/valid_data.ini"

		err := parser.LoadFromFile(filePath)

		if err != nil {
			t.Errorf(err.Error())
		}

		assertAreEqual(t, parser.parsedData, validParsedContent)
	})

	var loadInvalidLoadTests = []struct {
		testName string
		filePath string
		want     error
	}{
		{"load from file with empty section name", "./testdata/empty_section_name.ini", ErrSectionIsEmpty},
		{"load from file with empty key name", "./testdata/empty_key_name.ini", ErrKeyIsEmpty},
		{"load from file with empty value", "./testdata/empty_value.ini", ErrValueIsEmpty},
		{"load from file with comment on new line", "./testdata/comment_on_new_line.ini", ErrCommentOnNewLine},
	}

	for _, tt := range loadInvalidLoadTests {
		t.Run(tt.testName, func(t *testing.T) {
			err := parser.LoadFromFile(tt.filePath)

			if err == nil {
				t.Error("Expected error, but got nil")
			}

			if err != nil {
				switch err {
				case ErrOpeningFile:
					t.Error(err)
				default:
					assertError(t, err, tt.want)
				}
			}
		})
	}
}

func ExampleParser_LoadFromFile() {
	parser := NewParser()
	_ = parser.LoadFromFile("./testdata/valid_data.ini")
	val, _ := parser.Get("owner", "name")
	fmt.Println(val)
	// Output: John Doe
}

func TestParser_LoadFromString(t *testing.T) {
	parser := NewParser()

	var loadTests = []struct {
		testName string
		input    string
		want     map[string]map[string]string
		err      error
	}{
		{"load from string with valid data", validInput, validParsedContent, nil},
		{"load from string with empty section name", invalidEmptySectionNameInput, nil, ErrSectionIsEmpty},
		{"load from string with empty key name", invalidEmptyKeyNameInput, nil, ErrKeyIsEmpty},
		{"load from string with empty value", invalidEmptyValueInput, nil, ErrValueIsEmpty},
		{"load from string with comment on new line", inValidCommentOnNewLineInput, nil, ErrCommentOnNewLine},
	}

	for _, tt := range loadTests {
		t.Run(tt.testName, func(t *testing.T) {
			err := parser.LoadFromString(tt.input)

			if err == nil {
				assertAreEqual(t, parser.parsedData, tt.want)
			}
			assertError(t, err, tt.err)
		})
	}
}

func ExampleParser_LoadFromString() {
	parser := NewParser()

	_ = parser.LoadFromString(validInput)
}

func TestParser_GetSectionNames(t *testing.T) {
	parser := NewParser()

	parser.parsedData = validParsedContent

	sectionNames := parser.GetSectionNames()

	validSectionNames := []string{"owner", "database"}

	if !(len(sectionNames) == len(validSectionNames) && (slices.Contains(sectionNames, "owner") && slices.Contains(sectionNames, "database"))) {
		t.Errorf("got %q want %q", sectionNames, validSectionNames)
	}
}

func ExampleParser_GetSectionNames() {
	parser := NewParser()
	_ = parser.LoadFromFile("./testdata/valid_data.ini")
	fmt.Println(parser.GetSectionNames())
}

func TestParser_String(t *testing.T) {
	parser := NewParser()

	t.Run("return string with valid data", func(t *testing.T) {
		parser.parsedData = validParsedContent

		str := parser.String()
		_ = parser.LoadFromString(str)

		assertAreEqual(t, parser.parsedData, validParsedContent)
	})
}

func ExampleParser_String() {
	parser := NewParser()
	_ = parser.LoadFromFile("./testdata/valid_data.ini")
	str := parser.String()
	fmt.Println(str)
}

func TestParser_Get(t *testing.T) {
	parser := NewParser()

	parser.parsedData = validParsedContent

	var getTests = []struct {
		testName string
		section  string
		key      string
		want     string
		err      error
	}{
		{"get value from existing section and key", "owner", "name", "John Doe", nil},
		{"get value from existing section and non-existing key", "owner", "config", "", ErrKeyNotFound},
		{"get value from non-existing section", "config", "", "", ErrSectionNotFound},
	}

	for _, tt := range getTests {
		t.Run(tt.testName, func(t *testing.T) {
			value, err := parser.Get(tt.section, tt.key)

			if err != nil {
				assertError(t, err, tt.err)
			}

			if value != tt.want {
				t.Errorf("got %q want %q", value, tt.want)
			}

		})

	}
}

func ExampleParser_Get() {
	parser := NewParser()
	_ = parser.LoadFromFile("./testdata/valid_data.ini")
	val, _ := parser.Get("owner", "name")
	fmt.Println(val)
	// Output: John Doe
}

func TestParser_Set(t *testing.T) {

	parser := NewParser()

	parser.parsedData = validParsedContent

	var setTests = []struct {
		testName string
		section  string
		key      string
		value    string
		want     string
		err      error
	}{
		{"set value for existing section and key", "owner", "name", "Eyad", "Eyad", nil},
		{"set value for existing section and non-existing key", "owner", "config", "data", "data", ErrKeyNotFound},
		{"set value for non-existing section", "config", "database", "192.178.292.1", "192.178.292.1", ErrSectionNotFound},
		{"set value on empty section", "", "", "", "", ErrSectionIsEmpty},
		{"set value on empty key", "owner", "", "", "", ErrKeyIsEmpty},
	}

	for _, tt := range setTests {
		t.Run(tt.testName, func(t *testing.T) {
			err := parser.Set(tt.section, tt.key, tt.value)
			if err != nil {
				assertError(t, err, tt.err)
			}
			value, _ := parser.Get(tt.section, tt.key)
			if value != tt.want {
				t.Errorf("got %q want %q", value, tt.want)
			}
		})
	}
}

func ExampleParser_Set() {
	parser := NewParser()
	_ = parser.LoadFromFile("./testdata/valid_data.ini")
	_ = parser.Set("owner", "name", "Eyad")
	val, _ := parser.Get("owner", "name")
	fmt.Println(val)
	// Output: Eyad
}

func TestParser_GetSections(t *testing.T) {

	t.Run("return sections with valid data", func(t *testing.T) {
		parser := NewParser()
		parser.parsedData = validParsedContent

		sections := parser.GetSections()

		assertAreEqual(t, sections, validParsedContent)
	})

	t.Run("return sections with empty data", func(t *testing.T) {
		parser := NewParser()

		sections := parser.GetSections()

		assertAreEqual(t, sections, make(map[string]map[string]string))
	})
}

func ExampleParser_GetSections() {
	parser := NewParser()
	_ = parser.LoadFromFile("./testdata/valid_data.ini")
	sections := parser.GetSections()

	fmt.Println(sections)
}

func TestParser_SaveToFile(t *testing.T) {
	parser := NewParser()
	filePath := "./testdata/output.ini"

	t.Run("save to file with valid data", func(t *testing.T) {
		parser.parsedData = validParsedContent

		err := parser.SaveToFile(filePath)

		if err != nil {
			t.Error(err)
		}

		err = parser.LoadFromFile(filePath)

		if err != nil {
			t.Error(err)
		}

		assertAreEqual(t, parser.parsedData, validParsedContent)
	})
}

func ExampleParser_SaveToFile() {
	parser := NewParser()
	_ = parser.LoadFromFile("./testdata/valid_data.ini")
	_ = parser.Set("owner", "name", "Eyad")
	_ = parser.SaveToFile("./testdata/output.ini")
}

func assertError(t testing.TB, got, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertAreEqual(t testing.TB, got, want any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q want %q", got, want)
	}
}
