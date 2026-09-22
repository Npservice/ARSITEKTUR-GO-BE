package inventory

import "errors"

// Aggregate Root terpisah dari Order. DDD menegaskan: satu transaksi = satu aggregate.
// Order TIDAK boleh memegang referensi langsung ke Inventory; keduanya hanya
// terhubung lewat ProductID (referensi by identity), dikoordinasikan lewat
// application layer atau domain service.
type Inventory struct {
	productID      string
	availableStock int
}

func NewInventory(productID string, availableStock int) (*Inventory, error) {
	if productID == "" {
		return nil, errors.New("productID is required")
	}
	if availableStock < 0 {
		return nil, errors.New("availableStock cannot be negative")
	}
	return &Inventory{productID: productID, availableStock: availableStock}, nil
}

func (i *Inventory) ProductID() string { return i.productID }
func (i *Inventory) AvailableStock() int { return i.availableStock }

// Reserve menegakkan invariant miliknya sendiri: stok tidak boleh negatif.
// Invariant ini murni tanggung jawab aggregate Inventory, bukan Order.
func (i *Inventory) Reserve(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	if quantity > i.availableStock {
		return errors.New("insufficient stock")
	}
	i.availableStock -= quantity
	return nil
}

func (i *Inventory) Release(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	i.availableStock += quantity
	return nil
}
