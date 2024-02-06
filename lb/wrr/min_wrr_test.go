package wrr

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinHeap(t *testing.T) {
	t.Run("Push", func(t *testing.T) {
		mh := minHeap{}
		assert.NotPanics(
			t,
			func() {
				for i := 0; i < 10; i++ {
					heap.Push(&mh, &Element{
						item:     i,
						priority: i,
					})
				}
			},
		)
	})

	t.Run("Pop", func(t *testing.T) {
		t.Run("panic", func(t *testing.T) {
			mh := minHeap{}
			assert.Panics(t, func() { heap.Pop(&mh) })
		})

		t.Run("success", func(t *testing.T) {
			mh := minHeap{}
			var out any
			assert.NotPanics(t, func() { mh.Push(&Element{}); out = heap.Pop(&mh) })
			assert.Equal(t, 0, mh.Len())
			assert.IsType(t, &Element{}, out)
			assert.Equal(t, &Element{}, out.(*Element))
		})
	})

	t.Run("peak", func(t *testing.T) {
		t.Run("panic, no element", func(t *testing.T) {
			mh := minHeap{}
			assert.Panics(t, func() { mh.peak() })
		})

		t.Run("success", func(t *testing.T) {
			mh := minHeap{}

			eles := make([]*Element, 10)
			for i := range eles {
				ele := &Element{
					item:     i,
					priority: i,
				}
				eles[i] = ele
				heap.Push(&mh, ele)
			}

			actual := mh.peak()
			assert.Equal(t, 10, mh.Len())
			assert.Equal(t, eles[0], actual)
		})
	})

	t.Run("update", func(t *testing.T) {
		mh := minHeap{}

		eles := make([]*Element, 10)
		for i := range eles {
			ele := &Element{
				item:     i,
				priority: i,
			}
			eles[i] = ele
			heap.Push(&mh, ele)
		}

		mh.update(eles[0], 10)
		assert.Equal(t, eles[1], mh.peak())
	})
}
