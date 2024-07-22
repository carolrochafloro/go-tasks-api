package tasks

import (
	"errors"
	"time"
)

type Task struct {
	Id          string `bson:"task_id,omitempty"`
	Title       string `bson:"title,omitempty"`
	Deadline    string `bson:"deadline,omitempty"`
	IsCompleted bool   `bson:"is_completed,omitempty"`
}

func (t *Task) Validade() error {

	if len(t.Title) < 3 {
		return errors.New("task title must have at least 3 characters")
	}

	parsedTime, _ := time.Parse(time.DateTime, t.Deadline)

	if parsedTime.Before(time.Now()) {
		return errors.New("the deadline must be set to a time in the future")
	}

	return nil
}