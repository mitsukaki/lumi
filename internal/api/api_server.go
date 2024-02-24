package api

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/mitsukaki/lumi/internal/db"
)

type APIServer struct {
	// router for the APIServer
	r *chi.Mux

	// database connection
	db *db.CouchDatabase

	// tables
	UserTable  *db.CouchTable
	AlbumTable *db.CouchTable
	PhotoTable *db.CouchTable

	// logger
	logger *zap.Logger
}

type APIConfig struct {
	StaticDir string
	DBConfig  *db.CouchDBConfig
}

// CreateAPI creates a new APIServer instance
func CreateAPIServer(config APIConfig) (*APIServer, error) {
	// create the APIServer instance
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	apiServer := &APIServer{
		r:      chi.NewRouter(),
		db:     db.CreateDatabase(config.DBConfig),
		logger: logger,
	}

	// create the tables
	apiServer.UserTable, err = apiServer.db.CreateTable("lumi_user")
	if err != nil {
		return nil, err
	}

	apiServer.AlbumTable, err = apiServer.db.CreateTable("lumi_album")
	if err != nil {
		return nil, err
	}

	apiServer.PhotoTable, err = apiServer.db.CreateTable("lumi_photo")
	if err != nil {
		return nil, err
	}

	// add middleware
	r := apiServer.r
	r.Use(middleware.Logger)

	// API subroute
	r.Route("/api", func(r chi.Router) {
		r.Get("/ping", apiServer.Ping)

		// user routes
		r.Route("/user", func(user chi.Router) {
			user.Use(apiServer.UserContext)

			user.Get("/{userId}", apiServer.GetUser)
			user.Put("/{userId}", apiServer.PutUser)
		})

		// album routes
		r.Route("/album", func(album chi.Router) {
			album.Post("/", apiServer.Unimplemented)

			// photo routes
			album.Route("/photos", func(photos chi.Router) {
				photos.Get("/", apiServer.Unimplemented)
				photos.Post("/", apiServer.Unimplemented)
				photos.Delete("/", apiServer.Unimplemented)
			})

			// album manipulation routes
			album.Route("/{albumId}", func(album chi.Router) {
				album.Use(apiServer.AlbumContext)

				album.Get("/", apiServer.Unimplemented)
				album.Put("/", apiServer.Unimplemented)
				album.Delete("/", apiServer.Unimplemented)
			})
		})
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(config.StaticDir, "index.html"))
	})

	// serve static content
	fs := http.FileServer(http.Dir(config.StaticDir))
	r.Handle("/*", fs)

	// return the APIServer instance
	return apiServer, nil
}

func (apiServer *APIServer) Router() *chi.Mux {
	return apiServer.r
}

func (apiServer *APIServer) Start() {
	// listen and serve
	port := strconv.Itoa(viper.GetInt("http.port"))

	apiServer.logger.Info("Starting server on port " + port)
	http.ListenAndServe(":"+port, apiServer.r)
}
