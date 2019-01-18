package stakesort

import (
	"container/heap"
	"fmt"
	"math/big"
	"testing"
)

func TestPq(t *testing.T){

	items := []*Item{
		{"a",big.NewInt(10),0},
		{"b",big.NewInt(11),1},
		{"c",big.NewInt(20),2},
		{"d",big.NewInt(21),2},
	}
	priorityqueue := make(PriorityQueue,0)
	heap.Init(&priorityqueue)
	for _,item := range items{
		heap.Push(&priorityqueue,item)

	}
	for priorityqueue.Len() > 0{
		item := heap.Pop(&priorityqueue).(*Item)
		fmt.Printf("addr : %s value %d\n",item.Address,item.Value)
	}
}