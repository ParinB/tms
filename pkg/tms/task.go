package tms

import "time"

//Task contains information about a task
type Task struct {
	ID             uint
	CreatedAt      time.Time
	Description    string
	StartTime      time.Time
	EndTime        time.Time
	Duration       time.Duration
	ReminderPeriod time.Time
}
