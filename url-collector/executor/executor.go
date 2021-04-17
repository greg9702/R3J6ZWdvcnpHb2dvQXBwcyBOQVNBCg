package executor

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

const (
	maxCoLimit = 5
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
	var err error

	limitOfExecutorFromEnv := os.Getenv("CONCURRENT_REQUESTS")

	if limitOfExecutorFromEnv != "" {
		executorLimit, err = strconv.Atoi(limitOfExecutorFromEnv)
		if err != nil {
			fmt.Println("Unable to parse CONCURRENT_REQUESTS to int value, using default value")
			executorLimit = maxCoLimit
		}
	}

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
		job.Execute()
	}
	fmt.Printf("Worker %d exiting\n", id)
}
