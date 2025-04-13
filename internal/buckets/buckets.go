package buckets

import (
	"container/list"
	"encoding/hex"

	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/config"
)

type Buckets struct {
	buckets_lists []*list.List
}

func NewBuckets() *Buckets {
	buc := make([]*list.List, 256)
	for i := 0; i < 256; i++ {
		buc[i] = list.New()
	}
	return &Buckets{
		buckets_lists: buc,
	}
}

func (buckets *Buckets) Insert_NodeID(node_id []byte) bool {
	if len(node_id) != 32 {
		return false
	}
	own_node_id := config.MetaData.NodeID

	bucket_num := 0
	for i := 255; i >= 0; i-- {
		idx := i / 8
		msk := i % 8

		if node_id[idx]&1<<msk != own_node_id[idx]&1<<msk {
			bucket_num = i
			break
		}
	}

	if buckets.buckets_lists[bucket_num].Len() < config.MetaData.BucketSize {
		for e := buckets.buckets_lists[bucket_num].Front(); e != nil; e = e.Next() {
			if hex.EncodeToString(e.Value.([]byte)) == hex.EncodeToString(node_id) {
				return true
			}
		}
		buckets.buckets_lists[bucket_num].PushBack(node_id)
		return true
	} else {
		// send ping and if not respond remove
		return false
	}
}
