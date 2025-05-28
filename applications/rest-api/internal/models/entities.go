package models

import "time"

type Post struct {
	Id        int
	Title     string
	Content   string
	CreatedAt time.Time
}
