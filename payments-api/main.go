package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	srv := newServer()
	srv.start()
}

type server struct {
	router *mux.Router
}

func newServer() server {
	return server{
		router: mux.NewRouter(),
	}
}

// routes - all routes go here
func (s *server) routes() {
	s.router.HandleFunc("/payments", s.paymentHandler())
}

func (s *server) start() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: s.router,
	}

	s.routes()

	go func() {
		log.Printf("Listening on http://0.0.0.0\n")

		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-stop // Wait for interrupt

	log.Printf("shutting down ...\n")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Printf("shut down gracefully\n")
}
