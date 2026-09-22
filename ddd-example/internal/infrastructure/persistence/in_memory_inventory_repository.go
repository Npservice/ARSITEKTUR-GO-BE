package persistence

import (
	"sync"

	domaininventory "ddd-example/internal/domain/inventory"
)

type InMemoryInventoryRepository struct {
	mu   sync.RWMutex
	data map[string]*domaininventory.Inventory
}

func NewInMemoryInventoryRepository() *InMemoryInventoryRepository {
	return &InMemoryInventoryRepository{data: make(map[string]*domaininventory.Inventory)}
}

func (r *InMemoryInventoryRepository) FindByProductID(productID string) (*domaininventory.Inventory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.data[productID], nil
}

func (r *InMemoryInventoryRepository) Save(inv *domaininventory.Inventory) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[inv.ProductID()] = inv
	return nil
}
