package storagenode

import (
	"net"

	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/connection"
)

type storagenode struct {
	node connection.Node
}

func (lookup *storagenode) Handel_conn(conn net.Conn) {

}

func NewStorageNode() *storagenode {
	x := &storagenode{
		node: *connection.NewNode(),
	}
	x.node.Parent_connector = x
	return x
}
