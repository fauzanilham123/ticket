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

func InitTalentService(db *gorm.DB) entity.ITalentService {
	talentRepo := repository.NewTalentRepository(db)
	talentService := service.NewTalentService(talentRepo)

	return talentService
}

func InitEventService(db *gorm.DB) entity.IEventService {
	eventRepo := repository.NewEventRepository(db)
	eventService := service.NewEventService(eventRepo)

	return eventService
}

func InitAuthService(db *gorm.DB) entity.IAuthService {
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo)

	return authService
}
