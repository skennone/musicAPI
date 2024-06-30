package data

import "time"

type Song struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Artist    string    `json:"artist"`
	Year      int32     `json:"year,omitempty"`
	Length    Length    `json:"length,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}
