package server

import (
	"net/http"
	"path"
	"spdb_manager/settings"

	"github.com/julienschmidt/httprouter"
)

func (srv *Server) initRouter() (err error) {
	router := httprouter.New()

	// System Handlers.
	router.POST("/system/shutdown", srv.hSystemShutdown)

	// Simple Handlers.
	router.GET("/", srv.hIndex)

	// Viewer's Handlers.
	router.GET("/actor", srv.hGetActor)
	router.GET("/codec", srv.hGetCodec)
	//router.GET("/comment", srv.hGetComment)
	router.GET("/container_type", srv.hGetContainerType)
	router.GET("/genre", srv.hGetGenre)
	router.GET("/server", srv.hGetServer)
	router.GET("/tag", srv.hGetTag)
	//router.GET("/user", srv.hGetUser)

	// Movie.
	router.GET("/movie", srv.hGetMovie)
	router.GET("/movies/count", srv.hGetMoviesCount)

	// Editor's Handlers.
	if srv.isEditorModeUsed {
		router.PUT("/actor", srv.hAddActor)
		router.PUT("/codec", srv.hAddCodec)
		//router.PUT("/comment", srv.hAddComment)
		router.PUT("/container_type", srv.hAddContainerType)
		router.PUT("/genre", srv.hAddGenre)
		router.PUT("/server", srv.hAddServer)
		router.PUT("/tag", srv.hAddTag)
		//router.PUT("/user", srv.hAddUser)

		// Movie.
		router.POST("/movie/title", srv.hSetMovieTitle)
		router.POST("/movie/description", srv.hSetMovieDescription)
		router.POST("/movie/views_count", srv.hIncMovieViewsCount)
		router.POST("/movie/hardcore_level", srv.hSetMovieHardcoreLevel)
		router.POST("/movie/tag/assign", srv.hAssignTagToMovie)
		router.POST("/movie/tag/retract", srv.hRetractTagFromMovie)
	}

	guiPath := path.Join("/", settings.GuiUrlFolder, "*filepath")
	router.ServeFiles(guiPath, http.Dir(srv.guiFolder))

	router.GlobalOPTIONS =
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Access-Control-Request-Method") != "" {
				// Set CORS headers.
				header := w.Header()
				header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
				header.Set("Access-Control-Allow-Origin", "*")
			}

			// Adjust status code to 204.
			w.WriteHeader(http.StatusNoContent)
		})

	srv.httpRouter = router

	return nil
}
