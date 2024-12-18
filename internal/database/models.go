// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"database/sql"
	"time"
)

type Feed struct {
	ID            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	LastFetchedAt sql.NullTime
	Name          string
	Url           string
	UserID        string
}

type FeedFollow struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    string
	FeedID    string
}

type Post struct {
	ID          string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       sql.NullString
	Url         string
	Description sql.NullString
	PublishedAt time.Time
	FeedID      string
}

type User struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
}
