package routes

import (
	"log"
	"net/http"
	"sync"

	"webapp/controllers"
	"webapp/middleware"
	"webapp/models"

	"github.com/gorilla/mux"
)

func Register(s *models.Services, wg *sync.WaitGroup, l *log.Logger) *mux.Router {
	r := mux.NewRouter()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(s.User, r, l)
	galleriesC := controllers.NewGalleries(s.Gallery, s.Image, r, l)

	ar := middleware.NewAwaitRequest(wg)
	um := middleware.NewUser(s.User)
	ru := middleware.NewRequireUser(*um)

	r.Use(ar.Middleware)
	r.Use(um.Middleware)
	r.Handle("/", staticC.HomeView).Methods(http.MethodGet)
	r.Handle("/contact", staticC.ContactView).Methods(http.MethodGet)
	r.Handle("/signup", usersC.SignupView).Methods(http.MethodGet)
	r.HandleFunc("/signup", usersC.Signup).Methods(http.MethodPost)
	r.Handle("/login", usersC.LoginView).Methods(http.MethodGet)
	r.HandleFunc("/login", usersC.Login).Methods(http.MethodPost)
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods(http.MethodGet).Name(controllers.GalleryShowURL)

	authR := r.NewRoute().Subrouter()
	authR.Use(ru.Middleware)
	authR.Handle("/galleries/new", galleriesC.NewView).Methods(http.MethodGet)
	authR.HandleFunc("/galleries", galleriesC.Create).Methods(http.MethodPost)
	authR.HandleFunc("/galleries", galleriesC.Index).Methods(http.MethodGet).Name(controllers.GalleriesIndexURL)
	authR.HandleFunc("/galleries/{id:[0-9]+}/edit", galleriesC.Edit).Methods(http.MethodGet).Name(controllers.GalleryEditURL)
	authR.HandleFunc("/galleries/{id:[0-9]+}/update", galleriesC.Update).Methods(http.MethodPost)
	authR.HandleFunc("/galleries/{id:[0-9]+}/delete", galleriesC.Delete).Methods(http.MethodPost)
	authR.HandleFunc("/galleries/{id:[0-9]+}/images", galleriesC.ImageUpload).Methods(http.MethodPost)

	return r
}
