package main

import (
	iniparser "github.com/codescalersinternships/INIParser-EYAD_HUSSEIN/pkg"
)

func main() {
	parser := iniparser.IniParser{}
	parser.LoadFromFile("./data.ini")
}
