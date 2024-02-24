package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func main() {
	apiC := make(chan string, 1)
	const numWorkers = 70

	// goroutine that simulates the worker processor
	go func() {
		for workerId := 0; workerId < numWorkers; workerId++ {

			// create  worker goroutine
			go func() {
				for {
					msg := <-apiC
					response, err := http.Post("http://127.0.0.1:8080", "plain/text", strings.NewReader(msg))
					if err != nil {
						fmt.Println("Error:", err)
						return
					}
					response.Body.Close()
				}
			}()
		}
	}()

	// simulate sending message
	i := 0
	cntMessagesSent := 0
	startTime := time.Now()
	for {

		msg := fmt.Sprint("Hello ", i)
		apiC <- msg
		i++
		cntMessagesSent++
		if time.Since(startTime).Seconds() >= 1 {
			fmt.Println("Messages sent: ", cntMessagesSent)
			startTime = time.Now()
			cntMessagesSent = 0
		}
	}
}
