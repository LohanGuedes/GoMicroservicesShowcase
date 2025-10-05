package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/lohanguedes/movie-microservices/metadata/internal/controller/metadata"
	httphandler "github.com/lohanguedes/movie-microservices/metadata/internal/handler/http"
	"github.com/lohanguedes/movie-microservices/metadata/internal/repository/memory"
	"github.com/lohanguedes/movie-microservices/pkg/discovery"
	"github.com/lohanguedes/movie-microservices/pkg/discovery/consul"
)

var serviceName = "metadata"

func main() {
	var port int
	flag.IntVar(&port, "port", 8082, "API handler port")
	flag.Parse()

	log.Printf("Staring the movie metadata service on port %d\n", port)
	registry, err := consul.NewResgistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}
	go func() {
		for {
			if err := registry.ReportHealhyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state:", err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)
	repo := memory.New()
	ctrl := metadata.New(repo)
	h := httphandler.New(ctrl)

	http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
