package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require" //nolint:all
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		// Write me

		c := NewCache(5)

		for i := range 5 {
			k := strconv.Itoa(i)
			existed := c.Set(Key(k), i)
			require.False(t, existed)
		}

		// Size of cache is 5, we inserted 5 elements. 1st element (oldest) must be in cache.
		// Read it

		elem, ok := c.Get("0")
		require.True(t, ok)
		require.Equal(t, 0, elem.(int))

		// Put one more element and check that oldest (1) was pushed out.
		ok = c.Set("5", 5)
		require.False(t, ok)

		elem, ok = c.Get("1")
		require.False(t, ok)
		require.Nil(t, elem)

		// Queue now should be 5 0 4 3 2, next insert should delete 2, but not 3
		ok = c.Set("6", 6)
		require.False(t, ok)

		// Queue now should be 6 5 0 4 3
		elem, ok = c.Get("2")
		require.False(t, ok)
		require.Nil(t, elem)

		elem, ok = c.Get("3")
		require.True(t, ok)
		require.NotNil(t, elem)
		require.Equal(t, 3, elem.(int))

		// Queue now should be 3 6 5 0 4
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
