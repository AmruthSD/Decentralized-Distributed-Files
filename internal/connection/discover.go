package connection

import (
	"encoding/hex"

	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/config"
)

func (node *Node) Handel_discover() {
	node.Bucket.Insert_NodeID(config.MetaData.NodeID)
	if config.MetaData.WellKnownPort != config.MetaData.Port {
		node_id := make([]byte, 32)
		copy(node_id, config.MetaData.NodeID)
		for i := 1; i < 256; i++ {
			id := i / 8
			msk := i % 8

			node_id[id] = node_id[id] ^ 1<<msk
			list_node_id := node.get_closest_nodes(node_id, config.MetaData.WellKnownListeningAddress)
			node_id[id] = node_id[id] ^ 1<<msk

			for v := range len(list_node_id) {
				if node.Bucket.Insert_NodeID(list_node_id[v].Node_id) {
					NodeIDtoNetConn[hex.EncodeToString(list_node_id[v].Node_id)] = string(list_node_id[v].Node_id)
				}
			}
		}
	}
}
