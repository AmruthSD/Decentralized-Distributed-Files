package lookupnode

import (
	"net"

	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/connection"
)

type lookupnode struct {
	node connection.Node
}

func (lookup *lookupnode) Handel_conn(conn net.Conn) {

}

func NewLookUpNode() *lookupnode {
	x := &lookupnode{
		node: *connection.NewNode(),
	}
	x.node.Parent_connector = x
	return x
}
