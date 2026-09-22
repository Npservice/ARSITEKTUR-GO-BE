package order

import (
	"errors"

	"ddd-example/internal/domain/shared"
)

// Value Object: bagian dari agregat Order, tidak punya identitas sendiri.
type OrderItem struct {
	ProductID string
	Quantity  int
	UnitPrice shared.Money
}

func NewOrderItem(productID string, quantity int, unitPrice shared.Money) (OrderItem, error) {
	if productID == "" {
		return OrderItem{}, errors.New("productID is required")
	}
	if quantity <= 0 {
		return OrderItem{}, errors.New("quantity must be positive")
	}
	return OrderItem{ProductID: productID, Quantity: quantity, UnitPrice: unitPrice}, nil
}

func (i OrderItem) Subtotal() shared.Money {
	return i.UnitPrice.Multiply(float64(i.Quantity))
}
