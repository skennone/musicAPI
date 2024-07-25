package data

import (
	"context"
	"database/sql"
	"errors"
	"songs/internal/validator"
	"time"

	"github.com/lib/pq"
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

type SongModel struct {
	DB *sql.DB
}

func (s SongModel) Insert(song *Song) error {
	query := `Insert into songs (title, artist,year,length,genres)
			Values ($1,$2,$3,$4,$5)
			Returning id,created_at,version`
	args := []any{song.Title, song.Artist, song.Year, song.Length, pq.Array(song.Genres)}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return s.DB.QueryRowContext(ctx, query, args...).Scan(&song.ID, &song.CreatedAt, &song.Version)
}

func (s SongModel) Get(id int64) (*Song, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `Select id, created_at, title, artist, year, length,genres,version
			From songs
			where id = $1`
	var song Song

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := s.DB.QueryRowContext(ctx, query, id).Scan(
		&song.ID,
		&song.CreatedAt,
		&song.Title,
		&song.Artist,
		&song.Year,
		&song.Length,
		pq.Array(&song.Genres),
		&song.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &song, nil
}
func (s SongModel) Update(song *Song) error {
	query := `Update songs
			Set title = $1, artist = $2,year = $3, length = $4, genres = $5,version = version + 1
			where id = $6 and version = $7
			Returning version`
	args := []any{
		song.Title,
		song.Artist,
		song.Year,
		song.Length,
		pq.Array(song.Genres),
		song.ID,
		song.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := s.DB.QueryRowContext(ctx, query, args...).Scan(&song.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}
func (s SongModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `Delete from songs
			where id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := s.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
func (s SongModel) GetAll(title, artist string, genres []string, filters Filters) ([]*Song, error) {
	query := `Select id, created_at,title,artist,year,length,genres,version
 			from songs
			where (to_tsvector('simple',title) @@ plainto_tsquery('simple',$1) or $1 = '')
			and (lower (artist) = lower($2) or $2='')
			and (genres @> $3 or $3='{}')
   			order by id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := s.DB.QueryContext(ctx, query, title, artist, pq.Array(genres))
	if err != nil {
		return nil, err
	}
	//Ensure that the result is closed before GetAll() returns
	defer rows.Close()
	songs := []*Song{}

	for rows.Next() {
		var song Song
		err := rows.Scan(
			&song.ID,
			&song.CreatedAt,
			&song.Title,
			&song.Artist,
			&song.Year,
			&song.Length,
			pq.Array(&song.Genres),
			&song.Version,
		)
		if err != nil {
			return nil, err
		}
		songs = append(songs, &song)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return songs, nil
}
