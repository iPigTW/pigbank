package initialize

import (
	"os"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
func ConnectToDb() {
	dsn := os.Getenv("DB")
	DB, _ = gorm.Open(postgres.Open(dsn), &gorm.Config{})

}