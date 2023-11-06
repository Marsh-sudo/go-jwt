package initializers

import (
	"github.com/marsh-sudo/go-jwt/models"
)

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}