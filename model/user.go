package model

import "time"

type User struct {
	ID        int
	Username  string
	Password  string
	Phone     string
	Avatar    string
	Sex       string
	CreatedAt time.Time
	Token     string
	Salt      string
	IsOnline  int
}
