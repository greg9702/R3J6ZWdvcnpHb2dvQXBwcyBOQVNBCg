package fetcher

import (
	"fmt"
	"sync"
)

const (
	maxCoLimit = 2
)

type ExecutorInterface interface {
	AddNewJob(job JobInterface)
}

type LimitAwareExecutor struct {
	limit int
	jobs  chan JobInterface
}

func NewLimitAwareExecutor() ExecutorInterface {
	executorLimit := maxCoLimit

	// TODO Get limit from config
	le := LimitAwareExecutor{executorLimit, make(chan JobInterface, 100)}

	var initWaitGroup sync.WaitGroup

	initWaitGroup.Add(le.limit)
	for i := 1; i <= le.limit; i++ {
		go le.spawnWorker(i, &initWaitGroup)
	}
	initWaitGroup.Wait()
	return &le
}

func (le *LimitAwareExecutor) AddNewJob(job JobInterface) {
	le.jobs <- job
}

func (le *LimitAwareExecutor) spawnWorker(id int, initWaitGroup *sync.WaitGroup) {
	fmt.Printf("Spawning worker %d\n", id)
	initWaitGroup.Done()
	for job := range le.jobs {
		fmt.Printf("[%d] executing new task\n", id)
		job.Execute()
	}
	fmt.Printf("Worker %d exiting\n", id)
}
