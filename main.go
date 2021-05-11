package main

import (
	"sync"

	"github.com/nicholasjackson/env"

	"soramon0/webapp/lib"
	"soramon0/webapp/models"
	"soramon0/webapp/routes"
	"soramon0/webapp/utils"
)

func main() {
	utils.Must(env.Parse())

	wg := &sync.WaitGroup{}
	l := lib.InitLogger()

	services := models.NewServices()
	utils.Must(services.AutoMigrate())
	defer services.Close()

	r := routes.Register(services, wg, l)
	s := lib.NewServer(l, wg, r)

	go s.Start()

	s.GracefulShutdown()
}
