package connection

import (
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/config"
)

func (node *Node) Handle_persist() {
	dir := "./files/" + strconv.Itoa(int(config.MetaData.Port)) + "/storage/"
	for {
		entries, err := os.ReadDir(dir)
		if err != nil {
			fmt.Println("Error reading directory:", err)
			continue
		}
		for _, v := range entries {
			fileName := v.Name()
			cid, _ := hex.DecodeString(fileName)

			nodes := node.get_closest_nodes(cid)
			delete_cid := 1
			file, _ := os.Open(dir + fileName)
			for _, x := range nodes {
				if hex.EncodeToString(x.Node_id) == hex.EncodeToString(config.MetaData.NodeID) {
					delete_cid = 0
					continue
				}
				conn, err := net.Dial("tcp", x.Address)
				if err != nil {
					continue
				}
				conn.Write([]byte(fmt.Sprintf("DOYOUHAVE %s CHECK\n", fileName)))
				readbuff := make([]byte, config.MetaData.ChunkSize)
				n, _ := conn.Read(readbuff)
				if string(readbuff[:n]) == "NO\n" {
					n, _ := file.Read(readbuff)
					conn.Write(readbuff[:n])
				}
			}
			if delete_cid == 1 {
				os.Remove(dir + fileName)
			}
		}

		time.Sleep(time.Duration(config.MetaData.TimeOut) * time.Minute)
	}
}
