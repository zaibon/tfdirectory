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

	ms, err := mongo.NewSession(ctx, "mongodb://localhost:27017")
	if err != nil {
		log.Fatalln("unable to connect to mongodb")
	}
	defer ms.Close(ctx)

	us := mongo.NewNodeService(ms, "tfdirectory", "node")
	fs := mongo.NewFarmerService(ms, "tfdirectory", "farmer")
	s := server.NewServer(us, fs)

	s.Start(":8081")
}
