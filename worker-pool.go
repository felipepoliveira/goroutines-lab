package main

import (
	"fmt"
	"net/http"
	"time"
)

type MessageProcessorPayload struct {
	message   string
	serverUrl string
}

func main2() {
	apiC := make(chan MessageProcessorPayload, 1)
	const numWorkers = 400

	// goroutine that simulates the worker processor
	for workerId := 0; workerId < numWorkers; workerId++ {

		// create  worker goroutine
		go func() {
			cntMessagesSent := 0
			// startTime := time.Now()
			for {
				msg := <-apiC
				response, err := http.Get(msg.serverUrl)
				time.Sleep(time.Second)

				// if time.Since(startTime).Seconds() >= 10 {
				// 	fmt.Printf("Worker %d sent %d msgs in 10 seconds\n", workerId, cntMessagesSent)
				// 	startTime = time.Now()
				// 	cntMessagesSent = 0
				// }

				if err != nil {
					fmt.Println("Error:", err)
					continue
				}
				cntMessagesSent++
				response.Body.Close()
			}
		}()
	}

	// simulate sending message
	i := 0
	serversUrls := []string{
		"http://localhost:8080",
		// "http://127.0.0.1:8081",
		// "http://127.0.0.1:8082",
		// "http://127.0.0.1:8083",
		// "http://127.0.0.1:8084",
		// "http://127.0.0.1:8085",
		// "http://127.0.0.1:8086",
		// "http://127.0.0.1:8087",
		// "http://127.0.0.1:8088",
		// "http://127.0.0.1:8089",
	}
	cntMessagesSent := 0
	startTime := time.Now()
	serverIndex := 0
	for {
		msg := fmt.Sprint("Hello ", i)
		apiC <- MessageProcessorPayload{message: msg, serverUrl: serversUrls[serverIndex]}
		i++
		cntMessagesSent++
		serverIndex++
		if serverIndex >= len(serversUrls) {
			serverIndex = 0
		}
		if time.Since(startTime).Seconds() >= 1 {
			fmt.Println("Messages sent: ", cntMessagesSent)
			startTime = time.Now()
			cntMessagesSent = 0
		}

		time.Sleep(time.Millisecond * 10)

	}
}
