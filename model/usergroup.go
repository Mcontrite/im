package model

import "time"

type UserGroup struct {
	ID        int
	UserID    int
	GroupID   int
	CreatedAt time.Time
}
