package lib

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/soramon0/webapp/utils"
)

func NewServer(l *log.Logger, wg *sync.WaitGroup, r *mux.Router) *Server {
	return &Server{
		Server: http.Server{
			Addr:         utils.GetBindAdress(),
			Handler:      r,
			ErrorLog:     l,                 // set the logger for the server
			ReadTimeout:  5 * time.Second,   // max time to read request from the client
			WriteTimeout: 10 * time.Second,  // max time to write response to the client
			IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
		},
		l:  l,
		wg: wg,
	}
}

// Start will starts the server and a wait group
// for running jobs
func (s *Server) Start() {
	s.l.Println("Starting server on port 3000")

	if err := s.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			s.l.Println("waiting for running requests to finish")
			s.wg.Wait()
			s.l.Println("requests finished. exiting")
			return
		}

		s.l.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}

// GracefulShutdown listens for signal (SIGINT, SIGTERM, SIGHUB)
// and gracefully shutsdown the server when it recieves one
func (s *Server) GracefulShutdown() {
	// trap interupt, sigterm or sighub and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	// Block until a signal is received.
	sig := <-c
	s.l.Printf("Recieved %s, graceful shutdown...\n", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		s.l.Fatal(err)
	}
}

type Server struct {
	http.Server
	l  *log.Logger
	wg *sync.WaitGroup
}
