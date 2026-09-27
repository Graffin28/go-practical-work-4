package main

import (
	"fmt"
	"time"
)

func main() {
	limiter := time.NewTicker(200 * time.Millisecond)
	defer limiter.Stop()

	startedAt := time.Now()
	for requestNumber := 1; requestNumber <= 15; requestNumber++ {
		<-limiter.C
		fmt.Printf("Запрос №%d выполнен спустя %v\n", requestNumber, time.Since(startedAt).Round(time.Millisecond))
	}

	fmt.Println("Обработка 15 запросов завершена")
}
