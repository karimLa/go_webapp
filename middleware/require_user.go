package middleware

import (
	"net/http"

	"soramon0/webapp/context"
	"soramon0/webapp/models"
)

type user struct {
	models.UserService
}

func NewUser(us models.UserService) *user {
	return &user{us}
}

// Middleware function, which will be called for each request
func (mw *user) Middleware(next http.Handler) http.Handler {
	return mw.ApplyFn(next.ServeHTTP)
}

func (mw *user) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

func (mw *user) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			next(w, r)
			return
		}

		user, err := mw.ByRemember(cookie.Value)
		if err != nil {
			next(w, r)
			return
		}

		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)

		next(w, r)
	}
}

type requireUser struct {
	user
}

// RequireUser needs the user middleware
// otherwise it will not work correctly.
func NewRequireUser(u user) *requireUser {
	return &requireUser{user: u}
}

// Middleware needs the user middleware
// otherwise it will not work correctly.
func (mw *requireUser) Middleware(next http.Handler) http.Handler {
	return mw.ApplyFn(next.ServeHTTP)
}

// Apply needs the user middleware
// otherwise it will not work correctly.
func (mw *requireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

// ApplyFn needs the user middleware
// otherwise it will not work correctly.
func (mw *requireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return mw.user.ApplyFn(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next(w, r)
	})
}
