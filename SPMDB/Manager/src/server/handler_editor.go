package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"spdb_manager/helper"
	"spdb_manager/models"
	"spdb_manager/settings"

	"github.com/julienschmidt/httprouter"
)

func (srv *Server) hAddActor(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	var actor = new(models.Actor)
	err = json.NewDecoder(r.Body).Decode(actor)
	if err != nil {
		srv.hBadRequest(w, err)
		return
	}

	var id models.ID
	id, err = srv.dbo.AddActor(actor)
	if err != nil {
		srv.hInternalServerError(w, err)
		return
	}

	_, err = fmt.Fprint(w, id)
	srv.logError(err)
}

func (srv *Server) hAddCodec(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	var codec = new(models.Codec)
	err = json.NewDecoder(r.Body).Decode(codec)
	if err != nil {
		srv.hBadRequest(w, err)
		return
	}

	var id models.ID
	id, err = srv.dbo.AddCodec(codec)
	if err != nil {
		srv.hInternalServerError(w, err)
		return
	}

	_, err = fmt.Fprint(w, id)
	srv.logError(err)
}

func (srv *Server) hAddContainerType(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	var containerType = new(models.ContainerType)
	err = json.NewDecoder(r.Body).Decode(containerType)
	if err != nil {
		srv.hBadRequest(w, err)
		return
	}

	var id models.ID
	id, err = srv.dbo.AddContainerType(containerType)
	if err != nil {
		srv.hInternalServerError(w, err)
		return
	}

	_, err = fmt.Fprint(w, id)
	srv.logError(err)
}

func (srv *Server) hAddGenre(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	var genre = new(models.Genre)
	err = json.NewDecoder(r.Body).Decode(genre)
	if err != nil {
		srv.hBadRequest(w, err)
		return
	}

	var id models.ID
	id, err = srv.dbo.AddGenre(genre)
	if err != nil {
		srv.hInternalServerError(w, err)
		return
	}

	_, err = fmt.Fprint(w, id)
	srv.logError(err)
}

func (srv *Server) hAddServer(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	var server = new(models.Server)
	err = json.NewDecoder(r.Body).Decode(server)
	if err != nil {
		srv.hBadRequest(w, err)
		return
	}

	var id models.ID
	id, err = srv.dbo.AddServer(server)
	if err != nil {
		srv.hInternalServerError(w, err)
		return
	}

	_, err = fmt.Fprint(w, id)
	srv.logError(err)
}

func (srv *Server) hAddTag(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	var tag = new(models.Tag)
	err = json.NewDecoder(r.Body).Decode(tag)
	if err != nil {
		srv.hBadRequest(w, err)
		return
	}

	var id models.ID
	id, err = srv.dbo.AddTag(tag)
	if err != nil {
		srv.hInternalServerError(w, err)
		return
	}

	_, err = fmt.Fprint(w, id)
	srv.logError(err)
}

func (srv *Server) hSetMovieTitle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	var movie *models.Movie
	movie, err = srv.getMovieFromBody(r)
	if err != nil {
		srv.hBadRequest(w, err)
		return
	}

	err = srv.dbo.SetMovieTitle(movie)
	if err != nil {
		srv.hInternalServerError(w, err)
		return
	}
}

func (srv *Server) hSetMovieDescription(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	var movie *models.Movie
	movie, err = srv.getMovieFromBody(r)
	if err != nil {
		srv.hBadRequest(w, err)
		return
	}

	err = srv.dbo.SetMovieDescription(movie)
	if err != nil {
		srv.hInternalServerError(w, err)
		return
	}
}

func (srv *Server) hIncMovieViewsCount(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	var movie *models.Movie
	movie, err = srv.getMovieFromBody(r)
	if err != nil {
		srv.hBadRequest(w, err)
		return
	}

	err = srv.dbo.IncMovieViewsCount(movie)
	if err != nil {
		srv.hInternalServerError(w, err)
		return
	}
}

func (srv *Server) hSetMovieHardcoreLevel(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	var movie *models.Movie
	movie, err = srv.getMovieFromBody(r)
	if err != nil {
		srv.hBadRequest(w, err)
		return
	}

	err = srv.dbo.SetMovieHardcoreLevel(movie)
	if err != nil {
		srv.hInternalServerError(w, err)
		return
	}
}

func (srv *Server) hAssignTagToMovie(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	var movieId models.ID
	movieId, err = helper.GetQueryParameterAsId(r, settings.QueryParameterMovieId)
	if err != nil {
		srv.hBadRequest(w, err)
		return
	}

	var tagId models.ID
	tagId, err = helper.GetQueryParameterAsId(r, settings.QueryParameterTagId)
	if err != nil {
		srv.hBadRequest(w, err)
		return
	}

	err = srv.dbo.AssignTagToMovie(movieId, tagId)
	if err != nil {
		srv.hInternalServerError(w, err)
		return
	}
}

func (srv *Server) hRetractTagFromMovie(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	var movieId models.ID
	movieId, err = helper.GetQueryParameterAsId(r, settings.QueryParameterMovieId)
	if err != nil {
		srv.hBadRequest(w, err)
		return
	}

	var tagId models.ID
	tagId, err = helper.GetQueryParameterAsId(r, settings.QueryParameterTagId)
	if err != nil {
		srv.hBadRequest(w, err)
		return
	}

	err = srv.dbo.RetractTagFromMovie(movieId, tagId)
	if err != nil {
		srv.hInternalServerError(w, err)
		return
	}
}
