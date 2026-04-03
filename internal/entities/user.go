package entities

import "time"

type User struct {
	ID        int64
	Login     string
	Password  string
	CreatedAt time.Time
}
