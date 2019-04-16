package server

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/zaibon/tfdirectory"
)

type Server struct {
	router *mux.Router
}

func NewServer(n tfdirectory.NodeService, f tfdirectory.FarmerService) *Server {
	s := Server{router: mux.NewRouter()}
	NewUserRouter(n, s.newSubrouter("/nodes"))
	NewFarmerRouter(f, s.newSubrouter("/farmers"))
	return &s
}

func (s *Server) Start(listen string) {
	log.Printf("Listening on %s\n", listen)
	if err := http.ListenAndServe(listen, handlers.LoggingHandler(os.Stdout, s.router)); err != nil {
		log.Fatal("http.ListenAndServe: ", err)
	}
}

func (s *Server) newSubrouter(path string) *mux.Router {
	return s.router.PathPrefix(path).Subrouter()
}
