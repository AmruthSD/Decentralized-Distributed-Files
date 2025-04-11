package connection

import (
	"fmt"
	"net"
	"strconv"

	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/config"
)

type Node struct {
	conn              net.Conn
	listening_address string
	parent_connector  Connector
}

type Connector interface {
	handel_conn(net.Conn)
}

func (node *Node) start() error {
	l, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(int(config.MetaData.Port)))
	if err != nil {
		return err
	}

	node.listening_address = l.Addr().String()

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go node.parent_connector.handel_conn(conn)
	}
}
