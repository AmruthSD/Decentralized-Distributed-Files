package connection

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/buckets"
	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/config"
)

var MapMutex sync.Mutex
var NodeIDtoNetConn = map[string]string{}

type Node struct {
	Bucket buckets.Buckets
}

func NewNode() *Node {
	return &Node{
		Bucket: *buckets.NewBuckets(),
	}
}

func (node *Node) Start() error {
	l, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(int(config.MetaData.Port)))
	if err != nil {
		return err
	}
	fmt.Println("Started Listening At:", l.Addr().String())

	config.MetaData.ListeningAddress = l.Addr().String()
	h := hex.EncodeToString(config.MetaData.NodeID)
	MapMutex.Lock()
	NodeIDtoNetConn[h] = config.MetaData.ListeningAddress
	MapMutex.Unlock()
	node.Dial_Well_Known()

	defer l.Close()

	node.Handel_discover()

	fmt.Println("Finished Discover")
	go node.Handle_Client()
	go node.Handle_persist()
	go node.Handle_KeepAlive()
	go node.Handle_DeleteExpire()

	fmt.Println("Started to Accpet")
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

	// fmt.Println("New Connection", conn.RemoteAddr().String())
	for {
		msg, err := reader.ReadString('\n')
		msg = strings.TrimSuffix(msg, "\n")
		if err != nil {
			// fmt.Println("Connection closed or error:", err)
			return
		}

		// fmt.Println("Received:", msg)
		msg = node.parse(msg, conn)
		if msg == "STOP" {
			break
		}
		conn.Write([]byte(msg + "\n"))
	}
}

func (node *Node) Dial_Well_Known() {
	if config.MetaData.WellKnownPort == config.MetaData.Port {
		return
	}
	conn, err := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(int(config.MetaData.WellKnownPort)))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	conn.Write([]byte(fmt.Sprintf("SEND_NODE_ID %s %s\n", hex.EncodeToString(config.MetaData.NodeID), config.MetaData.ListeningAddress)))
	reader := bufio.NewReader(conn)
	msg, err := reader.ReadString('\n')
	msg = strings.TrimSuffix(msg, "\n")
	if err != nil {
		fmt.Println("Connection closed or error:", err)
		return
	}
	fmt.Println("Received:", msg)
	parts := strings.Split(msg, " ")

	if len(parts) == 2 {
		id_decoded, err := hex.DecodeString(parts[0])
		config.MetaData.WellKnownListeningAddress = parts[1]
		if err == nil {
			if node.Bucket.Insert_NodeID(id_decoded) {
				MapMutex.Lock()
				NodeIDtoNetConn[parts[0]] = parts[1]
				MapMutex.Unlock()
			}
		}
	}
}
