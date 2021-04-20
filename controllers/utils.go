package controllers

import (
	"net/http"

	"github.com/gorilla/schema"
)

func parseForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	d := schema.NewDecoder()
	if err := d.Decode(dst, r.PostForm); err != nil {
		return err
	}

	return nil
}
