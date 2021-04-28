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
