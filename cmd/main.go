package main

import (
	"log"

	iniparser "github.com/codescalersinternships/INIParser-EYAD_HUSSEIN/pkg"
)

func main() {
	parser := iniparser.IniParser{}
	err := parser.LoadFromFile("./data.ini")

	if err != nil {
		log.Fatal(err)
	}
}
