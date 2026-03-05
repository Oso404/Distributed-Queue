package job

import (
	"time"

	"github.com/google/uuid"
)

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
		Retries:            3,
	}
	return genereated_job
}
