package config

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"math/rand"
	"os"
)

type metadata struct {
	NodeID [32]byte
	Port   uint16
}

var MetaData metadata

func InitConfig() {
	var restart string
	var port_num int
	flag.StringVar(&restart, "restart", "none", "Restarting or completely new")
	flag.IntVar(&port_num, "port", 0, "The port to listen to")
	flag.Parse()

	if restart == "flase" {
		// TODO read config from a file and fill metadata

		return
	} else if restart != "none" {
		fmt.Println("Unknown restart flag")
		os.Exit(1)
	}

	if port_num <= 49151 && port_num >= 1024 {
		MetaData.Port = uint16(port_num)
	} else {
		fmt.Println("port is not valid")
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
