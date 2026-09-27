package main

import (
	"fmt"
	"sync"
)

func squareWorker(workerNumber int, taskQueue <-chan int, output chan<- int, group *sync.WaitGroup) {
	defer group.Done()

	for value := range taskQueue {
		square := value * value
		fmt.Printf("Обработчик %d вычислил квадрат числа %d\n", workerNumber, value)
		output <- square
	}
}

func main() {
	taskQueue := make(chan int, 10)
	output := make(chan int, 10)

	var workers sync.WaitGroup
	for number := 1; number <= 3; number++ {
		workers.Add(1)
		go squareWorker(number, taskQueue, output, &workers)
	}

	for value := 1; value <= 10; value++ {
		taskQueue <- value
	}
	close(taskQueue)

	go func() {
		workers.Wait()
		close(output)
	}()

	fmt.Println("Полученные квадраты:")
	for square := range output {
		fmt.Println(square)
	}
}
