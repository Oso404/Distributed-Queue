package worker

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	internal "github.com/Oso404/distributed-queue/internal/queue"
	"github.com/google/uuid"
)

type Worker struct {
	Worker_ID      string
	Current_Job_ID string
}

func Create_Worker() *Worker {
	return &Worker{
		Worker_ID: uuid.New().String(),
	}
}

func (worker *Worker) Start(queue *internal.Queue) {
	/*
		worker should ask queue if there is an available job
		queue will check if its length is > 0
		if queue length > 0 worker should aquire next available job
		dequeue returns *Job! (dequeue handles queue fields)
		worker now has current job (Current_Job_ID)
		at the time that job is dequeued its 2 fields (startTime and visibilityDeadline are set)
		*need to find a way to check visiblityDeadline and that we are not over
		generate random value between 0 and 35; "time to perform job"
		if random value is > 30 we will
		1. set job status to failed
		2. set currentjobid to ""
		3. send back to queue
		if random value is < 30 we will
		1. set job status to successful
		2. set currentjobi d to ""
		3. send back to queue
		once job is in back with queue queue will...
		remove from job from Processing Jobs map
		check job status
		1. if successful delete job
		2. if failed
		2a. update retries
		2b. check if retries > max_retries
		2b1. if retries > max_retries send to deadletter queue
		2b2. if tries < max_tries send back of line in queue and send to PendingQueue



	*/
	for {
		j := queue.Dequeue() //j is pointer to job
		if j == nil {
			//here no job available let worker wait 1 second and ask again
			time.Sleep(1 * time.Second)
			// fmt.Println("No job available for", worker.Worker_ID)
			// fmt.Println(worker.Worker_ID, "resting for one second...")
			continue
		}
		fmt.Printf("Worker %s got job %s\n", worker.Worker_ID, j.ID)
		//cant perform actual job rn so will simulate with random number
		jobDuration := time.Duration(rand.Intn(15)) * time.Second
		fmt.Println(worker.Worker_ID, "works for", jobDuration)

		// time.Sleep(jobDuration)

		//visibility deadline set in dequeue() so we can use that to check if job is still alive
		ctx_with_deadline, cancel := context.WithDeadline(context.Background(), j.VisibilityDeadline)
		defer cancel()
		select {
		//will we reach deadline before we finish job?
		case <-ctx_with_deadline.Done():
			fmt.Println("Job", j.ID, "has exceeded visibility deadline and is considered failed.")
			j.Status = "failed"
			//willl job complete before deadline?
		case <-time.After(jobDuration):
			fmt.Println("Worker", worker.Worker_ID, "finished job", j.ID)
			j.Status = "completed"
		}
		//rather than if statement i want to check realtime
		//if statement job can finish and mark as incomplete
		//i have job.VisibilityDeadline

		// if jobDuration <= 10*time.Second {
		// 	j.Status = "completed"
		// } else {
		// 	j.Status = "failed"
		// }

		//send job back to queue
		queue.HandleJobCompletion(j, worker.Worker_ID)
	}
}
