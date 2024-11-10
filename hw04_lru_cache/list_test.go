package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require" //nolint:all
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("front", func(t *testing.T) {
		l := NewList()
		it1 := l.PushBack(1)

		require.Nil(t, it1.Prev)
		require.Nil(t, it1.Next)

		it2 := l.PushBack(2)

		require.Nil(t, it1.Prev)
		require.Equal(t, it2, it1.Next)

		require.Equal(t, it1, it2.Prev)
		require.Nil(t, it2.Next)

		require.Equal(t, 2, l.Len())

		f1 := l.Front()
		require.Equal(t, it1, f1)

		l.Remove(f1)
		f1 = l.Front()
		require.Equal(t, 1, l.Len())
		require.Nil(t, f1.Prev)
		require.Nil(t, f1.Next)

		l.Remove(f1)
		require.Equal(t, 0, l.Len())
		f1 = l.Front()
		require.Nil(t, f1)

		f1 = l.Back()
		require.Nil(t, f1)
	})

	t.Run("back", func(t *testing.T) {
		l := NewList()
		it1 := l.PushFront(1)

		require.Nil(t, it1.Prev)
		require.Nil(t, it1.Next)

		it2 := l.PushFront(2)

		require.Nil(t, it2.Prev)
		require.Equal(t, it2, it1.Prev)

		require.Equal(t, it1, it2.Next)
		require.Nil(t, it2.Prev)

		require.Equal(t, 2, l.Len())

		f1 := l.Back()
		require.Equal(t, it1, f1)

		l.Remove(f1)
		f1 = l.Back()
		require.Equal(t, 1, l.Len())
		require.Nil(t, f1.Next)
		require.Nil(t, f1.Next)

		l.Remove(f1)
		require.Equal(t, 0, l.Len())

		f1 = l.Back()
		require.Nil(t, f1)

		f1 = l.Front()
		require.Nil(t, f1)
	})
}
