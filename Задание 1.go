package main

import (
	"fmt"
	"sync"
)

func main() {
	const count = 1000

	var views int
	var mutex sync.Mutex
	var group sync.WaitGroup

	for i := 0; i < count; i++ {
		group.Add(1)
		go func() {
			defer group.Done()

			mutex.Lock()
			views++
			mutex.Unlock()
		}()
	}

	group.Wait()
	fmt.Println("Количество просмотров:", views)
}
