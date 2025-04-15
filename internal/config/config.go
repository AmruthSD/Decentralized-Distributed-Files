package config

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

type metadata struct {
	NodeID                    []byte
	Port                      uint16
	WellKnownPort             uint16
	BucketSize                int
	ListeningAddress          string
	WellKnownListeningAddress string
	SearchAlpha               int
	ChunkSize                 int
	TimeOutPersist            int
	TimeOutKeepAlive          int
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
	fmt.Println("NodeID:", hex.EncodeToString(MetaData.NodeID))
	MetaData.WellKnownPort = 8000
	MetaData.BucketSize = 20
	MetaData.SearchAlpha = 3
	MetaData.ChunkSize = 4 * 1024
	MetaData.WellKnownListeningAddress = "[::]:8000"
	MetaData.TimeOutPersist = 10
	MetaData.TimeOutKeepAlive = 20

	dir := "./files/" + strconv.Itoa(int(MetaData.Port)) + "/"
	err := os.MkdirAll(dir+"downloaded", 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
	}
	err = os.MkdirAll(dir+"hashed", 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
	}
	err = os.MkdirAll(dir+"storage", 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
	}
	file, err := os.OpenFile(dir+"storage.json", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsExist(err) {
			fmt.Println("Storage Json already exists, not overwriting.")
			return
		}
		fmt.Println("Error creating file:", err)
		return
	}
	file.Write([]byte("{}"))
	file.Close()
}

func (MetaData *metadata) generate_new_node_id() {
	bytes_random := make([]byte, 32)
	_, err := rand.Read(bytes_random)
	if err != nil {
		fmt.Println("error in random read")
		os.Exit(1)
	}

	arr := sha256.Sum256(bytes_random)
	MetaData.NodeID = arr[:]
}
