package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Malachy-Olua/social-platform/cmd/api/handlers"
	"github.com/Malachy-Olua/social-platform/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// serve() or run()

type dbConfig struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  time.Duration
}

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr string
	db   dbConfig
	env  string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	// mux:= http.NewServeMux()
	// mux.HandleFunc("GET /v1/health", app.healthcheckHandler)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	// r.Use(middleware.RealIP) // this is deprectaed
	r.Use(middleware.Logger)
	// r.Get("/v1/health", app.healthcheckHandler)

	r.Use(middleware.Timeout(60 * time.Second))

	// initialize handlers
	postHandler := handlers.NewPostHandler(app.store)
	userHandler := handlers.NewUserHandler(app.store)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthcheckHandler)
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", postHandler.CreatePostHandler)
			r.Get("/", postHandler.ListPostsHandler)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", postHandler.GetPostHandler)
				r.Put("/", postHandler.UpdatePostHandler)
				r.Delete("/", postHandler.DeletePostHandler)
			})
		})
		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.CreateUserHandler)
		})
	})

	return r

}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Starting server on %s", app.config.addr)

	return srv.ListenAndServe()
}
