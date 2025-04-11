package connection

import (
	"fmt"
	"net"
	"strconv"

	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/config"
)

type Node struct {
	Listening_address string
	Parent_connector  Connector
}

func NewNode() *Node {
	return &Node{
		Listening_address: "",
	}
}

type Connector interface {
	Handel_conn(net.Conn)
}

func (node *Node) start() error {
	l, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(int(config.MetaData.Port)))
	if err != nil {
		return err
	}

	node.Listening_address = l.Addr().String()

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go node.Parent_connector.Handel_conn(conn)
	}
}
