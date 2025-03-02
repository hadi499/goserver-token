
package database

import (
	"server-crud/models"
	"fmt"
	"log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Ganti dengan kredensial database yang sesuai
	dsn := "hadi:admin123@tcp(127.0.0.1:3306)/go_crud?charset=utf8mb4&parseTime=True&loc=Local"

	// Membuka koneksi ke database
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrasi otomatis (pastikan model sudah benar)
	err = database.AutoMigrate(&models.Product{}, &models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Set database global
	DB = database

	fmt.Println("✅ Database connected successfully!")
}

// GetDB mengembalikan instance database
func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("❌ Database not initialized. Call ConnectDatabase() first.")
	}
	return DB
}
