package service

import (
	"fmt"

	"hexagonal-example/internal/core/domain"
	"hexagonal-example/internal/core/port"
)

type productService struct {
	// Bergantung hanya pada port (abstraksi), bukan pada adapter konkret.
	repo port.ProductRepository
}

func NewProductService(repo port.ProductRepository) port.ProductService {
	return &productService{repo: repo}
}

func (s *productService) CreateProduct(id, name string, price float64, stock int) (*domain.Product, error) {
	product, err := domain.NewProduct(id, name, price, stock)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Save(product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *productService) ListProducts() ([]*domain.Product, error) {
	return s.repo.FindAll()
}

func (s *productService) Purchase(id string, quantity int) (*domain.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, fmt.Errorf("product %s not found", id)
	}
	if err := product.ReduceStock(quantity); err != nil {
		return nil, err
	}
	if err := s.repo.Save(product); err != nil {
		return nil, err
	}
	return product, nil
}
