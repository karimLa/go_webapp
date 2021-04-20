package controllers

import (
	"sync"

	"github.com/karimla/webapp/views"
)

func NewStatic(wg *sync.WaitGroup) *Static {
	return &Static{
		Home:    views.NewView(wg, "bootstrap", "static/index"),
		Contact: views.NewView(wg, "bootstrap", "static/contact"),
	}
}

type Static struct {
	Home    *views.View
	Contact *views.View
}
