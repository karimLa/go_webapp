package main

import (
	"sync"

	"github.com/nicholasjackson/env"

	"github.com/soramon0/webapp/lib"
	"github.com/soramon0/webapp/models"
	"github.com/soramon0/webapp/routes"
	"github.com/soramon0/webapp/utils"
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
