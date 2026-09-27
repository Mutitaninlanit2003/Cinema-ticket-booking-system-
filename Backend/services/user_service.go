package services

import (
	"time"

	"cinema-ticket-booking/models"
	"cinema-ticket-booking/repository"
)

type UserService struct {
	Repo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		Repo: repository.NewUserRepository(),
	}
}

// เพิ่มผู้ใช้
func (s *UserService) CreateUser(user models.User) error {

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return s.Repo.CreateUser(user)
}

// ดึงผู้ใช้ทั้งหมด
func (s *UserService) GetUsers() ([]models.User, error) {
	return s.Repo.GetUsers()
}
func (s *UserService) GetUserByID(
	userID string,
) (*models.User, error) {

	return s.Repo.FindByUserID(userID)
}

func (s *UserService) UpdateUser(
	userID string,
	user models.User,
) error {

	user.UpdatedAt = time.Now()

	return s.Repo.UpdateUser(
		userID,
		&user,
	)
}

func (s *UserService) DeleteUser(
	userID string,
) error {

	return s.Repo.DeleteUser(userID)
}
func (s *UserService) UpdateProfileImage(
	userID string,
	profileImage string,
) error {

	return s.Repo.UpdateProfileImage(
		userID,
		profileImage,
	)
}
