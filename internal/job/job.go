package job

import (
	"time"

	"github.com/google/uuid"
)

const defaultTries = 0

type Job struct {
	ID                 string
	Status             string
	Payload            []byte
	VisibilityDeadline time.Time
	StartTime          time.Time
	Retries            int
}

func Create_job(incoming_data []byte) *Job {
	genereated_job := &Job{
		ID:                 uuid.New().String(),
		Status:             "pending",
		Payload:            incoming_data,
		VisibilityDeadline: time.Time{},
		StartTime:          time.Time{},
		Retries:            defaultTries,
	}
	return genereated_job
}
