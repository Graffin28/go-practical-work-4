package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type siteReply struct {
	address string
	status  string
	err     error
}

func checkSites(addresses <-chan string, replies chan<- siteReply, client *http.Client, group *sync.WaitGroup) {
	defer group.Done()

	for address := range addresses {
		response, err := client.Get(address)
		if err != nil {
			replies <- siteReply{address: address, err: err}
			continue
		}

		replies <- siteReply{address: address, status: response.Status}
		response.Body.Close()
	}
}

func main() {
	sites := []string{
		"https://go.dev",
		"https://pkg.go.dev",
		"https://github.com",
		"https://example.com",
		"https://httpbin.org/status/404",
	}

	addresses := make(chan string)
	replies := make(chan siteReply)
	client := &http.Client{Timeout: 8 * time.Second}

	var pool sync.WaitGroup
	for count := 0; count < 3; count++ {
		pool.Add(1)
		go checkSites(addresses, replies, client, &pool)
	}

	go func() {
		for _, site := range sites {
			addresses <- site
		}
		close(addresses)
		pool.Wait()
		close(replies)
	}()

	for reply := range replies {
		if reply.err != nil {
			fmt.Printf("%s: ошибка запроса\n", reply.address)
		} else {
			fmt.Printf("%s: %s\n", reply.address, reply.status)
		}
	}
}
