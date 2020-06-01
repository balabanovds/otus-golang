package hw04_lru_cache //nolint:golint,stylecheck

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		t.Run("purge first element", func(t *testing.T) {
			c := initCache(t)

			was := c.Set("d", 4)
			assert.False(t, was)

			a, ok := c.Get("a")
			assert.False(t, ok)
			assert.Nil(t, a)
		})

		t.Run("purge less used element", func(t *testing.T) {
			c := initCache(t)

			_ = c.Set("a", 4)
			_, _ = c.Get("c")

			was := c.Set("d", 4)
			assert.False(t, was)

			b, ok := c.Get("b")
			assert.False(t, ok)
			assert.Nil(t, b)
		})
	})
}

func TestCacheMultithreading(t *testing.T) {
	// t.Skip() // Remove if task with asterisk completed

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(strconv.Itoa(i), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(strconv.Itoa(rand.Intn(1_000_000)))
		}
	}()

	wg.Wait()
}

func initCache(t *testing.T) Cache {
	t.Helper()

	c := NewCache(3)

	was := c.Set("a", 1)
	assert.False(t, was)

	was = c.Set("b", 2)
	assert.False(t, was)

	was = c.Set("c", 2)
	assert.False(t, was)

	return c
}
