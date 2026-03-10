package pool

import (
	"fmt"
	"time"

	Q "github.com/Oso404/distributed-queue/internal/queue"
	internal "github.com/Oso404/distributed-queue/internal/queue"
	W "github.com/Oso404/distributed-queue/internal/worker"
	// W "github.com/osos"
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
		Workers: make([]*W.Worker, Default_Pool_Size),
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
	go pool.Monitor_Queue(queue, 5, 20) //monitor queue and adjust number of workers between 1 and 10
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

func (pool *Pool) Monitor_Queue(queue *Q.Queue, minWorkers int, maxWorkers int) {
	//ill have pool check every .25 seconds the length of the queue and adjust number of workers accordingly
	//ill add a worker for every 5 jobs in the queue after 5 jobs in the queue
	/*
		im going to have this function run in an goroutine when the pool starts
	*/
	for {
		queueLength := len(queue.PendingQueue)
		if queueLength > 5 && len(pool.Workers) < maxWorkers {
			Add_Worker(pool)
			fmt.Println("Added worker. Total workers:", len(pool.Workers))
		} else if queueLength < 5 && len(pool.Workers) > minWorkers {
			Remove_Worker(pool)
			fmt.Println("Removed worker. Total workers:", len(pool.Workers))
		}
	}
}

func (p *Pool) MonitorQueue_(queue *internal.Queue, minWorkers, maxWorkers int) {
	ticker := time.NewTicker(500 * time.Millisecond) // check every 0.5s
	defer ticker.Stop()

	for range ticker.C {
		queueLength := len(queue.PendingQueue)
		currentWorkers := len(p.Workers)

		// Determine how many workers to add or remove
		// 1 worker per 5 pending jobs (scale gradually)
		desiredWorkers := queueLength/5 + 1

		if desiredWorkers > maxWorkers {
			desiredWorkers = maxWorkers
		}
		if desiredWorkers < minWorkers {
			desiredWorkers = minWorkers
		}

		if desiredWorkers > currentWorkers {
			// Add the difference
			toAdd := desiredWorkers - currentWorkers
			for i := 0; i < toAdd; i++ {
				newWorker := W.Create_Worker()
				p.Workers = append(p.Workers, newWorker)
				go newWorker.Start(queue)
			}
			fmt.Println("Added workers. Total:", len(p.Workers))
		} else if desiredWorkers < currentWorkers {
			// Remove idle workers
			toRemove := currentWorkers - desiredWorkers
			removed := 0
			for i := 0; i < len(p.Workers) && removed < toRemove; i++ {
				w := p.Workers[i]
				if w.Current_Job_ID == "" { // idle
					w.Stop() // implement Stop() to exit Start() loop
					// remove from slice
					p.Workers = append(p.Workers[:i], p.Workers[i+1:]...)
					i-- // adjust index after removing
					removed++
				}
			}
			if removed > 0 {
				fmt.Println("Removed idle workers. Total:", len(p.Workers))
			}
		}
	}
}

/*
i need to find a way to dynamically add and remove workers from pool
i should add a worker for every 5 jobs in the queue and remove a worker for every 5 jobs completed
i need a scheduler that checks the length of the queue and the number of completed jobs every 10 seconds and adds or removes workers accordingly
*/
