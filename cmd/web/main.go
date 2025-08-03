package main

import (
	"log"

	// Import the pgx driver
	"github.com/Archnick/go-ecommerce/Internal/api"
	"github.com/Archnick/go-ecommerce/Internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Define the PostgreSQL connection string.
	dsn := "host=localhost user=myuser password=mypassword dbname=etsydb port=5432 sslmode=disable"

	// 2. Open a GORM database connection.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	db.AutoMigrate(
		&models.User{},
		&models.Shop{},
		&models.Product{},
		&models.ProductImage{},
		&models.Category{},
		&models.Order{},
		&models.OrderItem{},
		&models.Review{},
	)

	// 3. Create and start the server.
	server := api.NewServer(db)
	if err := server.Start(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
