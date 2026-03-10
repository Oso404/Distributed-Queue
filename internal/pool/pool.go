package pool

import (
	"fmt"

	Q "github.com/Oso404/distributed-queue/internal/queue"
	W "github.com/Oso404/distributed-queue/internal/worker"
	// w2 "github.com/Oso404/distributed-queue/internal/worker"
)

/*
default pool of workers will have 4 workers
*/

const Default_Pool_Size = 4

type Pool struct {
	Workers []*W.Worker
}

// ceate pool with 4 workers and return pointer to pool struct
func Create_Pool() *Pool {
	//create pool struct and initialize with 4 workers
	pool := &Pool{
		Workers: make([]*W.Worker, 4),
	}
	//creation of workers and add to pool
	for i := 0; i < Default_Pool_Size; i++ {
		pool.Workers[i] = W.Create_Worker()
	}

	return pool
}

// (pool *Pool) serves as "this" keyword to reference current pool struct
func (pool *Pool) Start(queue *Q.Queue) {
	//start each worker in pool
	for _, worker := range pool.Workers {
		go worker.Start(queue) //issue: worker start has inf loop and will constantly ask queue for jobs
	}
}

func (pool *Pool) Show_Workers() {
	for _, worker := range pool.Workers {
		fmt.Printf("Worker ID: %s, Current Job ID: %s\n", worker.Worker_ID, worker.Current_Job_ID)
	}
}

func Add_Worker(pool *Pool) {
	//create new worker and add to pool
	newWorker := W.Create_Worker()
	pool.Workers = append(pool.Workers, newWorker)
}

func Remove_Worker(pool *Pool) {
	//remove last worker from pool
	if len(pool.Workers) > 0 {
		pool.Workers = pool.Workers[:len(pool.Workers)-1]
	}
}

/*
i need to find a way to dynamically add and remove workers from pool
i should add a worker for every 5 jobs in the queue and remove a worker for every 5 jobs completed
i need a scheduler that checks the length of the queue and the number of completed jobs every 10 seconds and adds or removes workers accordingly
*/
