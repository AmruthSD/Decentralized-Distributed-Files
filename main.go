package main

import (
	"github.com/AmruthSD/Decentralized-Distributed-Files/config"
)

func main() {
	config.InitConfig()

	if config.MetaData.Type == 1 {
		// go storage
	} else if config.MetaData.Type == 2 {
		// go lookup
	}
}
