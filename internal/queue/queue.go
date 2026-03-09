package internal

import (
	"fmt"
	"sync"
	"time"

	job "github.com/Oso404/distributed-queue/internal/job"
)

type Queue struct {
	Name           string              //name of of Queue; can be anything
	Jobs           map[string]*job.Job //list of all jobs that have been in the queue
	PendingQueue   []string            //list of job ids that are in the queue
	ProcessingJobs map[string]*job.Job //list of jobs that are currently being processed by workers
	DeadLetterJobs map[string]*job.Job //list of jobs that couldnt be processed even after maxRetries
	MaxRetries     int                 //default is 3
	Mutex          sync.Mutex          //prevent race conditions from workers
}

func Create_Queue(name string) *Queue {
	queue := &Queue{
		Name:           name,
		Jobs:           make(map[string]*job.Job),
		PendingQueue:   make([]string, 0),
		ProcessingJobs: make(map[string]*job.Job),
		DeadLetterJobs: make(map[string]*job.Job),
		MaxRetries:     3,
		Mutex:          sync.Mutex{},
	}
	return queue
}

// Enqueue
// (q *Queue) serves as "this" keyword in other languages
func (q *Queue) Enqueue(job *job.Job) {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()
	q.Jobs[job.ID] = job                            //add newly created Job to Jobs map
	q.PendingQueue = append(q.PendingQueue, job.ID) //add jobID to PendignQueue
	job.Status = "Pending"                          //change status
}

func (q *Queue) Dequeue() *job.Job {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	//check if PendingQueue has no jobs
	if len(q.PendingQueue) == 0 {
		return nil
	}
	//PendingQueue is []string IDs; retrieve next available job
	jobID := q.PendingQueue[0]
	//remove job that is going to be taken by worker
	q.PendingQueue = q.PendingQueue[1:]
	//retrieve entire *Job based on jobID
	job := q.Jobs[jobID]
	//change status to processing and then add to ProcessingJobs map
	job.Status = "processing"
	q.ProcessingJobs[jobID] = job
	//Initialize job startTime and VisibilityDeadline
	//VisibilityDeadline important helps us keep track if job is alive! (Maybe worker crashed which is why its taking forever; we have no way of finding out)
	job.StartTime = time.Now()
	job.VisibilityDeadline = time.Now().Add(10 * time.Second)
	fmt.Println("Job Dequeued at", job.StartTime, "and has deadline at", job.VisibilityDeadline)

	return job
}

// HandleJobCompletion will run when Worker finishes job (or after 30 seconds)
func (q *Queue) HandleJobCompletion(j *job.Job, workerID string) {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()
	delete(q.ProcessingJobs, j.ID)
	if j.Status == "completed" {
		// remove job completely
		delete(q.Jobs, j.ID)
		fmt.Printf("Worker %s completed job %s\n", workerID, j.ID)
	} else {
		// failed job
		fmt.Println(j.ID, "has failed!")
		j.Retries++
		if j.Retries > q.MaxRetries {
			j.Status = "deadletter"
			q.DeadLetterJobs[j.ID] = j
			fmt.Printf("Job %s failed too many times → DeadLetter\n", j.ID)
		} else {
			j.Status = "pending"
			q.PendingQueue = append(q.PendingQueue, j.ID)
			fmt.Printf("Job %s failed by worker %s → retry %d\n", j.ID, workerID, j.Retries)
		}
	}
}

/*
Issues
What if 2 workers try to get same job at the same time? -> mutex
What if a worker takes too long doing job? -> visiblity deadline
What if a worker crashes? (What happens to job?)-> visibility deadline
*/
