package wrr

import (
	"container/heap"
	"sync"
)

type minWRR struct {
	mu         sync.RWMutex
	elements   minHeap
	elementMap map[any]*Element // key is item, which is deliered from Add or Update
}

func NewMinWRR() WRR {
	h := minHeap{}
	heap.Init(&h)
	return &minWRR{
		elements:   h,
		elementMap: map[any]*Element{},
	}
}

func (mw *minWRR) Next() (item interface{}) {
	mw.mu.RLock()
	defer mw.mu.RUnlock()
	if len(mw.elements) == 0 {
		return nil
	}

	ele := mw.elements.peak()

	return ele.item
}

func (mw *minWRR) Add(item any, weight uint32) {
	mw.mu.Lock()
	defer mw.mu.Unlock()

	it, ok := mw.elementMap[item]
	if ok {
		mw.elements.update(it, int(weight))
		return
	}

	it = &Element{
		item:     item,
		priority: it.priority + int(weight),
	}
	heap.Push(&mw.elements, it)
	mw.elementMap[item] = it
}

func (mw *minWRR) Update(item any, weight uint32) {
	mw.mu.Lock()
	defer mw.mu.Unlock()

	it, ok := mw.elementMap[item]
	if ok {
		mw.elements.update(it, int(weight))
		return
	}

	it = &Element{
		item:     item,
		priority: int(weight),
	}
	heap.Push(&mw.elements, it)
	mw.elementMap[item] = it
}

type Element struct {
	item     any
	priority int
	index    int
}

type minHeap []*Element

func (mh minHeap) Len() int { return len(mh) }

func (mh minHeap) Less(i, j int) bool {
	return mh[i].priority > mh[j].priority
}

func (mh minHeap) Swap(i, j int) {
	mh[i], mh[j] = mh[j], mh[i]
	mh[i].index = i
	mh[j].index = j
}

// return any in order to follow heap.Interface
func (mh *minHeap) Push(x any) {
	n := len(*mh)
	item := x.(*Element)
	item.index = n
	*mh = append(*mh, item)
}

// return any in order to follow heap.Interface
func (mh *minHeap) Pop() any {
	old := *mh
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*mh = old[0 : n-1]
	return item
}

func (mh *minHeap) peak() *Element {
	return (*mh)[len(*mh)-1]
}

func (mh *minHeap) update(ele *Element, priority int) {
	ele.priority = priority
	heap.Fix(mh, ele.index)
}
