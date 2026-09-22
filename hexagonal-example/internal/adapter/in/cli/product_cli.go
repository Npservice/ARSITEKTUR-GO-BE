package cli

import (
	"fmt"
	"strconv"

	"hexagonal-example/internal/core/port"
)

// Inbound adapter kedua: menerjemahkan argumen CLI menjadi pemanggilan use case (input port
// yang sama persis dengan yang dipakai adapter HTTP). Core/service tidak tahu-menahu ini
// dipanggil dari CLI atau HTTP.
type ProductCLI struct {
	service port.ProductService
}

func NewProductCLI(service port.ProductService) *ProductCLI {
	return &ProductCLI{service: service}
}

// Run menerima args ala: create <id> <name> <price> <stock> | list | purchase <id> <qty>
func (c *ProductCLI) Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: create|list|purchase [args]")
	}

	switch args[0] {
	case "create":
		if len(args) != 5 {
			return fmt.Errorf("usage: create <id> <name> <price> <stock>")
		}
		price, err := strconv.ParseFloat(args[3], 64)
		if err != nil {
			return err
		}
		stock, err := strconv.Atoi(args[4])
		if err != nil {
			return err
		}
		product, err := c.service.CreateProduct(args[1], args[2], price, stock)
		if err != nil {
			return err
		}
		fmt.Printf("created: %+v\n", product)

	case "list":
		products, err := c.service.ListProducts()
		if err != nil {
			return err
		}
		for _, p := range products {
			fmt.Printf("%+v\n", p)
		}

	case "purchase":
		if len(args) != 3 {
			return fmt.Errorf("usage: purchase <id> <qty>")
		}
		qty, err := strconv.Atoi(args[2])
		if err != nil {
			return err
		}
		product, err := c.service.Purchase(args[1], qty)
		if err != nil {
			return err
		}
		fmt.Printf("purchased: %+v\n", product)

	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}

	return nil
}
