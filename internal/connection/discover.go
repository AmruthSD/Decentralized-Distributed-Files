package connection

import (
	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/config"
)

func (node *Node) Handel_discover() {
	node.Bucket.Insert_NodeID(config.MetaData.NodeID)
	if config.MetaData.WellKnownPort != config.MetaData.Port {
		// node_id := make([]byte, 32)
		// copy(node_id, config.MetaData.NodeID)
		// for i := 1; i < 256; i++ {

		// 	// do node.get nodes with k nearest

		// 	var list_node_id [][]byte
		// 	for v := range len(list_node_id) {
		// 		node.Bucket.Insert_NodeID(list_node_id[v])
		// 	}
		// }
	}
}
