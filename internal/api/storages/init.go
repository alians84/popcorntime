package storages

import (
	"gorm.io/gorm"
	"popcorntime-project/internal/config"
	"popcorntime-project/internal/repository/storages"
)

func Init(db *gorm.DB, b2client *config.B2Client) *Handler {
	storageService := storages.NewService(db, b2client)

	return NewHandler(storageService)
}
