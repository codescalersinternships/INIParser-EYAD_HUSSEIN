# INI Parser using Go

## Project Description

This is a library that provides as ini file parser using Go language, it implements several methods to help parse, access, edit, save, and write ini files. Important to note that it works with ini files with the following specifications:

- assume there're no global keys, every keys need to be part of a section
- assume the key value separator is just =
- keys and values should have spaces trimmed
- comments are only valid at the beginning of the line

### Features

- Load and parse string-like ini configuration file
- Load and parse a full ini file
- Easily get all section names
- Serialize the raw content to map-like structure
- Get a value of a key inside any section
- Set a value of a key inside any section
- Convert map-like structure containing ini data to string
- Write out and save parsed data to a file

## How to Use

1- import package

```golang
import github.com/codescalersinternships/INIParser-EYAD_HUSSEIN
```

2- create a new parser struct using NewParser()

```golang
parser := NewParser()
```

3- load from a file

```golang
_ = parser.LoadFromFile("./test-files/test_file_1.ini")
val, _ := parser.Get("owner", "name")
fmt.Println(val)
// Output: John Doe
```

4- load from a string

```golang
_ = parser.LoadFromString(validStringInput)
```

5- get a key value from a section

```golang
val, _ := parser.Get("owner", "name")
fmt.Println(val)
// Output: John Doe
```

6- set a value for a key in a section

```golang
_ = parser.Set("owner", "name", "Eyad")
val, _ := parser.Get("owner", "name")
fmt.Println(val)
// Output: Eyad
```

7- get section names

```golang
sectionsNames, _ := parser.GetSectionNames()
```

8- get parsed data

```golang
sections, _ := parser.GetSections()
```

9- convert data to string

```golang
str, _ := parser.ToString()
```

10- save data to file

```golang
_ = parser.LoadFromFile("./test-files/test_file_1.ini")
_ = parser.Set("owner", "name", "Eyad")
_ = parser.SaveToFile("./test-files/output.ini")
```

## How to Test

- run go test ./... in root directory

```golang
go test ./...
```

- add the -v flag for more details about the specific tests that are running

```golang
go test -v ./...
```
