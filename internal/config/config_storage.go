package config

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
)

func restart_node(port int) bool {
	dir := "./files/" + strconv.Itoa(int(port)) + "/"
	configFile := dir + "config.txt"

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return false
	}
	file, err := os.Open(configFile)
	if err != nil {
		return true
	}
	readID := ""
	portID := 0
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		readID = scanner.Text()
	} else {
		return false
	}
	if scanner.Scan() {
		portID, err = strconv.Atoi(scanner.Text())
		if err != nil {
			panic(err)
		}
	} else {
		return false
	}
	if portID == port {
		MetaData.NodeID, err = hex.DecodeString(readID)
		if err != nil {
			return false
		}
	}
	file, _ = os.OpenFile(configFile, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	file.Write([]byte(fmt.Sprintf("%s\n%d\n", hex.EncodeToString(MetaData.NodeID), MetaData.Port)))
	file.Close()
	fmt.Println("Loading from previous session with same port")
	return true
}
