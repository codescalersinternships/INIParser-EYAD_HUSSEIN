package iniparser

import (
	"fmt"
	"reflect"
	"testing"
)

var validParsedContent = map[string]map[string]string{
	"owner":    {"name": "John Doe", "organization": "Acme Widgets Inc."},
	"database": {"server": "192.0.2.62", "port": "143", "file": "payroll.dat"},
}

func TestParser_LoadFromFile(t *testing.T) {
	t.Run("load from file with valid data", func(t *testing.T) {
		filePath := "./test_file_1.ini"

		parser := NewParser()

		err := parser.LoadFromFile(filePath)

		if err != nil {
			switch err {
			case ErrReadingFile:
			case ErrOpeningFile:
				t.Errorf("%q\n, file name given:\n%q", ErrReadingFile.Error(), filePath)
			default:
				t.Errorf(err.Error())
			}
		}

		assertTwoMaps(t, parser.parsedData, validParsedContent)
	})

	t.Run("load from file with invalid data", func(t *testing.T) {
		filePath := "./test_file_2.ini"

		parser := NewParser()

		err := parser.LoadFromFile(filePath)

		if err != nil {
			switch err {
			case ErrReadingFile:
			case ErrOpeningFile:
				t.Errorf("%q\n, file name given:\n%q", ErrReadingFile.Error(), filePath)
			default:
				t.Errorf(err.Error())
			}
		}

		if reflect.DeepEqual(parser.parsedData, validParsedContent) {
			t.Errorf(ErrParsedDataMatching.Error())
		}
	})
}

func assertTwoMaps(t testing.TB, got, want map[string]map[string]string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got:\n\t%q \nwant:\n\t%q", got, want)
	}
}

func assertStrings(t testing.TB, got, want string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got:\n\t%q \nwant:\n\t%q", got, want)
	}
}

func ExampleParser_LoadFromFile() {
	parser := NewParser()
	_ = parser.LoadFromFile("./test_file_1.ini")
	val, _ := parser.Get("owner", "name")
	fmt.Println(val)
	// Output: John Doe
}

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

func TestParser_LoadFromString(t *testing.T) {
	t.Run("load from string with valid string input", func(t *testing.T) {
		parser := NewParser()

		parser.LoadFromString(validStringInput)

		assertTwoMaps(t, parser.parsedData, validParsedContent)
	})

	t.Run("load from string with invalid string input", func(t *testing.T) {
		parser := NewParser()

		parser.LoadFromString(inValidStringInput)

		if reflect.DeepEqual(parser.parsedData, validParsedContent) {
			t.Errorf(ErrParsedDataMatching.Error())
		}
	})
}

func ExampleParser_LoadFromString() {
	parser := NewParser()

	parser.LoadFromString(validStringInput)
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

			assertStrings(t, value, tt.want)

		})

	}
}

func assertError(t testing.TB, got, want error) {
	t.Helper()
	if got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}
}

func ExampleParser_Get() {
	parser := NewParser()
	_ = parser.LoadFromFile("./test_file_1.ini")
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
		{"set value on empty section", "", "", "", "", ErrSectionIsEmtpy},
		{"set value on empty key", "owner", "", "", "", ErrKeyIsEmtpy},
	}

	for _, tt := range setTests {
		t.Run(tt.testName, func(t *testing.T) {
			err := parser.Set(tt.section, tt.key, tt.value)
			if err != nil {
				assertError(t, err, tt.err)
			}
			value, _ := parser.Get(tt.section, tt.key)
			assertStrings(t, value, tt.want)
		})
	}
}

func ExampleParser_Set() {
	parser := NewParser()
	_ = parser.LoadFromFile("./test_file_1.ini")
	_ = parser.Set("owner", "name", "Eyad")
	val, _ := parser.Get("owner", "name")
	fmt.Println(val)
	// Output: Eyad
}
