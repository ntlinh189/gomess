package models

import "time"

type Friend struct {
	User1ID   int64
	User2ID   int64
	CreatedAt time.Time
}
