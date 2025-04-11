package config

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"math/rand"
	"os"
)

type metadata struct {
	Type   int
	NodeID [32]byte
}

var MetaData metadata

func InitConfig() {
	var type_flag, restart string
	flag.StringVar(&type_flag, "type", "none", "Type of node either storage/lookup")
	flag.StringVar(&restart, "restart", "none", "Restarting or completely new")
	flag.Parse()

	if restart == "flase" {
		// TODO read config from a file and fill metadata

		return
	} else if restart != "none" {
		fmt.Println("Unknown restart flag")
		os.Exit(1)
	}

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
	MetaData.generate_new_node_id()
}

func (MetaData *metadata) generate_new_node_id() {
	bytes_random := make([]byte, 32)
	_, err := rand.Read(bytes_random)
	if err != nil {
		fmt.Println("error in random read")
		os.Exit(1)
	}

	MetaData.NodeID = sha256.Sum256(bytes_random)

}
