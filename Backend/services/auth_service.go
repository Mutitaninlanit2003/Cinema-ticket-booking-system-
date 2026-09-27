package services

import (
	"errors"
	"os"
	"time"

	"cinema-ticket-booking/models"
	"cinema-ticket-booking/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	jwtSecret         []byte
	passwordResetRepo *repository.PasswordResetRepository
}

func NewAuthService() (*AuthService, error) {
	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		return nil, errors.New("JWT_SECRET is not configured")
	}

	return &AuthService{
		jwtSecret:         []byte(secret),
		passwordResetRepo: repository.NewPasswordResetRepository(),
	}, nil
}

// --------------------------------------------------
// Generate JWT Token
// --------------------------------------------------

func (s *AuthService) GenerateToken(
	userID string,
	role string,
) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(s.jwtSecret)
}

// --------------------------------------------------
// Password Hash
// --------------------------------------------------

func HashPassword(password string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// --------------------------------------------------
// Check Password
// --------------------------------------------------

func CheckPassword(
	password string,
	passwordHash string,
) bool {

	err := bcrypt.CompareHashAndPassword(
		[]byte(passwordHash),
		[]byte(password),
	)

	return err == nil
}

// --------------------------------------------------
// Create Password Reset Token
// --------------------------------------------------

// CreatePasswordResetToken creates a new password reset token
// for the specified user.
//
// The token will expire after 30 minutes.
func (s *AuthService) CreatePasswordResetToken(
	userID string,
) (*models.PasswordResetToken, error) {

	// Remove any previous reset tokens belonging to this user.
	if err := s.passwordResetRepo.DeleteByUserID(userID); err != nil {
		return nil, err
	}

	resetToken := &models.PasswordResetToken{
		Token:     uuid.NewString(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

	if err := s.passwordResetRepo.Create(resetToken); err != nil {
		return nil, err
	}

	return resetToken, nil
}

// --------------------------------------------------
// Find Password Reset Token
// --------------------------------------------------

// FindPasswordResetToken verifies that the token exists
// and has not expired.
func (s *AuthService) FindPasswordResetToken(
	token string,
) (*models.PasswordResetToken, error) {

	return s.passwordResetRepo.FindByToken(token)
}

// --------------------------------------------------
// Delete Password Reset Token
// --------------------------------------------------

// DeletePasswordResetToken removes a reset token after
// it has been successfully used.
func (s *AuthService) DeletePasswordResetToken(
	token string,
) error {

	return s.passwordResetRepo.DeleteByToken(token)
}
