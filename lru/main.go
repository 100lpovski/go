package main

import (
	"container/list"
	"fmt"
)

type entry[K comparable, V any] struct {
	key   K
	value V
}

type LRUCache[K comparable, V any] struct {
	capacity int
	items    map[K]*list.Element
	order    *list.List
}

func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	if capacity <= 0 {
		panic("capacity must be > 0")
	}
	return &LRUCache[K, V]{
		capacity: capacity,
		items:    make(map[K]*list.Element, capacity),
		order:    list.New(),
	}
}

func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	if el, ok := c.items[key]; ok {
		c.order.MoveToFront(el)
		return el.Value.(*entry[K, V]).value, true
	}
	var zero V
	return zero, false
}

func (c *LRUCache[K, V]) Put(key K, value V) {
	if el, ok := c.items[key]; ok {
		el.Value.(*entry[K, V]).value = value
		c.order.MoveToFront(el)
		return
	}

	el := c.order.PushFront(&entry[K, V]{key: key, value: value})
	c.items[key] = el

	if c.order.Len() > c.capacity {
		c.evict()
	}
}

func (c *LRUCache[K, V]) evict() {
	back := c.order.Back()
	if back == nil {
		return
	}
	c.order.Remove(back)
	kv := back.Value.(*entry[K, V])
	delete(c.items, kv.key)
}

func (c *LRUCache[K, V]) Len() int {
	return c.order.Len()
}

func main() {
	cache := NewLRUCache[string, int](3)

	cache.Put("a", 1)
	cache.Put("b", 2)
	cache.Put("c", 3)

	fmt.Println(cache.Get("a"))

	cache.Put("d", 4)

	if _, ok := cache.Get("b"); !ok {
		fmt.Println("b вытеснен")
	}
	fmt.Println(cache.Get("c"))
	fmt.Println(cache.Get("d"))
}
