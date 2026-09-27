package services

import (
	"time"

	"cinema-ticket-booking/models"
	"cinema-ticket-booking/repository"
)

type ScreenService struct {
	Repo *repository.ScreenRepository
}

func NewScreenService() *ScreenService {

	return &ScreenService{
		Repo: repository.NewScreenRepository(),
	}
}

func (s *ScreenService) CreateScreen(screen models.Screen) error {

	theater, err := s.Repo.GetTheaterByUUID(screen.TheaterID)

	if err != nil {
		return err
	}

	screen.TheaterRefID = theater.ID

	screen.CreatedAt = time.Now()
	screen.UpdatedAt = time.Now()

	return s.Repo.CreateScreen(screen)
}

func (s *ScreenService) GetScreens() ([]models.Screen, error) {

	return s.Repo.GetScreens()
}

func (s *ScreenService) UpdateScreen(
	screenID string,
	screen models.Screen,
) error {

	screen.UpdatedAt = time.Now()

	return s.Repo.UpdateScreen(screenID, screen)
}

func (s *ScreenService) DeleteScreen(
	screenID string,
) error {

	return s.Repo.DeleteScreen(screenID)
}

// ForceDeleteScreen ลบ screen พร้อม showtimes, bookings, payments, seats
// ใช้กับข้อมูลทดสอบ/dev เท่านั้น
func (s *ScreenService) ForceDeleteScreen(
	screenID string,
) error {

	return s.Repo.ForceDeleteScreen(screenID)
}
