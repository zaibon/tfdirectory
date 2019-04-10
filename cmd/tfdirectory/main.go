package main

import (
	"context"
	"log"

	"github.com/zaibon/tfdirectory/mongo"
	"github.com/zaibon/tfdirectory/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ms, err := mongo.NewSession(ctx, "127.0.0.1:27017")
	if err != nil {
		log.Fatalln("unable to connect to mongodb")
	}
	defer ms.Close(ctx)

	ns := mongo.NewNodeService(ms, "tfdirectory", "node")
	s := server.NewServer(ns)

	s.Start()
}
