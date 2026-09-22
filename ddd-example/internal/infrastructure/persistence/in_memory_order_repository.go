package persistence

import (
	"fmt"
	"sync"

	domainorder "ddd-example/internal/domain/order"
)

// Infrastructure layer: implementasi konkret dari domainorder.Repository.
type InMemoryOrderRepository struct {
	mu      sync.RWMutex
	data    map[string]*domainorder.Order
	counter int
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{data: make(map[string]*domainorder.Order)}
}

func (r *InMemoryOrderRepository) NextID() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counter++
	return fmt.Sprintf("order-%d", r.counter)
}

func (r *InMemoryOrderRepository) FindByID(id string) (*domainorder.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	order, ok := r.data[id]
	if !ok {
		return nil, nil
	}
	return order, nil
}

func (r *InMemoryOrderRepository) Save(order *domainorder.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[order.ID()] = order
	return nil
}
