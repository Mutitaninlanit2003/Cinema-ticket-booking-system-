package repository

import (
	"errors"

	"cinema-ticket-booking/config"
	"cinema-ticket-booking/models"

	"gorm.io/gorm"
)

type TheaterRepository struct {
	DB *gorm.DB
}

func NewTheaterRepository() *TheaterRepository {
	return &TheaterRepository{
		DB: config.DB,
	}
}

func (r *TheaterRepository) GetTheaters() ([]*models.Theater, error) {
	var theaters []*models.Theater

	if err := r.DB.Find(&theaters).Error; err != nil {
		return nil, err
	}

	return theaters, nil
}

func (r *TheaterRepository) CreateTheater(theater *models.Theater) error {
	return r.DB.Create(theater).Error
}

func (r *TheaterRepository) UpdateTheater(
	theaterID string,
	theater *models.Theater,
) error {

	result := r.DB.
		Model(&models.Theater{}).
		Where("theater_id = ?", theaterID).
		Updates(map[string]interface{}{
			"name":       theater.Name,
			"address":    theater.Address,
			"status":     theater.Status,
			"location":   theater.Location,
			"capacity":   theater.Capacity,
			"updated_at": theater.UpdatedAt,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("theater not found")
	}

	return nil
}

func (r *TheaterRepository) DeleteTheater(
	theaterID string,
) error {

	result := r.DB.
		Where("theater_id = ?", theaterID).
		Delete(&models.Theater{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("theater not found")
	}

	return nil
}
