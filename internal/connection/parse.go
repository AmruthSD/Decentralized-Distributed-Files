package connection

import (
	"encoding/hex"
	"fmt"
	"net"
	"strings"

	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/config"
)

func (node *Node) parse(msg string, conn net.Conn) string {
	var parse_func = map[string]func([]string, net.Conn) string{
		"PING":         node.handle_ping,
		"unknown":      node.handle_unknown,
		"SEND_NODE_ID": node.handle_node_id,
		"CLOSEST":      node.handel_closest,
		"STORE":        node.handel_store,
		"DONE":         node.handle_done,
		"DOYOUHAVE":    node.handle_doyouhave,
		"KEEPALIVE":    node.handle_keepalive,
	}
	parts := strings.Split(msg, " ")
	f, ex := parse_func[parts[0]]
	if ex {
		return f(parts, conn)
	}
	return "unknown"
}

func (node *Node) handle_ping(parts []string, conn net.Conn) string {
	return "PONG"
}

func (node *Node) handle_unknown(msg []string, conn net.Conn) string {
	return "STOP"
}

func (node *Node) handle_done(msg []string, conn net.Conn) string {
	return "STOP"
}

func (node *Node) handle_node_id(parts []string, conn net.Conn) string {
	if len(parts) == 3 {
		node_id, _ := hex.DecodeString(parts[1])
		node_listening := parts[2]
		if node.Bucket.Insert_NodeID(node_id) {
			MapMutex.Lock()
			NodeIDtoNetConn[hex.EncodeToString(node_id)] = node_listening
			MapMutex.Unlock()
			fmt.Println("Inserting:", parts[1], parts[2])
		}
	}
	return hex.EncodeToString(config.MetaData.NodeID) + " " + config.MetaData.ListeningAddress
}
