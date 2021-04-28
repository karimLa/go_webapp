package main

import (
	"sync"

	"github.com/nicholasjackson/env"

	"github.com/karimla/webapp/lib"
	"github.com/karimla/webapp/models"
	"github.com/karimla/webapp/routes"
	"github.com/karimla/webapp/utils"
)

func main() {
	utils.Must(env.Parse())

	wg := &sync.WaitGroup{}
	l := lib.InitLogger()

	services := models.NewServices()
	utils.Must(services.AutoMigrate())
	defer services.Close()

	r := routes.Register(services, wg)
	s := lib.NewServer(l, wg, r)

	go s.Start()

	s.GracefulShutdown()
}
