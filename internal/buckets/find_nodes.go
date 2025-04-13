package buckets

import (
	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/config"
)

type heap_node struct {
	node_id  []byte
	xor_dist []byte
}

type BigIntHeap []heap_node

func (h BigIntHeap) Len() int { return len(h) }
func (h BigIntHeap) Less(i, j int) bool {
	for idx := 31; idx >= 0; idx-- {
		if int(h[i].xor_dist[idx]) < int(h[j].xor_dist[idx]) {
			return true
		} else if int(h[i].xor_dist[idx]) > int(h[j].xor_dist[idx]) {
			return false
		}
	}
	return true
}

func (h BigIntHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *BigIntHeap) Push(x heap_node) {
	*h = append(*h, x)
}

func (h *BigIntHeap) Pop() heap_node {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func xor_dist(node_id1 []byte, node_id2 []byte) []byte {
	new_byte := make([]byte, 32)
	for i := 0; i < 32; i++ {
		new_byte[i] = node_id1[i] ^ node_id2[i]
	}
	return new_byte
}

func (buckets *Buckets) Find_Nodes(node_id []byte) [][]byte {

	nodes := make([][]byte, 0)
	pq := make(BigIntHeap, 0)

	for idx := 0; idx < 256; idx++ {
		for e := buckets.buckets_lists[idx].Front(); e != nil; e = e.Next() {
			k := e.Value.([]byte)
			pq.Push(heap_node{node_id: k, xor_dist: xor_dist(node_id, k)})
		}
	}
	for i := 0; i < config.MetaData.BucketSize && pq.Len() > 0; i++ {
		nodes = append(nodes, pq.Pop().node_id)
	}
	return nodes
}
