package model

import "time"

type Group struct {
	ID        int
	LeaderID  int
	Groupname string
	Avatar    string
	CreatedAt time.Time
}
