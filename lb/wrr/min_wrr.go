package wrr

import (
	"container/heap"
	"sync"
)

type minWRR struct {
	mu      sync.RWMutex
	items   minHeap
	itemMap map[any]*Item
}

func NewMinWRR() WRR {
	h := minHeap{}
	heap.Init(&h)
	return &minWRR{
		items:   h,
		itemMap: map[any]*Item{},
	}
}

func (mw *minWRR) Next() (item interface{}) {
	mw.mu.RLock()
	defer mw.mu.RUnlock()
	if len(mw.items) == 0 {
		return nil
	}

	it, ok := mw.items.Peak().(*Item)
	if !ok {
		panic("invalid min heap")
	}

	return it.value
}

func (mw *minWRR) Add(item any, weight uint32) {
	mw.mu.Lock()
	defer mw.mu.Unlock()

	it, ok := mw.itemMap[item]
	if ok {
		mw.items.update(it, int(weight))
		return
	}

	it = &Item{
		value:    item,
		priority: it.priority + int(weight),
	}
	heap.Push(&mw.items, it)
	mw.itemMap[item] = it
}

func (mw *minWRR) Update(item any, weight uint32) {
	mw.mu.Lock()
	defer mw.mu.Unlock()

	it, ok := mw.itemMap[item]
	if ok {
		mw.items.update(it, int(weight))
		return
	}

	it = &Item{
		value:    item,
		priority: int(weight),
	}
	heap.Push(&mw.items, it)
	mw.itemMap[item] = it
}

type Item struct {
	value    any
	priority int
	index    int
}

type minHeap []*Item

func (mh minHeap) Len() int { return len(mh) }

func (mh minHeap) Less(i, j int) bool {
	return mh[i].priority > mh[j].priority
}

func (mh minHeap) Swap(i, j int) {
	mh[i], mh[j] = mh[j], mh[i]
	mh[i].index = i
	mh[j].index = j
}

func (mh *minHeap) Push(x any) {
	n := len(*mh)
	item := x.(*Item)
	item.index = n
	*mh = append(*mh, item)
}

func (mh *minHeap) Pop() any {
	old := *mh
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*mh = old[0 : n-1]
	return item
}

func (mh *minHeap) Peak() any {
	return (*mh)[len(*mh)-1]
}

func (mh *minHeap) update(item *Item, priority int) {
	item.priority = priority
	heap.Fix(mh, item.index)
}
