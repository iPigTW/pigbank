package initialize

import "github.com/ipigtw/api/models"

func SyncDB() {
	DB.AutoMigrate(&models.User{})
}