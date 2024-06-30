package main

import (
	"fmt"
	"net/http"
	"songs/internal/data"
	"time"
)

func (app *application) createSongHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new song")
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
