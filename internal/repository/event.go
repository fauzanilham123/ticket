package repository

import (
	"api-ticket/internal/entity"
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) entity.IEventRepository {
	return &EventRepository{db: db}
}

func (repo *EventRepository) Create(event entity.Event) error {
	return repo.db.Omit("updated_at").Create(&event).Error
}

func (repo *EventRepository) FindAll(filter entity.RequestGetEvent) ([]entity.Event, int64, error) {
	var events []entity.Event
	query := repo.db.Model(&entity.Event{})

	if filter.Title != nil {
		query = query.Where("title ILIKE ?", "%"+*filter.Title+"%")
	}

	if filter.Id_talent != nil {
		query = query.Where("id_title ILIKE ?", "%"+*filter.Id_talent+"%")
	}

	if filter.Desc != nil {
		query = query.Where("desc ILIKE ?", "%"+*filter.Desc+"%")
	}

	if filter.Location != nil {
		query = query.Where("location ILIKE ?", "%"+*filter.Location+"%")
	}

	if filter.Sk != nil {
		query = query.Where("sk ILIKE ?", "%"+*filter.Sk+"%")
	}

	if filter.Tag != nil {
		query = query.Where("tag ILIKE ?", "%"+*filter.Tag+"%")
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

	result := query.Find(&events)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, 0, result.Error
		}
	}

	return events, totalCount, nil
}

func (repo *EventRepository) FindByID(id string) (entity.Event, error) {
	// Cari event berdasarkan ID
	var events entity.Event
	err := repo.db.Where("id = ?", id).First(&events).Error
	return events, err
}

func (repo *EventRepository) Delete(id string) (entity.Event, error) {
	// Cari event berdasarkan ID
	var event entity.Event
	if err := repo.db.First(&event, "id = ?", id).Error; err != nil {
		return event, err
	}
	// Hapus event yang ditemukan berdasarkan ID
	if err := repo.db.Delete(&event).Error; err != nil {
		return event, err
	}
	return event, nil
}

func (repo *EventRepository) Update(id string, updatedEvent entity.Event) (entity.Event, error) {
	// Cari event berdasarkan ID terlebih dahulu
	var event entity.Event
	if err := repo.db.First(&event, "id = ?", id).Error; err != nil {
		return event, err
	}

	// Lakukan update dengan nilai-nilai dari updatedBanner
	if err := repo.db.Model(&event).Updates(updatedEvent).Error; err != nil {
		return event, err
	}

	return event, nil
}
