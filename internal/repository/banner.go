package repository

import (
	"api-ticket/internal/entity"
	"errors"
	"fmt"

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

	// Apply OrderBy and Sort
	if filter.OrderBy != nil {
		sortOrder := "ASC"
		if filter.Sort != nil && *filter.Sort == "desc" {
			sortOrder = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", *filter.OrderBy, sortOrder))
	}

	result := query.Find(&banners)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, 0, result.Error
		}
	}

	return banners, totalCount, nil
}

func (repo *BannerRepository) FindByID(id string) (entity.Banner, error) {
	// Cari banner berdasarkan ID
	var banners entity.Banner
	err := repo.db.Where("id = ?", id).First(&banners).Error
	return banners, err
}

func (repo *BannerRepository) Delete(id string) (entity.Banner, error) {
	// Cari banner berdasarkan ID
	var banner entity.Banner
	if err := repo.db.First(&banner, "id = ?", id).Error; err != nil {
		return banner, err
	}
	// Hapus banner yang ditemukan berdasarkan ID
	if err := repo.db.Delete(&banner).Error; err != nil {
		return banner, err
	}
	return banner, nil
}

func (repo *BannerRepository) Update(id string, updatedBanner entity.Banner) (entity.Banner, error) {
	// Cari banner berdasarkan ID terlebih dahulu
	var banner entity.Banner
	if err := repo.db.First(&banner, "id = ?", id).Error; err != nil {
		return banner, err
	}

	// Lakukan update dengan nilai-nilai dari updatedBanner
	if err := repo.db.Model(&banner).Updates(updatedBanner).Error; err != nil {
		return banner, err
	}

	return banner, nil
}
