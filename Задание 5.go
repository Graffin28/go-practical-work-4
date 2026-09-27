package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"sync"
)

type digestResult struct {
	filename string
	digest   []byte
	err      error
}

func calculateDigest(filename string) ([]byte, error) {
	openedFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer openedFile.Close()

	calculator := md5.New()
	if _, err = io.Copy(calculator, openedFile); err != nil {
		return nil, err
	}
	return calculator.Sum(nil), nil
}

func digestWorker(filenames <-chan string, completed chan<- digestResult, group *sync.WaitGroup) {
	defer group.Done()

	for filename := range filenames {
		digest, err := calculateDigest(filename)
		completed <- digestResult{filename: filename, digest: digest, err: err}
	}
}

func main() {
	fileList := []string{
		"Задание 1.go", "Задание 2.go", "Задание 3.go", "Задание 4.go",
		"Задание 5.go", "Задание 6.go", "Задание 7.go", "Задание 8.go",
	}

	filenames := make(chan string)
	completed := make(chan digestResult)

	var pool sync.WaitGroup
	for count := 0; count < 3; count++ {
		pool.Add(1)
		go digestWorker(filenames, completed, &pool)
	}

	go func() {
		for _, filename := range fileList {
			filenames <- filename
		}
		close(filenames)
		pool.Wait()
		close(completed)
	}()

	for result := range completed {
		if result.err != nil {
			fmt.Printf("Не удалось обработать %s: %v\n", result.filename, result.err)
			continue
		}
		fmt.Printf("%x  %s\n", result.digest, result.filename)
	}
}
