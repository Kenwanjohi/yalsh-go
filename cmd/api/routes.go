package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(app.enableCORS)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(app.authenticate)

	r.Get("/v1/healthcheck", app.healthcheckHandler)
	r.Post("/v1/accounts", app.registerUserHandler)
	r.Post("/v1/session", app.loginHandler)
	r.Delete("/v1/session", app.logoutHandler)

	r.Group(func(r chi.Router) {
		r.Use(app.requireAuthenticated)

		r.Get("/v1/account", app.showUserHandler)
		r.Patch("/v1/account", app.updateProfileHandler)
		r.Delete("/v1/account", app.deleteAccountHandler)

		r.Post("/v1/links", app.createLinkHandler)
		r.Get("/v1/links", app.listLinksHandler)
		r.Put("/v1/links/{id}", app.updateLinkHandler)
		r.Delete("/v1/links/{id}", app.deleteLinkHandler)
	})

	r.Get("/{key}", app.redirectLinkHandler)

	return r
}
