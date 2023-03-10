// Source: https://pkg.go.dev/container/heap#example-package-PriorityQueue

package worker

import (
	"container/heap"

	"github.com/badochov/distributed-shortest-path/src/libs/db"
)

// An Item is something we manage in a distance queue.
type Item struct {
	id       db.VertexId // The id of the item; arbitrary.
	distance float64     // The distance of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the lowest distance.
	return pq[i].distance < pq[j].distance
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the distance and id of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, distance float64) {
	item.distance = distance
	heap.Fix(pq, item.index)
}
