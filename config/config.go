package config

import (
	"flag"
	"fmt"
	"os"
)

type metadata struct {
	Type int
}

var MetaData metadata

func ReadFlags() {
	var type_flag string
	flag.StringVar(&type_flag, "type", "none", "Type of node either storage/lookup")

	flag.Parse()

	if type_flag == "storage" {
		MetaData.Type = 1
	} else if type_flag == "lookup" {
		MetaData.Type = 2
	} else if type_flag == "none" {
		fmt.Println("No type flag provided")
		os.Exit(1)
	} else {
		fmt.Println("type flag not recognised")
		os.Exit(1)
	}
}
