package tms

import (
	"context"
	"time"
)

//Storer is an interface providing access to records
type Storer interface {
	Close() error
	CreateUser(ctx context.Context, firstName, lastName, email string) error
	CreateTask(ctx context.Context, description string, userID uint, startTime, endTime, reminderPeriod time.Time) error
	DeleteUser(ctx context.Context, id uint) error
	DeleteTask(ctx context.Context, id uint) error
	GetTask(ctx context.Context, id uint) (*Task, error)
	GetTasks(ctx context.Context, limit, offset uint) ([]Task, error)
	GetUser(ctx context.Context, id uint) (*User, error)
	GetUsers(ctx context.Context, limit, offset uint) ([]User, error)
	GetUserTasks(ctx context.Context, userID, limit, offset uint) ([]Task, error)
}
