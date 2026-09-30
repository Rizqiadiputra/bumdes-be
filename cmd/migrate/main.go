package main

import (
	"log"

	"github.com/liyansasongko/bumdes-be/internal/config"
	"github.com/liyansasongko/bumdes-be/internal/database"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("gagal terhubung ke database: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("gagal migrasi database: %v", err)
	}

	log.Println("migrasi database berhasil")
}
