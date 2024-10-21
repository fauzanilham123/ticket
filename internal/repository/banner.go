package repository

import (
	"api-ticket/internal/entity"
	"errors"
	"gorm.io/gorm"
)

type BannerRepository struct {
	db *gorm.DB
}

func NewBannerRepository(db *gorm.DB) entity.IBannerRepository {
	return &BannerRepository{db: db}
}

func (repo *BannerRepository) Create(banner entity.Banner) error {
	return repo.db.Create(&banner).Error
}

func (repo *BannerRepository) FindAll(filter entity.RequestGetBanner) ([]entity.Banner, int64, error) {
	var banners []entity.Banner
	query := repo.db.Model(&entity.Banner{})

	if filter.Title != nil {
		query = query.Where("title ILIKE ?", "%"+*filter.Title+"%")
	}

	//Count
	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	//Apply Limit
	if filter.Limit != nil && filter.Offset != nil {
		limit := *filter.Limit
		page := (*filter.Offset / limit) + 1
		offset := (page - 1) * limit

		query = query.Limit(int(limit)).Offset(int(offset))
	} else if filter.Limit != nil {
		query = query.Limit(int(*filter.Limit))
	}

	result := query.Find(&banners)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, 0, result.Error
		}
	}

	return banners, totalCount, nil
}
