package connection

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/config"
)

// send keep alive message every keep alive time to make sure that the chunks are not deleted
func (node *Node) Handle_KeepAlive() {
	dir := "./files/" + strconv.Itoa(int(config.MetaData.Port)) + "/hashed/"
	for {
		entries, err := os.ReadDir(dir)
		if err != nil {
			log.Println("Error reading directory:", err)
			continue
		}
		for _, v := range entries {
			fileName := v.Name()
			file, err := os.Open(dir + fileName)
			if err != nil {
				log.Println("Error opening file", err)
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

// if the chuck is not kept alive then just delete every 48 hrs
func (node *Node) Handle_DeleteExpire() {
	dir := "./files/" + strconv.Itoa(int(config.MetaData.Port)) + "/"
	for {
		entries, err := os.ReadDir(dir + "storage")
		if err != nil {
			log.Println("Error reading directory:", err)
			continue
		}
		file, _ := os.Open(dir + "storage.json")
		data := map[string]time.Time{}
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&data); err != nil {
			log.Println("Error decoding JSON:", err)
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

// open the json and then update the keep alive for that chunk
func UpdateTimeStamp(cid string) {
	dir := "./files/" + strconv.Itoa(int(config.MetaData.Port)) + "/"
	file, _ := os.Open(dir + "storage.json")
	data := map[string]time.Time{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		log.Println("Error decoding JSON:", err)
		return
	}
	file.Close()

	expireTime := time.Now().Add(48 * time.Hour)
	outFile, err := os.OpenFile(dir+"storage.json", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Println("Error ", err)
		return
	}
	data[cid] = expireTime
	encoder := json.NewEncoder(outFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		log.Println("Error encoding JSON:", err)
	}
	outFile.Close()
}

// makes sure that the chunk is not deleted
func (node *Node) handle_keepalive(parts []string, conn net.Conn) string {
	UpdateTimeStamp(parts[1])
	return "STOP"
}
