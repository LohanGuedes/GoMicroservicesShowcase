package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/lohanguedes/movie-microservices/metadata/internal/controller/metadata"
	"github.com/lohanguedes/movie-microservices/metadata/internal/repository"
)

type Handler struct {
	ctrl *metadata.Controller
}

func New(ctrl *metadata.Controller) *Handler {
	return &Handler{ctrl}
}

func (h *Handler) GetMetadata(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m, err := h.ctrl.Get(r.Context(), id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("repository get error for movie %s: %v\n", id, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Printf("Response encoding error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
