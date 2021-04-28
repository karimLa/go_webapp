package controllers

import (
	"fmt"
	"net/http"
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

type CreateGalleryForm struct {
	Title string `schema:"title,required"`
}

// Create is used to process the create gallery form when a user
// submits it. This is used to create a new gallery.
//
// POST /galleries
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	g.wg.Add(1)
	defer g.wg.Done()

	var vd views.Data
	var form CreateGalleryForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.NewView.Render(w, vd)
		return
	}

	gallery := models.Gallery{
		Title: form.Title,
	}

	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.NewView.Render(w, vd)
		return
	}

	fmt.Fprintln(w, gallery)
}
