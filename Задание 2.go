package main

import (
	"fmt"
	"sync"
)

type Cache struct {
	mutex sync.RWMutex
	data  map[string]string
}

func (cache *Cache) Get(key string) (string, bool) {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	value, found := cache.data[key]
	return value, found
}

func (cache *Cache) Set(key, value string) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	cache.data[key] = value
}

func main() {
	cache := Cache{data: make(map[string]string)}
	cache.Set("страница", "старая версия")

	var group sync.WaitGroup
	for i := 1; i <= 5; i++ {
		group.Add(1)
		go func(id int) {
			defer group.Done()

			if value, found := cache.Get("страница"); found {
				fmt.Printf("Читатель %d: %s\n", id, value)
			}
		}(i)
	}

	group.Add(1)
	go func() {
		defer group.Done()
		cache.Set("страница", "новая версия")
	}()

	group.Wait()

	value, _ := cache.Get("страница")
	fmt.Println("В кэше после обновления:", value)
}
