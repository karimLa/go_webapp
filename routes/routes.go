package routes

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/karimla/webapp/controllers"
	"github.com/karimla/webapp/middleware"
	"github.com/karimla/webapp/models"
)

func Register(s *models.Services, wg *sync.WaitGroup) *mux.Router {
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(s.User)
	galleriesC := controllers.NewGalleries(s.Gallery)

	ar := middleware.NewAwaitRequest(wg)
	ru := middleware.NewRequireUser(s.User)

	r := mux.NewRouter()
	r.Use(ar.Middleware)
	r.Handle("/", staticC.HomeView).Methods(http.MethodGet)
	r.Handle("/contact", staticC.ContactView).Methods(http.MethodGet)
	r.Handle("/signup", usersC.SignupView).Methods(http.MethodGet)
	r.HandleFunc("/signup", usersC.Signup).Methods(http.MethodPost)
	r.Handle("/login", usersC.LoginView).Methods(http.MethodGet)
	r.HandleFunc("/login", usersC.Login).Methods(http.MethodPost)
	r.Handle("/galleries/new", ru.Apply(galleriesC.NewView)).Methods(http.MethodGet)
	r.HandleFunc("/galleries", ru.ApplyFn(galleriesC.Create)).Methods(http.MethodPost)

	return r
}
