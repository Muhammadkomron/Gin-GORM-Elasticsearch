package initializers

import "gin-gorm-tutorial/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{}, &models.Item{})
}
