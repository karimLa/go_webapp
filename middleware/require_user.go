package middleware

import (
	"net/http"

	"github.com/soramon0/webapp/context"
	"github.com/soramon0/webapp/models"
)

type requireUser struct {
	models.UserService
}

func NewRequireUser(us models.UserService) *requireUser {
	return &requireUser{us}
}

// Middleware function, which will be called for each request
func (mw *requireUser) Middleware(next http.Handler) http.Handler {
	return mw.ApplyFn(next.ServeHTTP)
}

func (mw *requireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

func (mw *requireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		user, err := mw.ByRemember(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)

		next(w, r)
	}
}
