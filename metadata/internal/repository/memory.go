package memory

import (
	"context"
	"sync"

	"github.com/lohanguedes/movie-microservices/metadata/pkg/model"
)

type Repository struct {
	sync.RWMutex
	data map[string]*model.Metadata
}

func NewRepository() *Repository {
	return &Repository{
		data: map[string]*model.Metadata{},
	}
}

func (r *Repository) Get(_ context.Context, id string) (*model.Metadata, error) {
}
