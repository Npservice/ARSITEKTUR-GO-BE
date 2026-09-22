package main

import (
	"fmt"
	"log"

	appOrder "ddd-example/internal/application/order"
	domaininventory "ddd-example/internal/domain/inventory"
	"ddd-example/internal/infrastructure/persistence"
)

func main() {
	orderRepo := persistence.NewInMemoryOrderRepository()
	inventoryRepo := persistence.NewInMemoryInventoryRepository()

	// Seed stok awal untuk aggregate Inventory.
	stock1, _ := domaininventory.NewInventory("product-1", 5)
	stock2, _ := domaininventory.NewInventory("product-2", 2)
	_ = inventoryRepo.Save(stock1)
	_ = inventoryRepo.Save(stock2)

	placeOrder := appOrder.NewPlaceOrderUseCase(orderRepo, inventoryRepo)
	cancelOrder := appOrder.NewCancelOrderUseCase(orderRepo, inventoryRepo)

	result, err := placeOrder.Execute(appOrder.PlaceOrderInput{
		CustomerID: "customer-1",
		Items: []appOrder.PlaceOrderItemInput{
			{ProductID: "product-1", Quantity: 2, UnitPrice: 15000, Currency: "IDR"},
			{ProductID: "product-2", Quantity: 1, UnitPrice: 25000, Currency: "IDR"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	total, _ := result.Total()
	fmt.Printf("Order %s placed with status %s, total %.2f %s\n",
		result.ID(), result.Status(), total.Amount(), total.Currency())

	afterReserve, _ := inventoryRepo.FindByProductID("product-1")
	fmt.Printf("product-1 stock after reserve: %d\n", afterReserve.AvailableStock())

	// Contoh invariant ditolak: order yang sudah PLACED tidak boleh melebihi stok
	// yang tersisa. Ini membuktikan Inventory menjaga aturannya sendiri.
	_, err = placeOrder.Execute(appOrder.PlaceOrderInput{
		CustomerID: "customer-2",
		Items: []appOrder.PlaceOrderItemInput{
			{ProductID: "product-2", Quantity: 10, UnitPrice: 25000, Currency: "IDR"},
		},
	})
	fmt.Printf("second order rejected as expected: %v\n", err)

	cancelled, err := cancelOrder.Execute(result.ID(), "customer requested refund")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Order %s status after cancel: %s\n", cancelled.ID(), cancelled.Status())

	afterRelease, _ := inventoryRepo.FindByProductID("product-1")
	fmt.Printf("product-1 stock after cancel (released back): %d\n", afterRelease.AvailableStock())
}
