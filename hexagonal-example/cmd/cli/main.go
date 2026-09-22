package main

import (
	"log"
	"os"

	incli "hexagonal-example/internal/adapter/in/cli"
	"hexagonal-example/internal/adapter/out/file"
	"hexagonal-example/internal/core/service"
)

// Composition root kedua: inbound adapter-nya CLI, bukan HTTP, tapi tetap
// memakai core/service.ProductService yang sama persis. Dipakai bersama
// file store supaya data konsisten dengan `cmd/app --store=file`.
func main() {
	repo := file.NewProductRepository("products.json")
	productService := service.NewProductService(repo)
	handler := incli.NewProductCLI(productService)

	if err := handler.Run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
