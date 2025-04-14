package connection

import (
	"bufio"
	"container/list"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/config"
)

type node_address struct {
	Node_id []byte
	Address string
}

func (node *Node) get_nodes(key []byte, peer_address string) []node_address {

	if peer_address == config.MetaData.ListeningAddress {
		return nil
	}
	conn, err := net.Dial("tcp", peer_address)

	if err != nil {
		// fmt.Println(err)
		return nil
	}
	defer conn.Close()

	conn.Write([]byte(fmt.Sprintf("CLOSEST %s\n", hex.EncodeToString(key))))
	reader := bufio.NewReader(conn)
	msg, err := reader.ReadString('\n')
	msg = strings.TrimSuffix(msg, "\n")
	if err != nil {
		fmt.Println("Connection closed or error:", err)
		return nil
	}
	// fmt.Println("msg", msg)
	ans := make([]node_address, 0)
	if num, e := strconv.Atoi(msg); e == nil {
		for i := 0; i < num; i++ {
			msg, err = reader.ReadString('\n')
			msg = strings.TrimSuffix(msg, "\n")
			// fmt.Println("msg", msg)
			if err != nil {
				fmt.Println("Connection closed or error:", err)
				return ans
			}
			parts := strings.Split(msg, " ")
			id, e := hex.DecodeString(parts[0])
			if e == nil {
				if node.Bucket.Insert_NodeID(id) {
					NodeIDtoNetConn[hex.EncodeToString(id)] = parts[1]
				}
				ans = append(ans, node_address{Node_id: id, Address: parts[1]})
			}
		}
	}
	return ans
}

func xor_dist(node_id1 []byte, node_id2 []byte) []byte {
	new_byte := make([]byte, 32)
	for i := 0; i < 32; i++ {
		new_byte[i] = node_id1[i] ^ node_id2[i]
	}
	return new_byte
}

func Comp(i []byte, j []byte) bool {
	for idx := 31; idx >= 0; idx-- {
		if int(i[idx]) < int(j[idx]) {
			return true
		} else if int(i[idx]) > int(j[idx]) {
			return false
		}
	}
	return true
}

func (node *Node) get_closest_nodes(key []byte) []node_address {
	nodes := node.Bucket.Find_Nodes(key)

	closest_list := list.New()
	for i := 0; i < len(nodes); i++ {
		closest_list.PushBack(node_address{Node_id: nodes[i], Address: NodeIDtoNetConn[hex.EncodeToString(nodes[i])]})

	}
	visited := make(map[string]bool, 0)

	for {
		dis := make([]byte, 32)

		var wg sync.WaitGroup
		new_nodes := map[string]node_address{}
		f := 0
		for it := 0; it < config.MetaData.SearchAlpha; it++ {
			var mi node_address
			for i := 0; i < 32; i++ {
				dis[i] = 1<<8 - 1
			}
			// find min dist
			for e := closest_list.Front(); e != nil; e = e.Next() {
				k := e.Value.(node_address)
				if !visited[hex.EncodeToString(k.Node_id)] && Comp(xor_dist(k.Node_id, key), dis) {
					mi = k
					dis = xor_dist(k.Node_id, key)
					f = 1
				}
			}

			if f != 0 {
				visited[hex.EncodeToString(mi.Node_id)] = true
				wg.Add(1)
				func(mi node_address) {
					defer wg.Done()
					new_grp := node.get_nodes(key, mi.Address)

					for i := 0; i < len(new_grp); i++ {
						new_nodes[hex.EncodeToString(new_grp[i].Node_id)] = new_grp[i]
						// fmt.Println("received ", hex.EncodeToString(new_grp[i].Node_id), new_grp[i].Address)
					}
				}(mi)
			} else {
				break
			}
		}
		if f == 0 {
			break
		}
		wg.Wait()
		f = 0

		for _, v := range new_nodes {
			for i := 0; i < 32; i++ {
				dis[i] = 0
			}
			if closest_list.Len() < config.MetaData.BucketSize {
				if !visited[hex.EncodeToString(v.Node_id)] {
					closest_list.PushBack(v)
				}
				continue
			}
			var mx node_address
			mxid := -1
			for e := closest_list.Front(); e != nil; e = e.Next() {
				k := e.Value.(node_address)
				if !Comp(xor_dist(key, k.Node_id), dis) {
					dis = xor_dist(key, k.Node_id)
					mx = k
					mxid = 1
				}
			}
			if mxid == 1 {
				if Comp(xor_dist(v.Node_id, key), dis) {
					for e := closest_list.Front(); e != nil; e = e.Next() {
						k := e.Value.(node_address)
						if hex.EncodeToString(k.Node_id) == hex.EncodeToString(mx.Node_id) {
							e.Value = v
							f = 1
						}
					}
				}
			}
		}
		if f == 0 {
			break
		}
	}

	vec := make([]node_address, 0)
	for e := closest_list.Front(); e != nil; e = e.Next() {
		k := e.Value.(node_address)
		vec = append(vec, node_address{Node_id: k.Node_id, Address: k.Address})
	}
	return vec
}

func (node *Node) handel_closest(parts []string, conn net.Conn) string {
	// CLOSEST hex_id
	id, err := hex.DecodeString(parts[1])
	if err == nil {
		nodes := node.Bucket.Find_Nodes(id)
		x := strconv.Itoa(len(nodes)) + "\n"
		for i := 0; i < len(nodes); i++ {
			x += hex.EncodeToString(nodes[i]) + " " + NodeIDtoNetConn[hex.EncodeToString(nodes[i])]
			if i != len(nodes)-1 {
				x += "\n"
			}
		}
		return x

	} else {
		return ""
	}

}
