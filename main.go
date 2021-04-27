package main

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/nicholasjackson/env"

	"github.com/karimla/webapp/controllers"
	"github.com/karimla/webapp/lib"
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

	staticC := controllers.NewStatic(wg)
	usersC := controllers.NewUsers(wg, services.User)
	galleriesC := controllers.NewGalleries(wg, services.Gallery)

	r := mux.NewRouter()
	r.Handle("/", staticC.HomeView).Methods(http.MethodGet)
	r.Handle("/contact", staticC.ContactView).Methods(http.MethodGet)
	r.Handle("/signup", usersC.SignupView).Methods(http.MethodGet)
	r.HandleFunc("/signup", usersC.Signup).Methods(http.MethodPost)
	r.Handle("/login", usersC.LoginView).Methods(http.MethodGet)
	r.HandleFunc("/login", usersC.Login).Methods(http.MethodPost)
	r.Handle("/galleries/new", galleriesC.NewView).Methods(http.MethodGet)

	s := lib.NewServer(l, wg, r)

	go s.Start()

	s.GracefulShutdown()
}
