package initialize

import "github.com/ipigtw/pigbank/models"

func SyncDB() {
	DB.AutoMigrate(&models.User{})
}