package iniparser

import (
	"reflect"
	"testing"
)

func TestLoadFromFile(t *testing.T) {

	filePath := "./t.ini"

	parsedContent := map[string]map[string]string{
		"forge": {"User": "hg", "Logger": "dksjdj"},
		"tail":  {"Hi": "no", "yes": "hola"},
	}

	parser := IniParser{}

	err := parser.LoadFromFile(filePath)

	if err != nil {
		switch err {
		case ErrReadingFile:
		case ErrOpeningFile:
			t.Errorf("%q, file name given: %q", ErrReadingFile.Error(), filePath)
		default:
			t.Errorf(err.Error())
		}
	}

	assertParsedData(t, parser.parsedData, parsedContent, ErrParsedDataNotMatching)
}

func assertParsedData(t testing.TB, got, want map[string]map[string]string, err error) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%q, got %q want %q", err.Error(), got, want)
	}
}

// func ExampleLoadFromFile() {
// 	parser := IniParser{}
// 	parser.LoadFromFile("./data.ini")
// }
