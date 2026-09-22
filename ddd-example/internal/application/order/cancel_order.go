package order

import (
	"fmt"

	domaininventory "ddd-example/internal/domain/inventory"
	domainorder "ddd-example/internal/domain/order"
)

type CancelOrderUseCase struct {
	orderRepository     domainorder.Repository
	inventoryRepository domaininventory.Repository
}

func NewCancelOrderUseCase(orderRepository domainorder.Repository, inventoryRepository domaininventory.Repository) *CancelOrderUseCase {
	return &CancelOrderUseCase{orderRepository: orderRepository, inventoryRepository: inventoryRepository}
}

func (uc *CancelOrderUseCase) Execute(orderID, reason string) (*domainorder.Order, error) {
	existingOrder, err := uc.orderRepository.FindByID(orderID)
	if err != nil {
		return nil, err
	}
	if existingOrder == nil {
		return nil, fmt.Errorf("order %s not found", orderID)
	}

	// Invariant pembatalan (status DRAFT/CANCELLED ditolak) ditegakkan di dalam aggregate.
	if err := existingOrder.Cancel(reason); err != nil {
		return nil, err
	}

	// Kompensasi: kembalikan stok yang sebelumnya direservasi ke aggregate Inventory.
	for _, item := range existingOrder.Items() {
		inv, err := uc.inventoryRepository.FindByProductID(item.ProductID)
		if err != nil {
			return nil, err
		}
		if inv == nil {
			continue
		}
		if err := inv.Release(item.Quantity); err != nil {
			return nil, err
		}
		if err := uc.inventoryRepository.Save(inv); err != nil {
			return nil, err
		}
	}

	if err := uc.orderRepository.Save(existingOrder); err != nil {
		return nil, err
	}

	for range existingOrder.PullEvents() {
		// no-op: contoh sederhana, tidak ada event handler nyata.
	}

	return existingOrder, nil
}
