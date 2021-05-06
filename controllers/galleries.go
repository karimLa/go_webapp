package controllers

import (
	"log"
	"net/http"
	"strconv"

	"webapp/context"
	"webapp/models"
	"webapp/views"

	"github.com/gorilla/mux"
)

const (
	GalleryShowURL    = "gallery_show"
	GalleryEditURL    = "gallery_edit"
	GalleriesIndexURL = "gallery_index"
)

// NewGalleries is used to create a new Gallery controller.
// This function will panic if the templates are not
// parsed correctly, and should only be used during
// initial setup.
func NewGalleries(gs models.GalleryService, r *mux.Router, l *log.Logger) *Galleries {
	return &Galleries{
		gs:        gs,
		r:         r,
		l:         l,
		IndexView: views.NewView("bootstrap", "galleries/index"),
		NewView:   views.NewView("bootstrap", "galleries/new"),
		ShowView:  views.NewView("bootstrap", "galleries/show"),
		EditView:  views.NewView("bootstrap", "galleries/edit"),
	}
}

type Galleries struct {
	gs        models.GalleryService
	r         *mux.Router
	l         *log.Logger
	IndexView *views.View
	NewView   *views.View
	ShowView  *views.View
	EditView  *views.View
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
		g.NewView.Render(w, r, vd)
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
		g.NewView.Render(w, r, vd)
		return
	}

	url, err := g.r.Get(GalleryEditURL).URL("id", strconv.Itoa(int(gallery.ID)))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	http.Redirect(w, r, url.Path, http.StatusFound)
}

// Index is used to show the user galleries.
//
// GET /galleries
func (g *Galleries) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	galleries, err := g.gs.ByUserID(user.ID)
	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	var vd views.Data
	vd.Yield = galleries
	g.IndexView.Render(w, r, vd)
}

// Show is used to show gallery.
//
// GET /galleries/:id
func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	vd := views.Data{Yield: gallery}
	g.ShowView.Render(w, r, vd)
}

// Edit is used to show the edit gallery view.
//
// GET /galleries/:id/edit
func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	var vd views.Data
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		vd.SetAlert(models.ErrNotFound)
		g.ShowView.Render(w, r, vd)
		return
	}

	vd.Yield = gallery
	g.EditView.Render(w, r, vd)
}

type UpdateGalleryForm struct {
	Title string `schema:"title,required"`
}

// Update is used to update a gallery.
//
// POST /galleries/:id/update
func (g *Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	var vd views.Data
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		vd.SetAlert(models.ErrNotFound)
		g.EditView.Render(w, r, vd)
		return
	}

	vd.Yield = gallery
	var form UpdateGalleryForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}

	gallery.Title = form.Title
	if err := g.gs.Update(gallery); err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}

	vd.Alert = &views.Alert{
		Level:   views.AlertLevelSucess,
		Message: "Gallery successfully updated!",
	}
	g.EditView.Render(w, r, vd)
}

// UpdateDelete is used to delete a gallery.
//
// POST /galleries/:id/delete
func (g *Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	var vd views.Data
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		vd.SetAlert(models.ErrNotFound)
		g.EditView.Render(w, r, vd)
		return
	}

	if err := g.gs.Delete(gallery.ID); err != nil {
		vd.Yield = gallery
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}

	url, err := g.r.Get(GalleriesIndexURL).URL()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	http.Redirect(w, r, url.Path, http.StatusMovedPermanently)
}

func (g *Galleries) galleryByID(w http.ResponseWriter, r *http.Request) (*models.Gallery, error) {
	var vd views.Data
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		g.l.Println("Error: converting id to int", err)
		vd.SetAlert(err)
		g.ShowView.Render(w, r, vd)
		return nil, err
	}

	gallery, err := g.gs.ByID(uint(id))
	if err != nil {
		g.l.Println("Error: fetching gallery", err)
		if err == models.ErrNotFound {
			vd.SetAlert(models.ErrNotFound)
			g.ShowView.Render(w, r, vd)
			return nil, err
		}

		vd.SetAlert(err)
		g.ShowView.Render(w, r, vd)
		return nil, err
	}

	return gallery, nil
}
