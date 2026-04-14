package main

import (
	"net/http"

	"github.com/Kenwanjohi/yalsh-go/internal/data"
	"github.com/golang-jwt/jwt/v5"
)

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("accessToken")
		if err != nil {
			r = app.contextSetUser(r, &data.User{})
			next.ServeHTTP(w, r)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (any, error) {
			return []byte(app.config.jwt.secret), nil
		})

		if err != nil || !token.Valid {
			app.errorResponse(w, r, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		user, err := app.models.Users.GetByID(claims.UserID)
		if err != nil {
			app.errorResponse(w, r, http.StatusUnauthorized, "user not found")
			return
		}

		r = app.contextSetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		if user == nil || user.ID == 0 {
			app.errorResponse(w, r, http.StatusUnauthorized, "you must be authenticated to access this resource")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// List of allowed origins
		allowedOrigins := []string{
			"http://localhost:5173", // Local dev (Vite default)
			"http://localhost:3000", // Local dev (Alternative)
			"https://app.yalsh.com", // Production frontend
		}

		for _, o := range allowedOrigins {
			if origin == o {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
