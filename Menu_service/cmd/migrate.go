package main

import (
	"log"
	"menu/config"
	"menu/internal/migration"
)

func main() {
	cfg := config.LoadConfig()
	db := config.ConnectToMongo(cfg.MongoURI, cfg.DatabaseName)

	log.Println("🚀 Running migrations...")
	migrations.Run(db)
}
