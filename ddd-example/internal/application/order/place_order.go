package order

import (
	"fmt"

	domaininventory "ddd-example/internal/domain/inventory"
	domainorder "ddd-example/internal/domain/order"
	"ddd-example/internal/domain/shared"
)

type PlaceOrderItemInput struct {
	ProductID string
	Quantity  int
	UnitPrice float64
	Currency  string
}

type PlaceOrderInput struct {
	CustomerID string
	Items      []PlaceOrderItemInput
}

// Application Service: mengorkestrasi DUA aggregate (Order dan Inventory).
// Ini letak yang benar untuk koordinasi lintas aggregate di DDD — bukan business
// rule (itu tetap di dalam masing-masing aggregate lewat Reserve()/AddItem()/Place()),
// melainkan urutan transaksi: reservasi stok dulu per item, baru order di-place.
// Kalau salah satu reservasi gagal, reservasi yang sudah jalan di-rollback (compensasi),
// karena dua aggregate = dua unit-of-work terpisah, tidak bisa satu transaksi DB.
type PlaceOrderUseCase struct {
	orderRepository     domainorder.Repository
	inventoryRepository domaininventory.Repository
}

func NewPlaceOrderUseCase(orderRepository domainorder.Repository, inventoryRepository domaininventory.Repository) *PlaceOrderUseCase {
	return &PlaceOrderUseCase{orderRepository: orderRepository, inventoryRepository: inventoryRepository}
}

func (uc *PlaceOrderUseCase) Execute(input PlaceOrderInput) (*domainorder.Order, error) {
	newOrder, err := domainorder.NewOrder(uc.orderRepository.NextID(), input.CustomerID)
	if err != nil {
		return nil, err
	}

	reserved := make([]PlaceOrderItemInput, 0, len(input.Items))
	rollback := func() {
		for _, item := range reserved {
			if inv, ferr := uc.inventoryRepository.FindByProductID(item.ProductID); ferr == nil && inv != nil {
				_ = inv.Release(item.Quantity)
				_ = uc.inventoryRepository.Save(inv)
			}
		}
	}

	for _, item := range input.Items {
		inv, err := uc.inventoryRepository.FindByProductID(item.ProductID)
		if err != nil {
			rollback()
			return nil, err
		}
		if inv == nil {
			rollback()
			return nil, fmt.Errorf("no inventory record for product %s", item.ProductID)
		}
		if err := inv.Reserve(item.Quantity); err != nil {
			rollback()
			return nil, fmt.Errorf("product %s: %w", item.ProductID, err)
		}
		if err := uc.inventoryRepository.Save(inv); err != nil {
			rollback()
			return nil, err
		}
		reserved = append(reserved, item)

		unitPrice, err := shared.NewMoney(item.UnitPrice, item.Currency)
		if err != nil {
			rollback()
			return nil, err
		}
		if err := newOrder.AddItem(item.ProductID, item.Quantity, unitPrice); err != nil {
			rollback()
			return nil, err
		}
	}

	if err := newOrder.Place(); err != nil {
		rollback()
		return nil, err
	}

	if err := uc.orderRepository.Save(newOrder); err != nil {
		rollback()
		return nil, err
	}

	// Di aplikasi nyata, event ini dipublikasikan ke message bus/event dispatcher.
	for range newOrder.PullEvents() {
		// no-op: contoh sederhana, tidak ada event handler nyata.
	}

	return newOrder, nil
}
