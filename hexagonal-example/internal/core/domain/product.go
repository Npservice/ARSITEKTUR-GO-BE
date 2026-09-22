package domain

import "errors"

type Product struct {
	ID    string
	Name  string
	Price float64
	Stock int
}

func NewProduct(id, name string, price float64, stock int) (*Product, error) {
	if price < 0 {
		return nil, errors.New("price cannot be negative")
	}
	if stock < 0 {
		return nil, errors.New("stock cannot be negative")
	}
	return &Product{ID: id, Name: name, Price: price, Stock: stock}, nil
}

func (p *Product) ReduceStock(quantity int) error {
	if quantity > p.Stock {
		return errors.New("insufficient stock")
	}
	p.Stock -= quantity
	return nil
}
