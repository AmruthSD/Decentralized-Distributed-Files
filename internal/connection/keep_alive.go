package connection

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/config"
)

func (node *Node) Handle_KeepAlive() {
	dir := "./files/" + strconv.Itoa(int(config.MetaData.Port)) + "/hashed/"
	for {
		entries, err := os.ReadDir(dir)
		if err != nil {
			fmt.Println("Error reading directory:", err)
			continue
		}
		for _, v := range entries {
			fileName := v.Name()
			file, err := os.Open(dir + fileName)
			if err != nil {
				fmt.Println("Error opening file", err)
				continue
			}
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				cid, _ := hex.DecodeString(line)
				nodes := node.get_closest_nodes(cid)
				for _, node := range nodes {
					conn, _ := net.Dial("tcp", node.Address)
					conn.Write([]byte(fmt.Sprintf("KEEPALIVE %s 48", line)))
				}
			}

			file.Close()
		}

		time.Sleep(time.Duration(config.MetaData.TimeOutKeepAlive) * time.Hour)
	}
}

func (node *Node) Handle_DeleteExpire() {
	dir := "./files/" + strconv.Itoa(int(config.MetaData.Port)) + "/"
	for {
		entries, err := os.ReadDir(dir + "storage")
		if err != nil {
			fmt.Println("Error reading directory:", err)
			continue
		}
		file, _ := os.Open(dir + "storage.json")
		data := map[string]time.Time{}
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&data); err != nil {
			fmt.Println("Error decoding JSON:", err)
			time.Sleep(48 * time.Hour)
			continue
		}
		file.Close()

		for _, v := range entries {
			fileName := v.Name()
			if time.Now().After(data[fileName]) {
				os.Remove(dir + "storage/" + "fileName")
			}
		}
		time.Sleep(48 * time.Hour)
	}
}

func UpdateTimeStamp(cid string) {
	dir := "./files/" + strconv.Itoa(int(config.MetaData.Port)) + "/"
	file, _ := os.Open(dir + "storage.json")
	data := map[string]time.Time{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	file.Close()

	expireTime := time.Now().Add(48 * time.Hour)
	outFile, err := os.OpenFile(dir+"storage.json", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error ", err)
		return
	}
	data[cid] = expireTime
	encoder := json.NewEncoder(outFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		fmt.Println("Error encoding JSON:", err)
	}
	outFile.Close()
}

func (node *Node) handle_keepalive(parts []string, conn net.Conn) string {
	UpdateTimeStamp(parts[1])
	return "STOP"
}
