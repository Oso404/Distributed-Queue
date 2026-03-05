package worker

import (
	internal "github.com/Oso404/distributed-queue/internal/queue"
)

type Worker struct {
	ID    int
	Queue *internal.Queue
	Done  chan bool
}
