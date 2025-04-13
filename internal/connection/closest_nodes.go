package connection

import (
	"encoding/hex"
	"strconv"
)

func (node *Node) get_closest_nodes(key []byte, peer_node_id []byte) [][]byte {

	return nil
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
