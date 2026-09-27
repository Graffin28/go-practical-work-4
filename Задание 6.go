package main

import (
	"fmt"
	"time"
)

type lookupAnswer struct {
	provider string
	value    string
	found    bool
}

func querySource(provider string, pause time.Duration, found bool, answers chan<- lookupAnswer) {
	time.Sleep(pause)

	answer := lookupAnswer{provider: provider, found: found}
	if found {
		answer.value = "Данные получены из источника: " + provider
	}
	answers <- answer
}

func main() {
	answers := make(chan lookupAnswer, 3)

	go querySource("локальный кэш", 120*time.Millisecond, false, answers)
	go querySource("основная база", 350*time.Millisecond, true, answers)
	go querySource("резервный API", 220*time.Millisecond, true, answers)

	for attempt := 0; attempt < 3; attempt++ {
		answer := <-answers
		if answer.found {
			fmt.Println("Первый успешный ответ:", answer.value)
			return
		}
		fmt.Println(answer.provider, "не вернул данные")
	}

	fmt.Println("Поиск не дал результатов")
}
