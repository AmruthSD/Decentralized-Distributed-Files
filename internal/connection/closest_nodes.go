package connection

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type node_address struct {
	Node_id []byte
	Address string
}

func (node *Node) get_closest_nodes(key []byte, peer_address string) []node_address {

	conn, err := net.Dial("tcp", peer_address)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	conn.Write([]byte(fmt.Sprintf("CLOSEST %s\n", hex.EncodeToString(key))))
	reader := bufio.NewReader(conn)
	msg, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Connection closed or error:", err)
		return nil
	}
	ans := make([]node_address, 0)
	if num, e := strconv.Atoi(msg); e == nil {
		for i := 0; i < num; i++ {
			msg, err = reader.ReadString('\n')
			if err != nil {
				fmt.Println("Connection closed or error:", err)
				return ans
			}
			parts := strings.Split(msg, " ")
			id, e := hex.DecodeString(parts[0])
			if e == nil {
				if node.Bucket.Insert_NodeID(id) {
					NodeIDtoNetConn[hex.EncodeToString(id)] = parts[1]
				}
				ans = append(ans, node_address{Node_id: id, Address: parts[1]})
			}
		}
	}
	return ans
}

func (node *Node) handel_closest(parts []string) string {
	// CLOSEST hex_id
	id, err := hex.DecodeString(parts[1])
	if err == nil {
		nodes := node.Bucket.Find_Nodes(id)
		x := strconv.Itoa(len(nodes)) + "\n"
		for i := 0; i < len(nodes); i++ {
			x += hex.EncodeToString(nodes[i]) + " " + NodeIDtoNetConn[hex.EncodeToString(nodes[i])]
			if i != len(nodes)-1 {
				x += "\n"
			}
		}
		return x

	} else {
		return ""
	}

}
