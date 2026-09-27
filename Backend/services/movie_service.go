package services

import (
	"time"

	"cinema-ticket-booking/models"
	"cinema-ticket-booking/repository"
)

type MovieService struct {
	Repo *repository.MovieRepository
}

func NewMovieService() *MovieService {
	return &MovieService{
		Repo: repository.NewMovieRepository(),
	}
}

func (s *MovieService) CreateMovie(movie models.Movie) error {

	movie.CreatedAt = time.Now()
	movie.UpdatedAt = time.Now()

	return s.Repo.CreateMovie(movie)
}

func (s *MovieService) GetMovies() ([]models.Movie, error) {

	return s.Repo.GetMovies()
}
func (s *MovieService) UpdateMovie(
	movieID string,
	movie models.Movie,
) error {

	movie.UpdatedAt = time.Now()

	return s.Repo.UpdateMovie(
		movieID,
		movie,
	)
}
func (s *MovieService) DeleteMovie(
	movieID string,
) error {

	return s.Repo.DeleteMovie(movieID)
}
func (s *MovieService) GetMovieByID(movieID string) (*models.Movie, error) {
	return s.Repo.GetMovieByID(movieID)
}
