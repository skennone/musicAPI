package data

import (
	"songs/internal/validator"
	"time"
)

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

func ValidateSong(v *validator.Validator, song *Song) {
	v.Check(song.Title != "", "title", "must be provided")
	v.Check(len(song.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(song.Artist != "", "artist", "must be provided")
	v.Check(len(song.Artist) <= 500, "artist", "must not be more than 500 bytes long")

	v.Check(song.Year != 0, "year", "must be provided")
	v.Check(song.Year >= 1888, "year", "must be greater than 1888")
	v.Check(song.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(song.Length != 0, "length", "must be provided")
	v.Check(song.Length > 0, "length", "must be a positive integer")

	v.Check(song.Genres != nil, "genres", "must be provided")
	v.Check(len(song.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(song.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(song.Genres), "genres", "must not contain duplicate values")
}
