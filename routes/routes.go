package routes

import (
	"log"
	"net/http"
	"sync"

	"soramon0/webapp/controllers"
	"soramon0/webapp/middleware"
	"soramon0/webapp/models"

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

	// Serving images
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))

	// Serving assets
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	baseR := r.NewRoute().Subrouter()
	baseR.Use(ar.Middleware)
	baseR.Use(um.Middleware)
	baseR.Handle("/", staticC.HomeView).Methods(http.MethodGet)
	baseR.Handle("/contact", staticC.ContactView).Methods(http.MethodGet)
	baseR.Handle("/signup", usersC.SignupView).Methods(http.MethodGet)
	baseR.HandleFunc("/signup", usersC.Signup).Methods(http.MethodPost)
	baseR.Handle("/login", usersC.LoginView).Methods(http.MethodGet)
	baseR.HandleFunc("/login", usersC.Login).Methods(http.MethodPost)
	baseR.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods(http.MethodGet).Name(controllers.GalleryShowURL)

	authR := baseR.NewRoute().Subrouter()
	authR.Use(ru.Middleware)
	authR.Handle("/galleries/new", galleriesC.NewView).Methods(http.MethodGet)
	authR.HandleFunc("/galleries", galleriesC.Create).Methods(http.MethodPost)
	authR.HandleFunc("/galleries", galleriesC.Index).Methods(http.MethodGet).Name(controllers.GalleriesIndexURL)
	authR.HandleFunc("/galleries/{id:[0-9]+}/edit", galleriesC.Edit).Methods(http.MethodGet).Name(controllers.GalleryEditURL)
	authR.HandleFunc("/galleries/{id:[0-9]+}/update", galleriesC.Update).Methods(http.MethodPost)
	authR.HandleFunc("/galleries/{id:[0-9]+}/delete", galleriesC.Delete).Methods(http.MethodPost)
	authR.HandleFunc("/galleries/{id:[0-9]+}/images", galleriesC.ImageUpload).Methods(http.MethodPost)
	authR.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete", galleriesC.ImageDelete).Methods(http.MethodPost)

	return r
}
