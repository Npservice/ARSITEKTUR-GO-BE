# Arsitektur Go — Contoh DDD & Hexagonal

Repo ini berisi dua contoh implementasi arsitektur di Go:

- [`ddd-example`](./ddd-example) — contoh Domain-Driven Design (aggregate `Order` & `Inventory`, use case pemesanan/pembatalan).
- [`hexagonal-example`](./hexagonal-example) — contoh Hexagonal Architecture (Ports & Adapters) untuk domain `Product`, dengan inbound adapter HTTP dan CLI.

## Persyaratan

- Go 1.22 atau lebih baru — cek dengan `go version`.
- Git (untuk clone repo).

## Setup Awal

```bash
git clone https://github.com/Npservice/ARSITEKTUR-GO-BE.git
cd ARSITEKTUR-GO-BE
```

Setiap folder (`ddd-example`, `hexagonal-example`) adalah Go module terpisah (masing-masing punya `go.mod` sendiri), jadi perintah `go run`/`go build` harus dijalankan dari dalam folder modul tersebut.

## 1. `ddd-example`

Contoh ini menjalankan skenario pemesanan langsung di `main.go` (tanpa server), lalu mencetak hasilnya ke terminal.

```bash
cd ddd-example
go run ./cmd/app
```

Yang akan terjadi:

- Seed stok awal untuk `product-1` dan `product-2`.
- Membuat order pertama (berhasil) lalu mengurangi stok.
- Mencoba order kedua yang melebihi stok (ditolak, membuktikan invariant domain).
- Membatalkan order pertama dan stok dikembalikan.

## 2. `hexagonal-example`

Contoh ini punya dua composition root: server HTTP (`cmd/app`) dan CLI (`cmd/cli`), keduanya memakai `core/service.ProductService` yang sama.

### Menjalankan server HTTP

```bash
cd hexagonal-example
go run ./cmd/app
```

Konfigurasi via environment variable (opsional, ada default):

| Env Var      | Default          | Keterangan                          |
|--------------|------------------|--------------------------------------|
| `HTTP_PORT`  | `8080`           | Port server HTTP                     |
| `STORE_KIND` | `memory`         | Storage adapter: `memory` atau `file`|
| `STORE_FILE` | `products.json`  | Path file saat `STORE_KIND=file`     |

Atau override lewat flag CLI, contoh pakai file store:

```bash
go run ./cmd/app --store=file --file=products.json
```

Endpoint yang tersedia:

- `GET /products` — daftar produk.
- `POST /products` — buat produk baru, body JSON: `{"id":"p1","name":"Produk A","price":10000,"stock":5}`.
- `POST /products/purchase` — beli produk, body JSON: `{"id":"p1","quantity":2}`.

Contoh test cepat:

```bash
curl -X POST localhost:8080/products -d '{"id":"p1","name":"Produk A","price":10000,"stock":5}'
curl localhost:8080/products
curl -X POST localhost:8080/products/purchase -d '{"id":"p1","quantity":2}'
```

### Menjalankan CLI

CLI ini selalu pakai file store (`products.json`) supaya datanya bisa konsisten dengan server saat dijalankan dengan `--store=file`.

```bash
cd hexagonal-example
go run ./cmd/cli create p1 "Produk A" 10000 5
go run ./cmd/cli list
go run ./cmd/cli purchase p1 2
```

## Build Binary (opsional)

```bash
cd ddd-example && go build -o ddd-app ./cmd/app
cd hexagonal-example && go build -o hex-app ./cmd/app && go build -o hex-cli ./cmd/cli
```

## Testing

Belum ada test otomatis di kedua modul ini. Untuk menjalankan test bila ditambahkan nanti:

```bash
go test ./...
```
