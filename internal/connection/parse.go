package connection

import (
	"strings"

	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/config"
)

var parse_func = map[string]func([]string) string{
	"PING":         handle_ping,
	"unknown":      handle_unknown,
	"SEND_NODE_ID": handle_node_id,
}

func parse(msg string) string {

	parts := strings.Split(msg, " ")
	f, ex := parse_func[parts[0]]
	if ex {
		return f(parts)
	}
	return "unknown"
}

func handle_ping(parts []string) string {
	return "PONG"
}

func handle_unknown(msg []string) string {
	return "STOP"
}

func handle_node_id(parts []string) string {
	return string(config.MetaData.NodeID)
}
