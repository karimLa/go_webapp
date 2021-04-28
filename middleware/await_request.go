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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		amw.wg.Add(1)
		defer amw.wg.Done()
		next.ServeHTTP(w, r)
	})
}

func (mw *awaitRequest) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

func (mw *awaitRequest) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mw.wg.Add(1)
		defer mw.wg.Done()
		next(w, r)
	}
}
