package main

import (
	"log"

	iniparser "github.com/codescalersinternships/INIParser-EYAD_HUSSEIN/pkg"
)

func main() {
	parser := iniparser.NewParser()
	err := parser.LoadFromFile("../pkg/test_file_1.ini")

	if err != nil {
		log.Fatal(err)
	}

	log.Println(parser.Get("owner", "name"))
}
