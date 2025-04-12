package main

import (
	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/config"
	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/connection"
)

func main() {
	config.InitConfig()

	node := connection.NewNode()

	node.Start()
}
