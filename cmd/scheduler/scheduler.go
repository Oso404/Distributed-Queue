package scheduler

import (
	"fmt"
	"time"

	J "github.com/Oso404/distributed-queue/internal/job"
	pool "github.com/Oso404/distributed-queue/internal/pool"
	Q "github.com/Oso404/distributed-queue/internal/queue"
)

// Queue is the queue that scheduler will manage
// Pool is the pool of workers that scheduler will manage
type Scheduler struct {
	Queue *Q.Queue
	Pool  *pool.Pool
}

func Create_Scheduler(queue *Q.Queue, pool *pool.Pool) *Scheduler {
	return &Scheduler{
		Queue: queue,
		Pool:  pool,
	}
}

func (s *Scheduler) Check_For_Job() (*J.Job, error) {
	//start pool of workers
	/*
		the issue with this is that pool.Start will execute worker.Start
		but worker.Start has the worker constantly asking the queue for jobs
		i need to remove the infinite loop from worker.Start and instead have the scheduler loop ask the pool to
		start each worker and then have the scheduler loop also monitor the queue for job completions and handle them accordingly
	*/
	s.Pool.Start(s.Queue)
	//check if there is a job available in the queue
	for {
		if s.Queue.JobAvailable() {
			fmt.Println("Job is available in the queue.")
		} else {
			fmt.Println("No job available in the queue. Checking again in 1 second...")
		}
		time.Sleep(1 * time.Second)
	}

}

func (s *Scheduler) Review_Availability_Workers() {
	sl := s.Pool.Workers
	for _, worker := range sl {
		if worker.Current_Job_ID == "" {
			fmt.Println("Worker", worker.Worker_ID, "is available.")
		} else {
			fmt.Println("Worker", worker.Worker_ID, "is currently working on job", worker.Current_Job_ID)
		}
	}
}
