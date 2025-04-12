package buckets

import (
	"container/list"

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

func (buckets *Buckets) Insert_NodeID(node_id []byte) {
	if len(node_id) != 32 {
		return
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
		buckets.buckets_lists[bucket_num].PushBack(node_id)
	} else {
		// send ping and if not respond remove
	}

}
