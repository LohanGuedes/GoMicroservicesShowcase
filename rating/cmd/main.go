package main

import (
	"log"
	"net/http"

	httphandler "github.com/lohanguedes/movie-microservices/metadata/internal/handler/http"
	"github.com/lohanguedes/movie-microservices/rating/internal/controller/rating"
	"github.com/lohanguedes/movie-microservices/rating/internal/repository/memory"
)

func main() {
	log.Println("Starting the rating service")
	repo := memory.New()
	ctrl := rating.New(repo)
	h := httphandler.New(ctrl)
	http.Handle("/rating", http.HandlerFunc(h.Handle))
	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}
}
