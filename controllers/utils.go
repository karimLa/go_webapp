package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

func parseForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	d := schema.NewDecoder()
	if err := d.Decode(dst, r.PostForm); err != nil {
		return parseError(err.Error())
	}

	return nil
}

func Reverse(path, fallback string, r *mux.Router, pathArgs ...string) string {
	url, err := r.Get(path).URL(pathArgs...)
	if err != nil {
		return fallback
	}

	return url.Path
}
