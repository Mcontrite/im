package model

import "time"

type Friend struct {
	ID        int
	UserID    int
	User2ID   int
	CreatedAt time.Time
}
