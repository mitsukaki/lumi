package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mitsukaki/lumi/internal/handler"
)

func main() {
	exePath, _ := os.Executable()
	binaryPath := filepath.Dir(exePath)

	router := chi.NewRouter()

	router.Use(middleware.Logger)

	// API subroute
	router.Route("/api", func(api chi.Router) {
		api.Get("/ping", handler.Unimplemented)

		// user routes
		api.Route("/user", func(user chi.Router) {
			// TODO: user.Use(r.UserCtx)

			user.Get("/{userId}", handler.Unimplemented)
			user.Put("/{userId}", handler.Unimplemented)
		})

		// album routes
		api.Route("/album", func(album chi.Router) {
			album.Post("/", handler.Unimplemented)

			// photo routes
			album.Route("/photos", func(photos chi.Router) {
				photos.Get("/", handler.Unimplemented)
				photos.Post("/", handler.Unimplemented)
				photos.Delete("/", handler.Unimplemented)
			})

			// album manipulation routes
			album.Route("/{albumId}", func(album chi.Router) {
				// TODO: album.Use(AlbumCtx)

				album.Get("/", handler.Unimplemented)
				album.Put("/", handler.Unimplemented)
				album.Delete("/", handler.Unimplemented)
			})
		})
	})

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(binaryPath, "public", "index.html"))
	})

	// serve static content
	fs := http.FileServer(http.Dir(filepath.Join(binaryPath, "public")))
	router.Handle("/*", fs)

	// listen and serve
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port
	}

	http.ListenAndServe(":"+port, router)
}
