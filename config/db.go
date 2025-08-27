package config

import (
	"fmt"
	"log"
	"os"
	"ticket/domain/entity"

	"gorm.io/driver/mysql"

	"github.com/joho/godotenv"

	"gorm.io/gorm"
)

func Db() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error load env")
	}

	user := os.Getenv("DB_USERNAME")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbName)

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&entity.User{}, &entity.Event{}, &entity.Ticket{})
	if err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}

	return db
}
