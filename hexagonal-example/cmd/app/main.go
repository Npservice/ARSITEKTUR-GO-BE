package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	inhttp "hexagonal-example/internal/adapter/in/http"
	"hexagonal-example/internal/adapter/out/file"
	"hexagonal-example/internal/adapter/out/memory"
	"hexagonal-example/internal/config"
	"hexagonal-example/internal/core/port"
	"hexagonal-example/internal/core/service"
)

// Composition root: satu-satunya tempat config (env var), adapter konkret,
// dan port domain "ditemukan" satu sama lain. core/service dan adapter HTTP
// tidak pernah mengimpor package config ini secara langsung.
func main() {
	cfg := config.Load() // baca env var: HTTP_PORT, STORE_KIND, STORE_FILE

	// Flag CLI dijadikan override tipis di atas config env, defaultnya ikut cfg.
	store := flag.String("store", cfg.StoreKind, "storage adapter: memory|file")
	filePath := flag.String("file", cfg.FilePath, "path used when --store=file")
	flag.Parse()

	var repo port.ProductRepository
	switch *store {
	case "file":
		repo = file.NewProductRepository(*filePath)
	default:
		repo = memory.NewProductRepository()
	}

	// productService dan handler hanya menerima port/nilai polos (string, interface),
	// tidak pernah tahu nilai itu tadinya datang dari env var HTTP_PORT.
	productService := service.NewProductService(repo)
	handler := inhttp.NewProductHandler(productService)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	addr := fmt.Sprintf(":%s", cfg.HTTPPort)
	log.Printf("hexagonal-example (store=%s) listening on %s", *store, addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
