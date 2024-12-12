package server

import (
	"encoding/json"
	"net/http"
	"spdb_manager/models"
)

func (srv *Server) getMovieFromBody(r *http.Request) (movie *models.Movie, err error) {
	var movieRaw = new(models.MovieRaw)
	err = json.NewDecoder(r.Body).Decode(movieRaw)
	if err != nil {
		return nil, err
	}

	movie, err = movieRaw.Parse()
	if err != nil {
		return nil, err
	}

	return movie, nil
}
