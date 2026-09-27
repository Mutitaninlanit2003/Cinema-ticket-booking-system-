package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"cinema-ticket-booking/models"
	"cinema-ticket-booking/repository"
	"cinema-ticket-booking/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var userService *services.UserService

// GetUsers godoc
// @Summary Get all users
// @Description Retrieve all users
// @Tags Users
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func GetUsers(c *gin.Context) {

	if userService == nil {
		userService = services.NewUserService()
	}

	users, err := userService.GetUsers()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, users)
}

// CreateUser godoc
// @Summary Create User
// @Description Create a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "User"
// @Success 201 {object} map[string]string
// @Router /users [post]
func CreateUser(c *gin.Context) {

	if userService == nil {
		userService = services.NewUserService()
	}

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := userService.CreateUser(user); err != nil {

		if strings.Contains(err.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, gin.H{
				"error": "email already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
	})
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param user body models.User true "User"
// @Success 200 {object} map[string]string
// @Router /users/{user_id} [put]
func UpdateUser(c *gin.Context) {

	if userService == nil {
		userService = services.NewUserService()
	}

	userID := c.Param("user_id")

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {

		if validationErrors, ok := err.(validator.ValidationErrors); ok {

			errors := make([]string, 0)

			for _, fieldError := range validationErrors {

				switch fieldError.Tag() {

				case "required":
					errors = append(errors, fieldError.Field()+" is required")

				case "email":
					errors = append(errors, fieldError.Field()+" must be a valid email")

				case "oneof":
					errors = append(errors, fieldError.Field()+" must be either user or admin")
				}
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"errors": errors,
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := userService.UpdateUser(userID, user); err != nil {

		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
	})
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete user by ID
// @Tags Users
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} map[string]string
// @Router /users/{user_id} [delete]
func DeleteUser(c *gin.Context) {

	if userService == nil {
		userService = services.NewUserService()
	}

	userID := c.Param("user_id")

	if err := userService.DeleteUser(userID); err != nil {

		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

// GetMe godoc
// @Summary Get current user
// @Description Get the currently authenticated user's profile
// @Tags Users
// @Produce json
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/me [get]
func GetMe(c *gin.Context) {
	if userService == nil {
		userService = services.NewUserService()
	}

	userIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "authentication required",
		})
		return
	}

	userID, ok := userIDValue.(string)

	if !ok || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user_id",
		})
		return
	}

	user, err := userService.GetUserByID(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UploadProfileImage godoc
// @Summary Upload profile image
// @Description Upload profile image for the authenticated user
// @Tags Users
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Profile image"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/me/profile-image [post]
func UploadProfileImage(c *gin.Context) {
	if userService == nil {
		userService = services.NewUserService()
	}

	// Get user_id from JWT
	userIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "authentication required",
		})
		return
	}

	userID, ok := userIDValue.(string)

	if !ok || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user_id",
		})
		return
	}

	// Get uploaded file
	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "profile image file is required",
		})
		return
	}

	// Validate file type
	contentType := file.Header.Get("Content-Type")

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}

	if !allowedTypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "only JPG, PNG and WEBP images are allowed",
		})
		return
	}

	// Maximum 5 MB
	const maxFileSize = 5 * 1024 * 1024

	if file.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "profile image must not exceed 5 MB",
		})
		return
	}

	// Create upload directory
	uploadDir := "uploads/profiles"

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create upload directory",
		})
		return
	}

	// Get extension
	ext := strings.ToLower(filepath.Ext(file.Filename))

	// One image per user
	filename := userID + ext

	filePath := filepath.Join(uploadDir, filename)

	// Save file
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save profile image",
		})
		return
	}

	// Save path to database
	imageURL := "/uploads/profiles/" + filename

	if err := userService.UpdateProfileImage(
		userID,
		imageURL,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "profile image uploaded successfully",
		"profile_image": imageURL,
	})
}

// --------------------------------------------------
// Update Me (Self-service)
// --------------------------------------------------

type UpdateMeRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// UpdateMe godoc
// @Summary Update current user's profile
// @Description Update the authenticated user's own name/email
// @Tags Users
// @Accept json
// @Produce json
// @Param request body UpdateMeRequest true "Profile fields"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /users/me [put]
func UpdateMe(c *gin.Context) {

	if userService == nil {
		userService = services.NewUserService()
	}

	userIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "authentication required",
		})
		return
	}

	userID, ok := userIDValue.(string)

	if !ok || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user_id",
		})
		return
	}

	var req UpdateMeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// ดึงข้อมูลปัจจุบันก่อน เพื่อไม่ให้ Role/ProfileImage หายไป
	// ตอน Update (เพราะ UpdateUser จะเขียนทับทุก field ที่ส่งมา)

	currentUser, err := userService.GetUserByID(userID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	currentUser.Name = req.Name
	currentUser.Email = req.Email

	if err := userService.UpdateUser(userID, *currentUser); err != nil {

		if strings.Contains(err.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, gin.H{
				"error": "email already in use",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	updatedUser, err := userService.GetUserByID(userID)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "profile updated successfully",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "profile updated successfully",
		"user":    updatedUser,
	})
}

// --------------------------------------------------
// Change My Password (Self-service)
// --------------------------------------------------

type ChangeMyPasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// ChangeMyPassword godoc
// @Summary Change current user's password
// @Description Change the authenticated user's own password
// @Tags Users
// @Accept json
// @Produce json
// @Param request body ChangeMyPasswordRequest true "Password fields"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /users/me/password [put]
func ChangeMyPassword(c *gin.Context) {

	userIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "authentication required",
		})
		return
	}

	userID, ok := userIDValue.(string)

	if !ok || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user_id",
		})
		return
	}

	var req ChangeMyPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userRepo := repository.NewUserRepository()

	user, err := userRepo.FindByUserID(userID)

	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	// Google-only accounts อาจไม่มี password เลย
	if user.PasswordHash == nil || *user.PasswordHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "this account does not have a password set",
		})
		return
	}

	if !services.CheckPassword(
		req.CurrentPassword,
		*user.PasswordHash,
	) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "current password is incorrect",
		})
		return
	}

	newHash, err := services.HashPassword(req.NewPassword)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	if err := userRepo.UpdatePassword(userID, newHash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "password changed successfully",
	})
}
