package file

import (
	"encoding/json"
	"os"
	"sync"

	"hexagonal-example/internal/core/domain"
	"hexagonal-example/internal/core/port"
)

// Outbound adapter kedua untuk port yang sama (port.ProductRepository).
// Membuktikan core/service sama sekali tidak berubah walau backing store
// diganti dari memori ke file JSON.
type productRepository struct {
	mu   sync.Mutex
	path string
}

func NewProductRepository(path string) port.ProductRepository {
	return &productRepository{path: path}
}

func (r *productRepository) load() (map[string]*domain.Product, error) {
	data := make(map[string]*domain.Product)
	bytes, err := os.ReadFile(r.path)
	if os.IsNotExist(err) {
		return data, nil
	}
	if err != nil {
		return nil, err
	}
	if len(bytes) == 0 {
		return data, nil
	}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (r *productRepository) persist(data map[string]*domain.Product) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.path, bytes, 0644)
}

func (r *productRepository) FindByID(id string) (*domain.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	data, err := r.load()
	if err != nil {
		return nil, err
	}
	return data[id], nil
}

func (r *productRepository) FindAll() ([]*domain.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	data, err := r.load()
	if err != nil {
		return nil, err
	}
	products := make([]*domain.Product, 0, len(data))
	for _, p := range data {
		products = append(products, p)
	}
	return products, nil
}

func (r *productRepository) Save(product *domain.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	data, err := r.load()
	if err != nil {
		return err
	}
	data[product.ID] = product
	return r.persist(data)
}
