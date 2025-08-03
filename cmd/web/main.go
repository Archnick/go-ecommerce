package main

import (
	"log"
	"log/slog"

	// Import the pgx driver
	"github.com/Archnick/go-ecommerce/Internal/api"
	"github.com/Archnick/go-ecommerce/Internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
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

	seedAdmin(db)

	// 3. Create and start the server.
	server := api.NewServer(db)
	if err := server.Start(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func seedAdmin(db *gorm.DB) {
	var adminUser models.User
	// Check if an admin user already exists
	err := db.Where("role = ?", models.AdminRole).First(&adminUser).Error

	if err != nil {
		// If no admin exists (record not found), create one
		if err == gorm.ErrRecordNotFound {
			slog.Info("No admin user found, creating one...")
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("adminpassword"), bcrypt.DefaultCost)

			admin := models.User{
				Email:    "admin@example.com",
				Password: string(hashedPassword),
				Role:     string(models.AdminRole),
			}
			db.Create(&admin)
			slog.Info("Admin user created successfully")
		}
	} else {
		slog.Info("Admin user already exists")
	}
}
