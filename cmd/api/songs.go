package main

import (
	"fmt"
	"net/http"
	"songs/internal/data"
	"songs/internal/validator"
	"time"
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

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showSongHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	song := data.Song{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Von dutch",
		Artist:    "Charlie XCX ",
		Length:    3,
		Genres:    []string{"Goth Pop", "Dark Pop", "Elektropop"},
		Version:   1,
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"song": song}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
