package api

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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

	// database connections
	db  *db.CouchDatabase
	svc *s3.S3

	// tables
	UserTable     *db.CouchTable
	AlbumTable    *db.CouchTable
	UsernameTable *db.CouchTable

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

	apiServer.UsernameTable, err = apiServer.db.CreateTable("lumi_uuid")
	if err != nil {
		return nil, err
	}

	// create amazon s3 client
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(viper.GetString("s3.region")),
		Endpoint: aws.String(viper.GetString("s3.endpoint")),
	})
	if err != nil {
		return nil, err
	}

	creds := credentials.NewStaticCredentials(
		viper.GetString("s3.access_key.id"),
		viper.GetString("s3.access_key.secret"),
		"",
	)

	apiServer.svc = s3.New(sess, &aws.Config{Credentials: creds})

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
			user.Put("/album", apiServer.PutAlbum)
		})

		// album routes
		r.Route("/album/{albumId}", func(album chi.Router) {
			album.Use(apiServer.AlbumContext)

			album.Put("/", apiServer.UploadAlbum)
			album.Get("/", apiServer.GetAlbum)
			album.Post("/", apiServer.Unimplemented)
			album.Delete("/", apiServer.Unimplemented)
		})
	})

	// r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, filepath.Join(config.StaticDir, "index.html"))
	// })

	// serve static content
	fs := http.FileServer(http.Dir(filepath.Join(config.StaticDir, "assets")))
	r.Handle("/assets/*", http.StripPrefix("/assets/", fs))

	// redirect any other requests to index.html
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(config.StaticDir, "index.html"))
	})

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
