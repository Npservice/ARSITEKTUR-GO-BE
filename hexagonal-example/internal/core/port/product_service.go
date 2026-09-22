package port

import "hexagonal-example/internal/core/domain"

// Input port: use case yang dipanggil oleh adapter masuk (inbound), mis. HTTP handler.
type ProductService interface {
	CreateProduct(id, name string, price float64, stock int) (*domain.Product, error)
	ListProducts() ([]*domain.Product, error)
	Purchase(id string, quantity int) (*domain.Product, error)
}
