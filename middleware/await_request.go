package middleware

import (
	"net/http"
	"sync"
)

type awaitRequest struct {
	wg *sync.WaitGroup
}

func NewAwaitRequest(wg *sync.WaitGroup) *awaitRequest {
	return &awaitRequest{wg: wg}
}

// Middleware function, which will be called for each request
func (amw *awaitRequest) Middleware(next http.Handler) http.Handler {
	return amw.ApplyFn(next.ServeHTTP)
}

func (amw *awaitRequest) Apply(next http.Handler) http.HandlerFunc {
	return amw.ApplyFn(next.ServeHTTP)
}

func (amw *awaitRequest) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		amw.wg.Add(1)
		defer amw.wg.Done()
		next(w, r)
	}
}
