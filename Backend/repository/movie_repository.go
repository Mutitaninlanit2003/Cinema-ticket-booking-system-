package repository

import (
	"cinema-ticket-booking/config"
	"cinema-ticket-booking/models"
	"errors"

	"gorm.io/gorm"

	"fmt"
)

type MovieRepository struct {
	DB *gorm.DB
}

func NewMovieRepository() *MovieRepository {
	return &MovieRepository{
		DB: config.DB,
	}
}

// Create
func (r *MovieRepository) CreateMovie(movie models.Movie) error {
	return r.DB.Create(&movie).Error
}

// Read
func (r *MovieRepository) GetMovies() ([]models.Movie, error) {

	var movies []models.Movie

	if err := r.DB.Find(&movies).Error; err != nil {
		return nil, err
	}

	return movies, nil
}

// Update
func (r *MovieRepository) UpdateMovie(
	movieID string,
	movie models.Movie,
) error {

	result := r.DB.
		Model(&models.Movie{}).
		Where("movie_id = ?", movieID).
		Updates(map[string]interface{}{
			"title":        movie.Title,
			"description":  movie.Description,
			"genre":        movie.Genre,
			"duration":     movie.Duration,
			"language":     movie.Language,
			"subtitle":     movie.Subtitle,
			"rating":       movie.Rating,
			"poster_url":   movie.PosterURL,
			"trailer_url":  movie.TrailerURL,
			"release_date": movie.ReleaseDate,
			"status":       movie.Status,
			"updated_at":   movie.UpdatedAt,
		})

	fmt.Println("MovieID =", movieID)
	fmt.Println("RowsAffected =", result.RowsAffected)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("movie not found")
	}

	return nil
}

// Delete
func (r *MovieRepository) DeleteMovie(
	movieID string,
) error {

	result := r.DB.
		Where("movie_id = ?", movieID).
		Delete(&models.Movie{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("movie not found")
	}

	return nil
}

func (r *MovieRepository) GetMovieByID(movieID string) (*models.Movie, error) {
	var movie models.Movie

	result := r.DB.
		Where("movie_id = ?", movieID).
		First(&movie)

	if result.Error != nil {
		return nil, result.Error
	}

	return &movie, nil
}
