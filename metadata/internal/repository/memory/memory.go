package memory

import (
	"context"
	"sync"

	"github.com/lohanguedes/movie-microservices/metadata/internal/repository"
	"github.com/lohanguedes/movie-microservices/metadata/pkg/model"
)

// Repository defines a memory movie metadata repository.
type Repository struct {
	sync.RWMutex
	data map[string]*model.Metadata
}

// New creates a new in-memory respository.
func New() *Repository {
	return &Repository{
		data: map[string]*model.Metadata{},
	}
}

// Get retrieves movie metadata my movie id.
func (r *Repository) Get(ctx context.Context, id string) (*model.Metadata, error) {
	r.RLock()
	defer r.RUnlock()
	m, ok := r.data[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return m, nil
}

// Put adds movie metadata for a given movie id.
func (r *Repository) Put(ctx context.Context, id string, metadata *model.Metadata) error {
	r.Lock()
	defer r.Unlock()
	r.data[id] = metadata
	return nil
}
