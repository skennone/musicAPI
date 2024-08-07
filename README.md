# REST JSON API Project
Simple JSON API for retrieving and managing information about music.

## Setup
To run the project locally, clone the repo
`git@github.com:skennone/songsAPI.git`

Run the following commands from the provided `Makefile` to run the project

`make redis`

`make run/server`

## Routes:
### Check it out at [routes](./cmd/api/routes.go)
```
GET /v1/songs
// Show the details of all songs

POST /v1/songs
// Create a new song

GET /v1/songs/:id
// Show the details of a specific song

DELETE /v1/songs/:id
// Delete a specific song

PATCH /v1/songs/:id
// Update the details of a specific song

POST /v1/users/
// Register a new user

PUT /v1/users/activated
// Activate a specific user
```

## TODO:
- [ ] Authentication
- [ ] User Activation
- [ ] ...
