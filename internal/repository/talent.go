package repository

import (
	"api-ticket/internal/entity"
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type TalentRepository struct {
	db *gorm.DB
}

func NewTalentRepository(db *gorm.DB) entity.ITalentRepository {
	return &TalentRepository{db: db}
}

func (repo *TalentRepository) Create(talent entity.Talent) error {
	return repo.db.Create(&talent).Error
}

func (repo *TalentRepository) FindAll(filter entity.RequestGetTalent) ([]entity.Talent, int64, error) {
	var talents []entity.Talent
	query := repo.db.Model(&entity.Talent{})

	if filter.Name != nil {
		query = query.Where("name ILIKE ?", "%"+*filter.Name+"%")
	}

	if filter.Id_promotor_created != nil {
		query = query.Where("id_promotor_created::text ILIKE ?", "%"+strconv.Itoa(*filter.Id_promotor_created)+"%")
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

	result := query.Find(&talents)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, 0, result.Error
		}
	}

	return talents, totalCount, nil
}

func (repo *TalentRepository) FindByID(id string) (entity.Talent, error) {
	// Cari talent berdasarkan ID
	var talents entity.Talent
	err := repo.db.Where("id = ?", id).First(&talents).Error
	return talents, err
}

func (repo *TalentRepository) Delete(id string) (entity.Talent, error) {
	// Cari talent berdasarkan ID
	var talent entity.Talent
	if err := repo.db.First(&talent, "id = ?", id).Error; err != nil {
		return talent, err
	}
	// Hapus talent yang ditemukan berdasarkan ID
	if err := repo.db.Delete(&talent).Error; err != nil {
		return talent, err
	}
	return talent, nil
}

func (repo *TalentRepository) Update(id string, updatedTalent entity.Talent) (entity.Talent, error) {
	// Cari talent berdasarkan ID terlebih dahulu
	var talent entity.Talent
	if err := repo.db.First(&talent, "id = ?", id).Error; err != nil {
		return talent, err
	}

	// Lakukan update dengan nilai-nilai dari updatedBanner
	if err := repo.db.Model(&talent).Updates(updatedTalent).Error; err != nil {
		return talent, err
	}

	return talent, nil
}
