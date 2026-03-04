package internal

import (
	"time"

	"github.com/google/uuid"
)

type Job struct {
	ID                 string
	Status             string
	Payload            []byte
	VisibilityDeadline time.Time
}

func Create_job(incoming_data []byte) *Job {
	genereated_job := &Job{
		ID:                 uuid.New().String(),
		Status:             "pending",
		Payload:            incoming_data,
		VisibilityDeadline: time.Now().Add(30 * time.Second),
	}
	return genereated_job
}
