package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/lohanguedes/movie-microservices/movie/internal/controller/movie"
	metadatagateway "github.com/lohanguedes/movie-microservices/movie/internal/gateway/metadata/http"
	ratinggateway "github.com/lohanguedes/movie-microservices/movie/internal/gateway/rating/http"
	httphandler "github.com/lohanguedes/movie-microservices/movie/internal/handler/http"
	"github.com/lohanguedes/movie-microservices/pkg/discovery"
	"github.com/lohanguedes/movie-microservices/pkg/discovery/consul"
)

var serviceName = "movie"

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "API handler port")
	flag.Parse()

	log.Printf("Staring the movie service on port %d\n", port)

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

	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)
	ctrl := movie.New(ratingGateway, metadataGateway)
	h := httphandler.New(ctrl)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
		panic(err)
	}
}
