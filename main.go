package main

import (
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
	db := lib.InitDB()
	us := models.NewUserService(db)

	staticC := controllers.NewStatic(wg)
	usersC := controllers.NewUsers(wg, us)

	r := mux.NewRouter()
	r.Handle("/", staticC.HomeView).Methods("GET")
	r.Handle("/contact", staticC.ContactView).Methods("GET")
	r.Handle("/signup", usersC.SignupView).Methods("GET")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/signup", usersC.Signup).Methods("POST")
	r.HandleFunc("/login", usersC.Login).Methods("POST")

	s := lib.NewServer(l, wg, r)

	go s.Start()

	s.GracefulShutdown()
}
