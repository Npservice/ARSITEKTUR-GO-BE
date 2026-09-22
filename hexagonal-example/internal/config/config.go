package config

import "os"

// Ini INFRASTRUKTUR, bukan core — isinya baca env var (I/O ke luar proses).
// Karena itu package ini HANYA boleh diimport dari composition root (cmd/*),
// tidak boleh diimport dari internal/core sama sekali.
type Config struct {
	HTTPPort string
	StoreKind string // "memory" | "file"
	FilePath  string
}

// Load membaca env var dan mengisi default kalau kosong.
// Nilai hasilnya berupa string/angka polos, jadi saat dioper ke core,
// core tidak pernah tahu bahwa nilai itu asalnya dari env var.
func Load() Config {
	return Config{
		HTTPPort:  getEnv("HTTP_PORT", "8080"),
		StoreKind: getEnv("STORE_KIND", "memory"),
		FilePath:  getEnv("STORE_FILE", "products.json"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
