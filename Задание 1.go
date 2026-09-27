package main

import (
	"fmt"
	"sync"
	"time"
)

func showSequence(group *sync.WaitGroup) {
	defer group.Done()

	for number := 1; number < 6; number++ {
		fmt.Printf("Число: %d\n", number)
		time.Sleep(time.Second)
	}
}

func main() {
	var completion sync.WaitGroup
	completion.Add(1)

	go showSequence(&completion)
	completion.Wait()

	fmt.Println("Вывод чисел завершён")
}
