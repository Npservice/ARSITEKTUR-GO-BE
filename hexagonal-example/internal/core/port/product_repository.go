package port

import "hexagonal-example/internal/core/domain"

// Output port: dibutuhkan oleh domain/service, diimplementasikan oleh adapter keluar (outbound).
type ProductRepository interface {
	FindByID(id string) (*domain.Product, error)
	FindAll() ([]*domain.Product, error)
	Save(product *domain.Product) error
}
