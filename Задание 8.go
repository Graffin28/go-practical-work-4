package main

import (
	"fmt"
	"sync"
	"time"
)

func service(instance int, inbox <-chan int, group *sync.WaitGroup) {
	defer group.Done()

	for requestNumber := range inbox {
		fmt.Printf("Экземпляр %d получил запрос %d\n", instance, requestNumber)
		time.Sleep(100 * time.Millisecond)
	}
}

func distribute(incoming <-chan int, instances []chan int) {
	next := 0
	for requestNumber := range incoming {
		instances[next] <- requestNumber
		next = (next + 1) % len(instances)
	}

	for _, inbox := range instances {
		close(inbox)
	}
}

func main() {
	const instanceCount = 3

	incoming := make(chan int)
	instances := make([]chan int, instanceCount)
	var services sync.WaitGroup

	for index := range instances {
		instances[index] = make(chan int, 4)
		services.Add(1)
		go service(index+1, instances[index], &services)
	}

	go distribute(incoming, instances)
	for requestNumber := 1; requestNumber <= 10; requestNumber++ {
		incoming <- requestNumber
	}
	close(incoming)

	services.Wait()
	fmt.Println("Балансировщик завершил работу")
}
