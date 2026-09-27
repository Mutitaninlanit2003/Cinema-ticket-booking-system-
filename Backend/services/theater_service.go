package services

import (
	"time"

	"cinema-ticket-booking/models"
	"cinema-ticket-booking/repository"
)

type TheaterService struct {
	Repo *repository.TheaterRepository
}

func NewTheaterService() *TheaterService {
	return &TheaterService{
		Repo: repository.NewTheaterRepository(),
	}
}

func (s *TheaterService) GetTheaters() ([]*models.Theater, error) {
	return s.Repo.GetTheaters()
}

func (s *TheaterService) CreateTheater(theater *models.Theater) error {
	theater.CreatedAt = time.Now()
	theater.UpdatedAt = time.Now()

	return s.Repo.CreateTheater(theater)
}

func (s *TheaterService) UpdateTheater(
	theaterID string,
	theater *models.Theater,
) error {
	theater.UpdatedAt = time.Now()

	return s.Repo.UpdateTheater(
		theaterID,
		theater,
	)
}

func (s *TheaterService) DeleteTheater(
	theaterID string,
) error {
	return s.Repo.DeleteTheater(theaterID)
}
