package server

import (
	"context"
	"fmt"
	"io/ioutil"
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

const staticDir = "/static"

var indexPage []byte

func NewServer(n tfdirectory.NodeService, f tfdirectory.FarmerService) *Server {
	r := mux.NewRouter().StrictSlash(true)
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

	// create API routes
	r.Use(handlers.CompressHandler)
	r.Use(func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, h)
	})
	r.Use(handlers.RecoveryHandler())

	apiRouter := s.router.PathPrefix("/api").Subrouter()
	apiRouter.Use(handlers.CORS())
	NewUserRouter(n, apiRouter.PathPrefix("/nodes").Subrouter())
	NewFarmerRouter(f, apiRouter.PathPrefix("/farmers").Subrouter())

	// var err error
	// indexPage, err = ioutil.ReadFile("static/index.html")
	// if err != nil {
	// 	log.Fatalf("cannot load index page: %v\n", err)
	// }

	// create static file route
	r.PathPrefix(staticDir).Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var err error
		indexPage, err = ioutil.ReadFile("static/index.html")
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot load index page: %v\n", err), http.StatusInternalServerError)
		}
		w.Write(indexPage)
	})
	return &s
}

// Start starts the server and make it listen to listen address
func (s *Server) Start(listen string) {
	log.Printf("Listening on %s\n", listen)

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
