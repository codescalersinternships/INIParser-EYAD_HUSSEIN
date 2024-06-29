package iniparser

import (
	"reflect"
	"testing"
)

var validParsedContent = map[string]map[string]string{
	"owner":    {"name": "John Doe", "organization": "Acme Widgets Inc."},
	"database": {"server": "192.0.2.62", "port": "143", "file": "payroll.dat"},
}

func TestLoadFromFile(t *testing.T) {
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

		assertValidParsedData(t, parser.parsedData, validParsedContent)
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

		assertInvalidParsedData(t, parser.parsedData, validParsedContent)
	})
}

func assertValidParsedData(t testing.TB, got, want map[string]map[string]string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%q:\n got:\n\t%q \nwant:\n\t%q", ErrParsedDataNotMatching.Error(), got, want)
	}
}

func assertInvalidParsedData(t testing.TB, got, want map[string]map[string]string) {
	t.Helper()
	if reflect.DeepEqual(got, want) {
		t.Errorf(ErrParsedDataMatching.Error())
	}
}

func ExampleIniParser_LoadFromFile() {
	parser := IniParser{}
	_ = parser.LoadFromFile("./data.ini")
	// Output:
}
