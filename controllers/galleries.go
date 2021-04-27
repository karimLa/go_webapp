package controllers

import (
	"sync"

	"github.com/karimla/webapp/models"
	"github.com/karimla/webapp/views"
)

// NewGalleries is used to create a new Gallery controller.
// This function will panic if the templates are not
// parsed correctly, and should only be used during
// initial setup.
func NewGalleries(wg *sync.WaitGroup, gs models.GalleryService) *Galleries {
	return &Galleries{
		gs:      gs,
		wg:      wg,
		NewView: views.NewView(wg, "bootstrap", "galleries/new"),
	}
}

type Galleries struct {
	gs      models.GalleryService
	wg      *sync.WaitGroup
	NewView *views.View
}
