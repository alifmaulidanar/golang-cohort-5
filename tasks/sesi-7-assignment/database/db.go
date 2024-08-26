package database

import (
	"fmt"
	"log"
	"os"
	"sesi-7-assignment/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	err error
)

// konek gorm database mysql
func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbport := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	config := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, dbport, dbname)
	db, err = gorm.Open(mysql.Open(config), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	db.Debug().AutoMigrate(&models.Order{}, &models.Item{})
}

func GetDB() *gorm.DB {
	return db
}
