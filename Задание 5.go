package main

import (
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	value     string
	expiresAt time.Time
}

type cache struct {
	mu            sync.Mutex
	entries       map[string]cacheEntry
	ttl           time.Duration
	refreshOnRead bool
}

func newCache(ttl time.Duration, refreshOnRead bool) *cache {
	return &cache{
		entries:       make(map[string]cacheEntry),
		ttl:           ttl,
		refreshOnRead: refreshOnRead,
	}
}

func (c *cache) set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = cacheEntry{value: value, expiresAt: time.Now().Add(c.ttl)}
}

func (c *cache) get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.entries[key]
	if !ok {
		return "", false
	}

	now := time.Now()
	if !now.Before(entry.expiresAt) {
		delete(c.entries, key)
		return "", false
	}

	if c.refreshOnRead {
		entry.expiresAt = now.Add(c.ttl)
		c.entries[key] = entry
	}
	return entry.value, true
}

func (c *cache) cleanExpired(stop <-chan struct{}) {
	ticker := time.NewTicker(c.ttl)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			c.mu.Lock()
			now := time.Now()
			for key, entry := range c.entries {
				if !now.Before(entry.expiresAt) {
					delete(c.entries, key)
				}
			}
			c.mu.Unlock()
		}
	}
}

func main() {
	ttl := 2 * time.Second
	fixed := newCache(ttl, false)
	sliding := newCache(ttl, true)

	stop := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	for _, c := range []*cache{fixed, sliding} {
		go func(c *cache) {
			defer wg.Done()
			c.cleanExpired(stop)
		}(c)
	}

	fixed.set("page", "Главная страница")
	sliding.set("page", "Главная страница")

	time.Sleep(time.Second)
	value, _ := fixed.get("page")
	fmt.Println("Первое чтение без продления:", value)
	value, _ = sliding.get("page")
	fmt.Println("Первое чтение с продлением:", value)

	time.Sleep(1500 * time.Millisecond)
	_, ok := fixed.get("page")
	fmt.Println("Без продления запись ещё доступна:", ok)
	_, ok = sliding.get("page")
	fmt.Println("С продлением запись ещё доступна:", ok)

	close(stop)
	wg.Wait()
}
