package repository

import (
	"errors"
	"time"

	"cinema-ticket-booking/config"
	"cinema-ticket-booking/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		DB: config.DB,
	}
}

// --------------------------------------------------
// Create User
// --------------------------------------------------

func (r *UserRepository) CreateUser(user models.User) error {
	return r.DB.Create(&user).Error
}

// --------------------------------------------------
// Get Users
// --------------------------------------------------

func (r *UserRepository) GetUsers() ([]models.User, error) {

	var users []models.User

	err := r.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

// --------------------------------------------------
// Update User
// --------------------------------------------------

func (r *UserRepository) UpdateUser(
	userID string,
	user *models.User,
) error {

	result := r.DB.
		Model(&models.User{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"name":          user.Name,
			"email":         user.Email,
			"role":          user.Role,
			"profile_image": user.ProfileImage,
			"updated_at":    user.UpdatedAt,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// --------------------------------------------------
// Delete User
// --------------------------------------------------

func (r *UserRepository) DeleteUser(
	userID string,
) error {

	result := r.DB.
		Where("user_id = ?", userID).
		Delete(&models.User{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// --------------------------------------------------
// Find User By Google ID
// --------------------------------------------------

func (r *UserRepository) FindByGoogleID(
	googleID string,
) (*models.User, error) {

	var user models.User

	err := r.DB.
		Where("google_id = ?", googleID).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// --------------------------------------------------
// Find User By Email
// --------------------------------------------------

// FindByEmail looks up a user by email, case-insensitively.
func (r *UserRepository) FindByEmail(
	email string,
) (*models.User, error) {

	var user models.User

	err := r.DB.
		Where("LOWER(email) = LOWER(?)", email).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// --------------------------------------------------
// Update Password
// --------------------------------------------------

func (r *UserRepository) UpdatePassword(
	userID string,
	passwordHash string,
) error {

	result := r.DB.
		Model(&models.User{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"password_hash": passwordHash,
			"updated_at":    time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// --------------------------------------------------
// Find User By User ID
// --------------------------------------------------

func (r *UserRepository) FindByUserID(
	userID string,
) (*models.User, error) {

	var user models.User

	err := r.DB.
		Where("user_id = ?", userID).
		First(&user).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}

		return nil, err
	}

	return &user, nil
}

// --------------------------------------------------
// Update Profile Image
// --------------------------------------------------

func (r *UserRepository) UpdateProfileImage(
	userID string,
	profileImage string,
) error {

	result := r.DB.
		Model(&models.User{}).
		Where("user_id = ?", userID).
		Update("profile_image", profileImage)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// ==================================================
// FORGOT PASSWORD
// ==================================================

// --------------------------------------------------
// Save Reset Token
// --------------------------------------------------

func (r *UserRepository) SaveResetToken(
	userID string,
	resetToken string,
	resetTokenExpiry time.Time,
) error {

	result := r.DB.
		Model(&models.User{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"reset_token":        resetToken,
			"reset_token_expiry": resetTokenExpiry,
			"updated_at":         time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// --------------------------------------------------
// Find User By Reset Token
// --------------------------------------------------

func (r *UserRepository) FindByResetToken(
	resetToken string,
) (*models.User, error) {

	var user models.User

	err := r.DB.
		Where(
			"reset_token = ? AND reset_token_expiry > ?",
			resetToken,
			time.Now(),
		).
		First(&user).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid or expired reset token")
		}

		return nil, err
	}

	return &user, nil
}

// --------------------------------------------------
// Clear Reset Token
// --------------------------------------------------

func (r *UserRepository) ClearResetToken(
	userID string,
) error {

	result := r.DB.
		Model(&models.User{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"reset_token":        nil,
			"reset_token_expiry": nil,
			"updated_at":         time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// --------------------------------------------------
// Link Google ID to an existing user
// --------------------------------------------------
// ใช้ตอนพบว่า email ตรงกับ user ที่เคยสมัครด้วย password ไว้แล้ว
// แต่ยังไม่เคย login ผ่าน Google — ผูก GoogleID เข้ากับบัญชีเดิม
// แทนที่จะสร้าง user ซ้ำ

func (r *UserRepository) LinkGoogleID(
	userID string,
	googleID string,
) error {

	result := r.DB.
		Model(&models.User{}).
		Where("user_id = ?", userID).
		Update("google_id", googleID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
