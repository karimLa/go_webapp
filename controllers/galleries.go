package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/soramon0/webapp/context"
	"github.com/soramon0/webapp/models"
	"github.com/soramon0/webapp/views"
)

const (
	GalleryShowURL = "gallery_show"
)

// NewGalleries is used to create a new Gallery controller.
// This function will panic if the templates are not
// parsed correctly, and should only be used during
// initial setup.
func NewGalleries(gs models.GalleryService, r *mux.Router, l *log.Logger) *Galleries {
	return &Galleries{
		gs:       gs,
		r:        r,
		l:        l,
		NewView:  views.NewView("bootstrap", "galleries/new"),
		ShowView: views.NewView("bootstrap", "galleries/show"),
	}
}

type Galleries struct {
	gs       models.GalleryService
	r        *mux.Router
	l        *log.Logger
	NewView  *views.View
	ShowView *views.View
}

type CreateGalleryForm struct {
	Title string `schema:"title,required"`
}

// Create is used to process the create gallery form when a user
// submits it. This is used to create a new gallery.
//
// POST /galleries
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form CreateGalleryForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.NewView.Render(w, vd)
		return
	}

	u := context.User(r.Context())

	if u == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	gallery := models.Gallery{
		Title:  form.Title,
		UserID: u.ID,
	}

	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.NewView.Render(w, vd)
		return
	}

	url, err := g.r.Get(GalleryShowURL).URL("id", strconv.Itoa(int(gallery.ID)))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	http.Redirect(w, r, url.Path, http.StatusFound)
}

// Show is used to show gallery.
//
// GET /galleries/:id
func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	var vd views.Data

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		g.l.Println("Error: converting id to int", err)
		vd.SetAlert(models.ErrNotFound)
		g.ShowView.Render(w, vd)
		return
	}

	gallery, err := g.gs.ByID(uint(id))
	if err != nil {
		g.l.Println("Error: fetching gallery", err)
		if err == models.ErrNotFound {
			vd.SetAlert(models.ErrNotFound)
			g.ShowView.Render(w, vd)
			return
		}

		vd.SetAlert(err)
		g.ShowView.Render(w, vd)
		return
	}

	vd.Yield = gallery
	g.ShowView.Render(w, vd)
}
