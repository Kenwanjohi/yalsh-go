package main

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Kenwanjohi/yalsh-go/internal/data"
	"github.com/Kenwanjohi/yalsh-go/internal/validator"
	"github.com/go-chi/chi/v5"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     "1.0.0",
		},
	}
	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Username: input.Username,
		Email:    input.Email,
	}
	err = user.Password.Set(input.Password)

	err = app.models.Users.Insert(user)
	if err != nil {
		if errors.Is(err, data.ErrDuplicateEmail) {
			app.errorResponse(w, r, http.StatusUnprocessableEntity, "a user with this email address already exists")
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	token, err := app.generateJWT(user.ID, user.Username)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.setAuthCookie(w, token, time.Now().Add(14*24*time.Hour))

	err = app.writeJSON(w, http.StatusCreated, envelope{"registered": true}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createLinkHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		URL string `json:"url"`
		Key string `json:"key"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)
	link := &data.Link{
		URL:    input.URL,
		Key:    input.Key,
		UserID: user.ID,
	}

	v := validator.New()
	data.ValidateLink(v, link)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Links.Insert(r.Context(), link)
	if err != nil {
		// Handle potential duplicate key errors here later
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"link": link}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) redirectLinkHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	url, err := app.models.Links.GetAndIncrementClicks(r.Context(), key)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)

}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		app.errorResponse(w, r, http.StatusUnauthorized, "invalid email or password")
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil || !match {
		app.errorResponse(w, r, http.StatusUnauthorized, "invalid email or password")
		return
	}

	token, err := app.generateJWT(user.ID, user.Username)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.setAuthCookie(w, token, time.Now().Add(14*24*time.Hour))

	app.writeJSON(w, http.StatusOK, envelope{"authenticated": true}, nil)
}

func (app *application) listLinksHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	links, err := app.models.Links.GetAllForUser(r.Context(), user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"links": links}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showUserHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	err := app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteLinkHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	user := app.contextGetUser(r)
	err = app.models.Links.Delete(r.Context(), id, user.ID)
	if err != nil {
		if errors.Is(err, errors.New("record not found or unauthorized")) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "link deleted successfully"}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	app.setAuthCookie(w, "", time.Time{})

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) updateProfileHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username *string `json:"username"`
		Password *string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	if input.Username != nil {
		user.Username = *input.Username
	}

	if input.Password != nil {
		err = user.Password.Set(*input.Password)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	err = app.models.Users.Update(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateLinkHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		URL *string `json:"url"`
		Key *string `json:"key"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	link, err := app.models.Links.GetByID(r.Context(), id, user.ID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	if input.URL != nil {
		link.URL = *input.URL
	}

	if input.Key != nil {
		link.Key = *input.Key
	}

	err = app.models.Links.Update(r.Context(), link)
	if err != nil {
		if errors.Is(err, data.ErrDuplicateKey) {
			app.errorResponse(w, r, http.StatusUnprocessableEntity, "this key is already in use")
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"link": link}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	err := app.models.Users.Delete(r.Context(), user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.setAuthCookie(w, "", time.Time{})

	err = app.writeJSON(w, http.StatusOK, envelope{"deleted": true}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
