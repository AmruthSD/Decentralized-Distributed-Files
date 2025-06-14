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

// sarting the node actions
func (node *Node) Start() error {
	l, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(int(config.MetaData.Port)))
	if err != nil {
		return err
	}
	fmt.Println("Started Listening At:", l.Addr().String())

	config.MetaData.ListeningAddress = l.Addr().String()
	h := hex.EncodeToString(config.MetaData.NodeID)

	// adding own listening address for own node id
	MapMutex.Lock()
	NodeIDtoNetConn[h] = config.MetaData.ListeningAddress
	MapMutex.Unlock()

	node.Dial_Well_Known()

	defer l.Close()

	node.Handel_discover()

	fmt.Println("Finished Discover")

	// starting the concurrent routines
	go node.Handle_Client()
	go node.Handle_persist()
	go node.Handle_KeepAlive()
	go node.Handle_DeleteExpire()

	fmt.Println("Started to Accpet")
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err.Error())
			continue
		}

		go node.Handel_conn(conn)
	}
}

// infite loop read and send msg
func (node *Node) Handel_conn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	log.Println("New Connection", conn.RemoteAddr().String())
	for {
		msg, err := reader.ReadString('\n')
		msg = strings.TrimSuffix(msg, "\n")
		if err != nil {
			log.Println("Connection closed or error:", err)
			return
		}

		log.Println("Received:", msg)
		msg = node.parse(msg, conn)
		if msg == "STOP" {
			break
		}
		conn.Write([]byte(msg + "\n"))
	}
}
