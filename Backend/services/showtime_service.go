package services

import (
	"context"
	"time"

	"cinema-ticket-booking/models"
	"cinema-ticket-booking/repository"
)

type ShowtimeService struct {
	repo *repository.ShowtimeRepository
}

func NewShowtimeService(
	repo *repository.ShowtimeRepository,
) *ShowtimeService {

	return &ShowtimeService{
		repo: repo,
	}

}

func (s *ShowtimeService) CreateShowtime(
	ctx context.Context,
	showtime models.Showtime,
) error {

	showtime.CreatedAt = time.Now()
	showtime.UpdatedAt = time.Now()

	return s.repo.Create(
		ctx,
		showtime,
	)

}

func (s *ShowtimeService) GetShowtimes(
	ctx context.Context,
) ([]models.Showtime, error) {

	return s.repo.FindAll(ctx)

}
func (s *ShowtimeService) UpdateShowtime(
	ctx context.Context,
	showtimeID string,
	showtime models.Showtime,
) error {

	showtime.UpdatedAt = time.Now()

	return s.repo.Update(
		ctx,
		showtimeID,
		showtime,
	)
}

func (s *ShowtimeService) DeleteShowtime(
	ctx context.Context,
	showtimeID string,
) error {

	return s.repo.Delete(
		ctx,
		showtimeID,
	)
}
