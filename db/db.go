package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var GlobalConnection *gorm.DB

func Connect() {
	var (
		DB_HOST     = os.Getenv("DB_HOST")
		DB_USER     = os.Getenv("DB_USER")
		DB_PASSWORD = os.Getenv("DB_PASSWORD")
		DB_NAME     = os.Getenv("DB_NAME")
		DB_PORT     = os.Getenv("DB_PORT")
	)

	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		DB_HOST,
		DB_USER,
		DB_PASSWORD,
		DB_NAME,
		DB_PORT,
	)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	GlobalConnection = db // make gorm PostgreSQL connection global
	if err != nil {
		panic(err)
	}

	log.Println("SUCCESS: Connected to PostgreSQL database.")

	migrationErr := GlobalConnection.AutoMigrate(&Todo{})
	if migrationErr != nil {
		log.Fatalf("ERROR: Failed to perform database migration: %s\n", err)
	}

	log.Println("SUCCESS: PostgreSQL migration completed (Some tables won't be created if they already exist).")
}

type Todo struct {
	ID        uuid.UUID      `json:"id" uri:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}
