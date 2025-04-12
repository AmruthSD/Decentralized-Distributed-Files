package connection

import (
	"bufio"
	"fmt"
	"net"
	"strconv"

	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/buckets"
	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/config"
)

type Node struct {
	Listening_address string
	Bucket            buckets.Buckets
}

func NewNode() *Node {
	return &Node{
		Listening_address: "",
		Bucket:            *buckets.NewBuckets(),
	}
}

func (node *Node) Start() error {
	l, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(int(config.MetaData.Port)))
	if err != nil {
		return err
	}

	node.Listening_address = l.Addr().String()

	defer l.Close()

	node.Handel_discover()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go node.Handel_conn(conn)
	}
}

func (node *Node) Handel_conn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection closed or error:", err)
			return
		}

		fmt.Println("Received:", msg)
		msg = parse(msg)
		if msg == "STOP" {
			break
		}
		conn.Write([]byte(msg + "\n"))
	}
}
