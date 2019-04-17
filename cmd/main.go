package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/zaibon/tfdirectory/mongo"
	"github.com/zaibon/tfdirectory/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ms, err := mongo.NewSession(ctx, "mongodb://localhost:27017")
	if err != nil {
		log.Fatalln("unable to connect to mongodb")
	}
	defer ms.Close(ctx)

	us := mongo.NewNodeService(ms, "tfdirectory", "node")
	fs := mongo.NewFarmerService(ms, "tfdirectory", "farmer")
	s := server.NewServer(us, fs)

	// start the http server
	s.Start(":8081")

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt, os.Kill)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for the http server to stop.
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	go func() {
		log.Println("shutting down")
		s.Shutdown(ctx)
	}()

	<-ctx.Done()
	log.Println("server is shutdown")

	os.Exit(0)
}
