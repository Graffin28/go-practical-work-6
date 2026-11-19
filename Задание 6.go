package main

import (
	"fmt"
	"sync"
)

type Logger struct {
	mutex sync.Mutex
}

func (logger *Logger) Print(id int, message string) {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	fmt.Printf("[Горутина %d] ", id)
	fmt.Printf("%s\n", message)
}

func main() {
	var logger Logger
	var group sync.WaitGroup

	for i := 1; i <= 5; i++ {
		group.Add(1)
		go func(id int) {
			defer group.Done()

			logger.Print(id, "Начало работы")
			logger.Print(id, "Работа завершена")
		}(i)
	}

	group.Wait()
}
