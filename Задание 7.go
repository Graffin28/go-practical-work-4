package main

import (
	"fmt"
	"sync"
)

type stateCommand struct {
	operation string
	key       string
	value     int
	answer    chan int
}

func runStateManager(commands <-chan stateCommand) {
	storage := make(map[string]int)

	for command := range commands {
		switch command.operation {
		case "write":
			storage[command.key] = command.value
		case "increase":
			storage[command.key]++
		case "read":
			command.answer <- storage[command.key]
		}
	}
}

func main() {
	commands := make(chan stateCommand)
	go runStateManager(commands)

	commands <- stateCommand{operation: "write", key: "users", value: 10}

	var senders sync.WaitGroup
	for count := 0; count < 100; count++ {
		senders.Add(1)
		go func() {
			defer senders.Done()
			commands <- stateCommand{operation: "increase", key: "visits"}
		}()
	}
	senders.Wait()

	answer := make(chan int)
	for _, key := range []string{"users", "visits"} {
		commands <- stateCommand{operation: "read", key: key, answer: answer}
		fmt.Printf("%s = %d\n", key, <-answer)
	}

	close(commands)
}
