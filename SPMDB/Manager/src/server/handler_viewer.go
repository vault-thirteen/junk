package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"spdb_manager/helper"
	"spdb_manager/models"
	"spdb_manager/settings"

	"github.com/julienschmidt/httprouter"
)

func (srv *Server) hGetActor(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	var id models.ID
	id, err = helper.GetQueryParameterAsId(r, settings.QueryParameterId)
	if err != nil {
		srv.hBadRequest(w, err)
		return
	}

	var actor *models.Actor
	actor, err = srv.dbo.GetActor(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			srv.hNotFound(w, err)
		} else {
			srv.hInternalServerError(w, err)
		}
		return
	}

	w.Header().Add(settings.HeaderContentType, settings.HeaderContentTypeJson)
	err = json.NewEncoder(w).Encode(actor)
	srv.logError(err)
}

func (srv *Server) hGetCodec(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	var id models.ID
	id, err = helper.GetQueryParameterAsId(r, settings.QueryParameterId)
	if err != nil {
		srv.hBadRequest(w, err)
		return
	}

	var codec *models.Codec
	codec, err = srv.dbo.GetCodec(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			srv.hNotFound(w, err)
		} else {
			srv.hInternalServerError(w, err)
		}
		return
	}

	w.Header().Add(settings.HeaderContentType, settings.HeaderContentTypeJson)
	err = json.NewEncoder(w).Encode(codec)
	srv.logError(err)
}

func (srv *Server) hGetContainerType(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	var id models.ID
	id, err = helper.GetQueryParameterAsId(r, settings.QueryParameterId)
	if err != nil {
		srv.hBadRequest(w, err)
		return
	}

	var containerType *models.ContainerType
	containerType, err = srv.dbo.GetContainerType(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			srv.hNotFound(w, err)
		} else {
			srv.hInternalServerError(w, err)
		}
		return
	}

	w.Header().Add(settings.HeaderContentType, settings.HeaderContentTypeJson)
	err = json.NewEncoder(w).Encode(containerType)
	srv.logError(err)
}

func (srv *Server) hGetGenre(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	var id models.ID
	id, err = helper.GetQueryParameterAsId(r, settings.QueryParameterId)
	if err != nil {
		srv.hBadRequest(w, err)
		return
	}

	var genre *models.Genre
	genre, err = srv.dbo.GetGenre(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			srv.hNotFound(w, err)
		} else {
			srv.hInternalServerError(w, err)
		}
		return
	}

	w.Header().Add(settings.HeaderContentType, settings.HeaderContentTypeJson)
	err = json.NewEncoder(w).Encode(genre)
	srv.logError(err)
}

func (srv *Server) hGetMovie(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	var id models.ID
	id, err = helper.GetQueryParameterAsId(r, settings.QueryParameterId)
	if err != nil {
		srv.hBadRequest(w, err)
		return
	}

	var movieRaw *models.MovieRaw
	movieRaw, err = srv.dbo.GetMovie(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			srv.hNotFound(w, err)
		} else {
			srv.hInternalServerError(w, err)
		}
		return
	}

	var movie *models.Movie
	movie, err = movieRaw.Parse()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			srv.hNotFound(w, err)
		} else {
			srv.hInternalServerError(w, err)
		}
		return
	}

	w.Header().Add(settings.HeaderContentType, settings.HeaderContentTypeJson)
	err = json.NewEncoder(w).Encode(movie)
	srv.logError(err)
}

func (srv *Server) hGetServer(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	var id models.ID
	id, err = helper.GetQueryParameterAsId(r, settings.QueryParameterId)
	if err != nil {
		srv.hBadRequest(w, err)
		return
	}

	var server *models.Server
	server, err = srv.dbo.GetServer(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			srv.hNotFound(w, err)
		} else {
			srv.hInternalServerError(w, err)
		}
		return
	}

	w.Header().Add(settings.HeaderContentType, settings.HeaderContentTypeJson)
	err = json.NewEncoder(w).Encode(server)
	srv.logError(err)
}

func (srv *Server) hGetTag(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	var id models.ID
	id, err = helper.GetQueryParameterAsId(r, settings.QueryParameterId)
	if err != nil {
		srv.hBadRequest(w, err)
		return
	}

	var tag *models.Tag
	tag, err = srv.dbo.GetTag(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			srv.hNotFound(w, err)
		} else {
			srv.hInternalServerError(w, err)
		}
		return
	}

	w.Header().Add(settings.HeaderContentType, settings.HeaderContentTypeJson)
	err = json.NewEncoder(w).Encode(tag)
	srv.logError(err)
}

func (srv *Server) hGetMoviesCount(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	var moviesCount uint
	moviesCount, err = srv.dbo.GetMoviesCount()
	if err != nil {
		srv.hInternalServerError(w, err)
		return
	}

	_, err = fmt.Fprint(w, moviesCount)
	srv.logError(err)
}
