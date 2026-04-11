package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(app.authenticate)

	r.Get("/v1/healthcheck", app.healthcheckHandler)
	r.Post("/v1/users", app.registerUserHandler)
	r.Post("/v1/tokens", app.loginHandler)
	r.Delete("/v1/tokens", app.logoutHandler)

	r.Group(func(r chi.Router) {
		r.Use(app.requireAuthenticated)

		r.Get("/v1/users/me", app.showUserHandler)
		r.Put("/v1/users/me", app.updateProfileHandler)

		r.Post("/v1/links", app.createLinkHandler)
		r.Put("/v1/links/{id}", app.updateLinkHandler)
		r.Get("/v1/links", app.listLinksHandler)
		r.Delete("/v1/links/{id}", app.deleteLinkHandler)
	})

	r.Get("/{key}", app.redirectLinkHandler)

	return r
}
