package main

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/nicholasjackson/env"

	"github.com/karimla/webapp/controllers"
	"github.com/karimla/webapp/lib"
	"github.com/karimla/webapp/middleware"
	"github.com/karimla/webapp/models"
	"github.com/karimla/webapp/utils"
)

func main() {
	utils.Must(env.Parse())

	wg := &sync.WaitGroup{}
	l := lib.InitLog()

	services := models.NewServices()
	utils.Must(services.AutoMigrate())
	defer services.Close()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery)

	ar := middleware.NewAwaitRequest(wg)

	r := mux.NewRouter()
	r.Handle("/", ar.Apply(staticC.HomeView)).Methods(http.MethodGet)
	r.Handle("/contact", ar.Apply(staticC.ContactView)).Methods(http.MethodGet)
	r.Handle("/signup", ar.Apply(usersC.SignupView)).Methods(http.MethodGet)
	r.HandleFunc("/signup", ar.ApplyFn(usersC.Signup)).Methods(http.MethodPost)
	r.Handle("/login", ar.Apply(usersC.LoginView)).Methods(http.MethodGet)
	r.HandleFunc("/login", ar.ApplyFn(usersC.Login)).Methods(http.MethodPost)
	r.Handle("/galleries/new", ar.Apply(galleriesC.NewView)).Methods(http.MethodGet)
	r.HandleFunc("/galleries", ar.ApplyFn(galleriesC.Create)).Methods(http.MethodPost)

	s := lib.NewServer(l, wg, r)

	go s.Start()

	s.GracefulShutdown()
}
