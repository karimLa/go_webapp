package routes

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/soramon0/webapp/controllers"
	"github.com/soramon0/webapp/middleware"
	"github.com/soramon0/webapp/models"
)

func Register(s *models.Services, wg *sync.WaitGroup) *mux.Router {
	r := mux.NewRouter()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(s.User)
	galleriesC := controllers.NewGalleries(s.Gallery, r)

	ar := middleware.NewAwaitRequest(wg)
	ru := middleware.NewRequireUser(s.User)

	r.Use(ar.Middleware)
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

	return r
}
