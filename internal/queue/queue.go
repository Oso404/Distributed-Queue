package internal

import "sync"

type Queue struct {
	Name           string
	Jobs           map[string]*Job
	PendingQueue   []string
	ProcessingJobs map[string]*Job
	DeadLetterJobs map[string]*Job
	Mutex          sync.Mutex
}

func Create_Queue(name string) *Queue {
	queue := &Queue{
		Name:           name,
		Jobs:           make(map[string]*Job),
		PendingQueue:   make([]string, 0),
		ProcessingJobs: make(map[string]*Job),
		DeadLetterJobs: make(map[string]*Job),
		Mutex:          sync.Mutex{},
	}
	return queue
}

/*
Issues
What if 2 workers try to get same job at the same time?
What if a worker takes too long doing job?
What if a worker crashes? (What happens to job?)
*/
