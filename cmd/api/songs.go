package main

import (
	"errors"
	"fmt"
	"net/http"
	"songs/internal/data"
	"songs/internal/validator"
)

func (app *application) createSongHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string      `json:"title"`
		Artist string      `json:"artist"`
		Year   int32       `json:"year"`
		Length data.Length `json:"length"`
		Genres []string    `json:"genres"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	song := &data.Song{
		Title:  input.Title,
		Artist: input.Artist,
		Year:   input.Year,
		Length: input.Length,
		Genres: input.Genres,
	}
	v := validator.New()

	if data.ValidateSong(v, song); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Songs.Insert(song)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/songs/%d", song.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"song": song}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showSongHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	song, err := app.models.Songs.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"song": song}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateSongHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	song, err := app.models.Songs.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
	}
	var input struct {
		Title  *string      `json:"title"`
		Artist *string      `json:"artist"`
		Year   *int32       `json:"year"`
		Length *data.Length `json:"length"`
		Genres []string     `json:"genres"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if input.Title != nil {
		song.Title = *input.Title
	}
	if input.Artist != nil {
		song.Artist = *input.Artist
	}
	if input.Year != nil {
		song.Year = *input.Year
	}
	if input.Length != nil {
		song.Length = *input.Length
	}
	if input.Genres != nil {
		song.Genres = input.Genres
	}
	v := validator.New()
	if data.ValidateSong(v, song); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Songs.Update(song)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"song": song}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteSongHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.Songs.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"messages": "song successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) listSongsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Artist string
		Genres []string
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()
	input.Title = app.readString(qs, "title", "")
	input.Artist = app.readString(qs, "artist", "")
	input.Genres = app.readCSV(qs, "genres", []string{})
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "artist", "year", "length", "-id", "-title", "-artist", "-year", "-length"}
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	songs, err := app.models.Songs.GetAll(input.Title, input.Artist, input.Genres, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"songs": songs}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
