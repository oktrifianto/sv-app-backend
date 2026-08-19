package main

import (
	"log"

	"sv-app/backend/internal/config"
	"sv-app/backend/internal/database"
	"sv-app/backend/internal/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	app := fiber.New()
	routes.Setup(app, db)

	log.Printf("Server starting on :%s", cfg.ServerPort)
	if err := app.Listen(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
