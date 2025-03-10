package database

import (
	"database/sql"
	"library_project/internal/model"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load("../../.env")

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	connectionString := os.Getenv("DATABASE_URL")
	DB, err = gorm.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	createDatabaseIfNotExists()

	err = DB.AutoMigrate(&model.Song{}).Error
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}

func createDatabaseIfNotExists() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}
	defer db.Close()

	dbName := os.Getenv("DB_NAME")
	_, err = db.Exec("CREATE DATABASE " + dbName)
	if err != nil {
		if err.Error() != "pq: database \""+dbName+"\" already exists" {
			log.Fatal("Failed to create database:", err)
		}
	}
}
