package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/config"
	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/connection"
)

func main() {
	config.InitConfig()

	dir := "./files/" + strconv.Itoa(int(config.MetaData.Port)) + "/"
	logFile, err := os.OpenFile(dir+"log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)

	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	node := connection.NewNode()

	node.Start()
}
