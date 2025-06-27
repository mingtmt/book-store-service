package domain

import "time"

type Book struct {
	ID        string
	Title     string
	Author    string
	Price     string
	CreatedAt time.Time
}
