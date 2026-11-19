package main

import (
	"fmt"
	"sync"
)

type taskQueue struct {
	mutex  sync.Mutex
	ready  *sync.Cond
	tasks  []int
	closed bool
}

func newTaskQueue() *taskQueue {
	queue := &taskQueue{}
	queue.ready = sync.NewCond(&queue.mutex)
	return queue
}

func (queue *taskQueue) push(task int) bool {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	if queue.closed {
		return false
	}

	queue.tasks = append(queue.tasks, task)
	queue.ready.Signal()
	return true
}

func (queue *taskQueue) pop() (int, bool) {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	for len(queue.tasks) == 0 && !queue.closed {
		queue.ready.Wait()
	}
	if len(queue.tasks) == 0 {
		return 0, false
	}

	task := queue.tasks[0]
	queue.tasks = queue.tasks[1:]
	return task, true
}

func (queue *taskQueue) close() {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	queue.closed = true
	queue.ready.Broadcast()
}

func handleTasks(worker int, queue *taskQueue, group *sync.WaitGroup) {
	defer group.Done()

	for {
		task, ok := queue.pop()
		if !ok {
			return
		}
		fmt.Printf("Обработчик %d выполнил задачу %d\n", worker, task)
	}
}

func main() {
	const taskCount = 10
	const producerCount = 2
	const workerCount = 3

	queue := newTaskQueue()
	var workers sync.WaitGroup
	var producers sync.WaitGroup

	workers.Add(workerCount)
	for worker := 1; worker <= workerCount; worker++ {
		go handleTasks(worker, queue, &workers)
	}

	producers.Add(producerCount)
	for producer := 1; producer <= producerCount; producer++ {
		go func(firstTask int) {
			defer producers.Done()

			for task := firstTask; task <= taskCount; task += producerCount {
				queue.push(task)
			}
		}(producer)
	}

	producers.Wait()
	queue.close()
	workers.Wait()
	fmt.Println("Все задачи обработаны")
}
