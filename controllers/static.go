package controllers

import (
	"sync"

	"github.com/karimla/webapp/views"
)

func NewStatic(wg *sync.WaitGroup) *Static {
	return &Static{
		HomeView:    views.NewView(wg, "bootstrap", "static/index"),
		ContactView: views.NewView(wg, "bootstrap", "static/contact"),
	}
}

type Static struct {
	HomeView    *views.View
	ContactView *views.View
}
