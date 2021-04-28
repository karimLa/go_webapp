package controllers

import (
	"github.com/karimla/webapp/views"
)

func NewStatic() *Static {
	return &Static{
		HomeView:    views.NewView("bootstrap", "static/index"),
		ContactView: views.NewView("bootstrap", "static/contact"),
	}
}

type Static struct {
	HomeView    *views.View
	ContactView *views.View
}
