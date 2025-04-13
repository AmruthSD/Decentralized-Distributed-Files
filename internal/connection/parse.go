package connection

import (
	"encoding/hex"
	"strings"

	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/config"
)

func (node *Node) parse(msg string) string {
	var parse_func = map[string]func([]string) string{
		"PING":         node.handle_ping,
		"unknown":      node.handle_unknown,
		"SEND_NODE_ID": node.handle_node_id,
		"CLOSEST":      node.handel_closest,
	}
	parts := strings.Split(msg, " ")
	f, ex := parse_func[parts[0]]
	if ex {
		return f(parts)
	}
	return "unknown"
}

func (node *Node) handle_ping(parts []string) string {
	return "PONG"
}

func (node *Node) handle_unknown(msg []string) string {
	return "STOP"
}

func (node *Node) handle_node_id(parts []string) string {
	if len(parts) == 3 {
		node_id := []byte(parts[1])
		node_listening := parts[2]
		if node.Bucket.Insert_NodeID(node_id) {
			NodeIDtoNetConn[hex.EncodeToString(node_id)] = node_listening
		}
	}
	return hex.EncodeToString(config.MetaData.NodeID) + " " + config.MetaData.ListeningAddress
}
