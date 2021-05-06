package main

import (
	"sync"

	"github.com/nicholasjackson/env"

	"webapp/lib"
	"webapp/models"
	"webapp/routes"
	"webapp/utils"
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
