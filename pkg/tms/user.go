package tms

import "time"

//User contains information about a user
type User struct {
	ID        uint
	FirstName string
	LastName  string
	Email     string
	CreatedAt time.Time
}
