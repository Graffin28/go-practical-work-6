package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	const count = 100
	candidates := []string{"Анна", "Борис", "Виктор"}
	votes := make(map[string]int)

	var mutex sync.Mutex
	var group sync.WaitGroup

	for i := 0; i < count; i++ {
		group.Add(1)
		go func() {
			defer group.Done()

			candidate := candidates[rand.Intn(len(candidates))]
			mutex.Lock()
			votes[candidate]++
			mutex.Unlock()
		}()
	}
	group.Wait()

	total := 0
	for _, candidate := range candidates {
		fmt.Printf("%s: %d\n", candidate, votes[candidate])
		total += votes[candidate]
	}
	fmt.Println("Всего голосов:", total)
}
