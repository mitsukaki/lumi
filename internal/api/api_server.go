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
	"github.com/mitsukaki/lumi/internal/endpoints"
	mid "github.com/mitsukaki/lumi/internal/middleware"
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

	authMiddleware := mid.CreateAuthHandler(
		logger,
		apiServer.UserTable,
	)

	r.Use(authMiddleware.Middleware)

	albumContext := mid.CreateAlbumContextHandler(
		logger,
		apiServer.AlbumTable,
	)

	// Create the endpoint handlers
	handler := endpoints.CreateEndpointHandler(
		logger,
		apiServer.UserTable,
		apiServer.AlbumTable,
		apiServer.svc,
	)

	// API subroute
	r.Route("/api/v1", func(r chi.Router) {
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

		r.Get("/ping", handler.Ping) // done

		// user endpoints
		r.Post("/users", handler.CreateUser)
		r.Post("/users/login", handler.LoginUser)
		r.Post("/users/logout", handler.LogoutUser)

		r.Route("/users/{userId}", func(user chi.Router) {
			user.Get("/", handler.GetUser)
			user.Put("/", handler.UpdateUser)
			user.Delete("/", handler.DeleteUser)

			user.Get("/albums", handler.GetUserAlbums)
			user.Delete("/albums/{albumId}", handler.DeleteUserAlbum)
		})

		// album routes
		r.Post("/albums", handler.CreateAlbum)

		r.Route("/albums/{albumId}", func(album chi.Router) {
			album.Use(albumContext.Middleware)

			album.Get("/", handler.GetAlbum)
			album.Put("/", handler.UpdateAlbum)
			album.Delete("/", handler.Unimplemented)

			album.Post("/photo", handler.UploadAlbumPhoto)
			album.Delete("/photo", handler.DeleteAlbumPhoto)
		})
	})

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
