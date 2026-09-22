package memory

import (
	"sync"

	"hexagonal-example/internal/core/domain"
	"hexagonal-example/internal/core/port"
)

// Outbound adapter: implementasi port.ProductRepository dengan penyimpanan di memori.
// Bisa diganti dengan adapter lain (mis. Postgres) tanpa menyentuh domain/service.
type productRepository struct {
	mu   sync.RWMutex
	data map[string]*domain.Product
}

func NewProductRepository() port.ProductRepository {
	return &productRepository{data: make(map[string]*domain.Product)}
}

func (r *productRepository) FindByID(id string) (*domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.data[id], nil
}

func (r *productRepository) FindAll() ([]*domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	products := make([]*domain.Product, 0, len(r.data))
	for _, p := range r.data {
		products = append(products, p)
	}
	return products, nil
}

func (r *productRepository) Save(product *domain.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[product.ID] = product
	return nil
}
