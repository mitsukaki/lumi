package api

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
		r.Use(cors.Handler(cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins: []string{"https://*", "http://*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))

		r.Get("/ping", apiServer.Ping)

		// user routes
		r.Route("/user/{userId}", func(user chi.Router) {
			user.Use(apiServer.UserContext)

			user.Get("/", apiServer.GetUser)
			user.Post("/", apiServer.PostUser)
		})

		// album routes
		r.Route("/album/{albumId}", func(album chi.Router) {
			album.Use(apiServer.AlbumContext)

			album.Get("/", apiServer.GetAlbum)
			album.Put("/", apiServer.Unimplemented)
			album.Post("/", apiServer.Unimplemented)
			album.Delete("/", apiServer.Unimplemented)
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
