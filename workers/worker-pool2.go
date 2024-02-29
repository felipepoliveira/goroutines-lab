package workers

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// The eagerWorker is used to control a goroutine process that runs over a for {} eternal loop
// that receives data from a channel. Every time the channel receive a data it will trigger the job
// provided in the EagerWorkerPool function
type eagerWorker[R any] struct {
	// flags if the worker is busy doing a job
	busy atomic.Bool
	// channel used to trigger the job inside the for{} goroutine
	receiverC chan R
}

// Represents a collection of eagerWorkers. The eagerWorkerPool will store the state of each
// goroutine that is triggered by a channel feeding (eagerWorker.receiverC). The pool will select
// available workers that will be used to execute the job.
// It is possible that all the workers from the worker pool is on a busy state. When this happen the pool
// will trigger the 'unbalanced' mode where it will buffer excessive data into each worker using a round-robin
// strategy, so all packages is sent successfully without errors. When a unbalanced mode is excessively used is recommended
// to increase the amount of workers in the eager pool
type eagerWorkerPool[R any] struct {
	// store each goroutine worker state
	workers []*eagerWorker[R]
	// when the unbalanced mode is activate this property will mark the last worker that received a buffered data in its channel
	indexToUseInUnbalancedMode int
	//An mutex lock that is locked whenever the EagerWorkerPool.Do() function is triggered. This will prevent data race between calls
	doingWorkLock sync.Mutex
}

// Execute the job of the eagerWorkerPool when the given data.
// The pool will select
// available workers that will be used to execute the job.
// It is possible that all the workers from the worker pool is on a busy state. When this happen the pool
// will trigger the 'unbalanced' mode where it will buffer excessive data into each worker using a round-robin
// strategy, so all packages is sent successfully without errors. When a unbalanced mode is excessively used is recommended
// to increase the amount of workers in the eager pool
func (wp *eagerWorkerPool[R]) Do(dat R) {
	// lock all race process until it finishes
	wp.doingWorkLock.Lock()
	defer wp.doingWorkLock.Unlock()

	// try to find a not busy worker to process the data
	for idx, worker := range wp.workers {
		if !worker.busy.Load() {
			fmt.Printf("Work %d is available for work\n", idx)
			worker.receiverC <- dat
			return
		}
	}

	fmt.Printf("All workers area busy. Using unbalanced mode on worker%d\n", wp.indexToUseInUnbalancedMode)

	// execute the 'unbalanced' mode
	// the unbalanced mode will buffer up the worker channels because all the channels are busy
	wp.workers[wp.indexToUseInUnbalancedMode].receiverC <- dat

	// go to the next index (reset to 0 if it will surpass the maximum amount of workers)
	wp.indexToUseInUnbalancedMode = (wp.indexToUseInUnbalancedMode + 1) % len(wp.workers)

}

// Error type used
type workerError struct {
	reason string
}

// Implement error interface
func (e *workerError) Error() string {
	return e.reason
}

// Create a EagerWorkerPool
//
// Parameters:
//
//	workers: The number of workers that will be created. Should be >= 1
//	chBufSize: The buffer size used in each worker channel
//	job: The job that will be executed every time eagerWorkerPool.Do() function is triggered
//
// Returns:
//
//	A eagerWorkerPool struct that will be used to execute a function in a goroutine or a error if it fails
func EagerWorkerPool[R any](workers int, chBufSize int, job func(dat R)) (*eagerWorkerPool[R], error) {
	// check for valid number of workers
	if workers < 1 {
		return nil, &workerError{"The number of workers should be >= 1"}
	}
	if chBufSize < 1 {
		return nil, &workerError{"The channel buffer size should be >= 1"}
	}

	// create the object that will contain the worker group and initialize its workers
	wp := eagerWorkerPool[R]{workers: make([]*eagerWorker[R], workers)}
	for wid := 0; wid < workers; wid++ {
		worker := eagerWorker[R]{receiverC: make(chan R, chBufSize)}

		// start a always-running goroutine that will execute the given job
		go func() {
			for {
				dat := <-worker.receiverC
				worker.busy.Store(true)
				job(dat)
				worker.busy.Store(false)
			}
		}()

		wp.workers[wid] = &worker
	}

	return &wp, nil
}
