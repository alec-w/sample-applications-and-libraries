package models

import "time"

type Post struct {
	Id        string
	Title     string
	Content   string
	CreatedAt time.Time
}
