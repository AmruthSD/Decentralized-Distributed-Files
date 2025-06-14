package connection

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/config"
)

// gets the well known node's port and id and listening
func (node *Node) Dial_Well_Known() {

	if config.MetaData.WellKnownPort == config.MetaData.Port {
		return
	}

	// connect to a well known port
	conn, err := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(int(config.MetaData.WellKnownPort)))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// send my own id for discovery
	conn.Write([]byte(fmt.Sprintf("SEND_NODE_ID %s %s\n", hex.EncodeToString(config.MetaData.NodeID), config.MetaData.ListeningAddress)))
	reader := bufio.NewReader(conn)

	// gets the id and listening address of the other person
	msg, err := reader.ReadString('\n')
	msg = strings.TrimSuffix(msg, "\n")
	if err != nil {
		log.Println("Connection closed or error:", err)
		return
	}
	log.Println("Received:", msg)
	parts := strings.Split(msg, " ")

	if len(parts) == 2 {
		// populate the well known port's address
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
