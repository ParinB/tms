package tms

import (
	"encoding/json"
	"time"
)

const (
	inputTimeFormat = "2006-01-02 15:04:05"
)

//CreateUserRequest contains a create user request
type CreateUserRequest struct {
	FirstName string `validate:"required,alpha,gt=1"`
	LastName  string `validate:"required,alpha,gt=1"`
	Email     string `validate:"required,email"`
}

//CreateTaskRequest contains a create task request
type CreateTaskRequest struct {
	CreatedAt      RequestTime
	Description    string `validate:"required,alpha,gt=1"`
	StartTime      RequestTime
	EndTime        RequestTime
	ReminderPeriod RequestTime
}

//RequestTime represents a request time
type RequestTime struct {
	t time.Time
}

//UnmarshalJSON decodes json
func (r *RequestTime) UnmarshalJSON(data []byte) error {
	var input string
	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}
	time, err := time.Parse(inputTimeFormat, input)
	if err != nil {
		return err
	}
	r.t = time
	return nil
}
