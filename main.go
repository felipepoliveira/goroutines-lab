package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/felipepoliveira/goroutines-lab/workers"
)

func main() {

	const messagesPerSecond int = 1
	const msgProcTimeoutInSecs = 10

	wp, err := workers.EagerWorkerPool[string](messagesPerSecond*msgProcTimeoutInSecs, 3, func(dat string) {
		fmt.Println(dat)
		time.Sleep(time.Second * time.Duration(rand.Intn(msgProcTimeoutInSecs)+1))
	})

	if err != nil {
		fmt.Println(err)
	}

	for {
		wp.Do("Working...")
		time.Sleep(time.Millisecond * 1000)
	}
}
