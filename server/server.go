package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/zaibon/tfdirectory"
)

type Server struct {
	srv    *http.Server
	router *mux.Router
}

func NewServer(n tfdirectory.NodeService, f tfdirectory.FarmerService) *Server {
	r := mux.NewRouter()
	srv := &http.Server{
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}
	s := Server{
		router: r,
		srv:    srv,
	}
	NewUserRouter(n, s.newSubrouter("/nodes"))
	NewFarmerRouter(f, s.newSubrouter("/farmers"))
	return &s
}

// Start starts the server and make it listen to listen address
func (s *Server) Start(listen string) {
	log.Printf("Listening on %s\n", listen)

	s.router.Use(handlers.CompressHandler)
	s.router.Use(func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, h)
	})
	s.router.Use(handlers.RecoveryHandler())

	s.srv.Addr = listen
	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			log.Println("http.ListenAndServe: ", err)
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) newSubrouter(path string) *mux.Router {
	return s.router.PathPrefix(path).Subrouter()
}
