package cmd

import (
	"api-ticket/internal/entity"
	"api-ticket/internal/repository"
	"api-ticket/internal/service"
	"gorm.io/gorm"
)

func InitBannerService(db *gorm.DB) entity.IBannerService {
	bannerRepo := repository.NewBannerRepository(db)
	bannerService := service.NewBannerService(bannerRepo)

	return bannerService
}
