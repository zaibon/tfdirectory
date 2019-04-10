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

func NewServer(u tfdirectory.NodeService) *Server {
	s := Server{router: mux.NewRouter()}
	NewUserRouter(u, s.newSubrouter("/nodes"))
	return &s
}

func (s *Server) Start() {
	log.Println("Listening on port 8081")
	if err := http.ListenAndServe(":8081", handlers.LoggingHandler(os.Stdout, s.router)); err != nil {
		log.Fatal("http.ListenAndServe: ", err)
	}
}

func (s *Server) newSubrouter(path string) *mux.Router {
	return s.router.PathPrefix(path).Subrouter()
}
