package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"cinema-ticket-booking/config"
	"cinema-ticket-booking/models"
	"cinema-ticket-booking/repository"
	"cinema-ticket-booking/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

// --------------------------------------------------
// Request Models
// --------------------------------------------------

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// --------------------------------------------------
// Auth Controller
// --------------------------------------------------

type AuthController struct {
	userRepo    *repository.UserRepository
	authService *services.AuthService
}

func NewAuthController(
	userRepo *repository.UserRepository,
	authService *services.AuthService,
) *AuthController {
	return &AuthController{
		userRepo:    userRepo,
		authService: authService,
	}
}

// --------------------------------------------------
// Register
// --------------------------------------------------

// Register godoc
// @Summary Register a new user
// @Description Create a new account with email/password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} map[string]interface{}
// @Router /auth/register [post]

func (a *AuthController) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Normalize email
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Check duplicate email
	if existing, err := a.userRepo.FindByEmail(req.Email); err == nil && existing != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "an account with this email already exists",
		})
		return
	}

	// Hash password
	passwordHash, err := services.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	// Create user
	user := models.User{
		UserID:       uuid.NewString(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: &passwordHash,
		Role:         "user",
	}

	if err := a.userRepo.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create user",
		})
		return
	}

	// Generate JWT
	appToken, err := a.authService.GenerateToken(
		user.UserID,
		user.Role,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "registration successful",
		"token":   appToken,
		"user": gin.H{
			"user_id": user.UserID,
			"name":    user.Name,
			"email":   user.Email,
			"role":    user.Role,
		},
	})
}

// --------------------------------------------------
// Login
// --------------------------------------------------

// Login godoc
// @Summary Log in with email/password
// @Description Authenticate and receive a JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Router /auth/login [post]

func (a *AuthController) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Find user
	user, err := a.userRepo.FindByEmail(req.Email)

	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	// Google-only users may not have a password
	if user.PasswordHash == nil || *user.PasswordHash == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "this account does not have a password login",
		})
		return
	}

	// Check password
	if !services.CheckPassword(
		req.Password,
		*user.PasswordHash,
	) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	// Generate JWT
	appToken, err := a.authService.GenerateToken(
		user.UserID,
		user.Role,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   appToken,
		"user": gin.H{
			"user_id": user.UserID,
			"name":    user.Name,
			"email":   user.Email,
			"role":    user.Role,
		},
	})
}

// --------------------------------------------------
// Forgot Password
// --------------------------------------------------

// ForgotPassword godoc
// @Summary Request password reset
// @Description Create a password reset token for a user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ForgotPasswordRequest true "Forgot password request"
// @Success 200 {object} map[string]interface{}
// @Router /auth/forgot-password [post]

func (a *AuthController) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Find user
	user, err := a.userRepo.FindByEmail(req.Email)

	// Do not reveal whether email exists
	if err != nil || user == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "If the email exists, a password reset link has been created.",
		})
		return
	}

	// Google-only accounts don't necessarily have a password
	if user.PasswordHash == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "If the email exists, a password reset link has been created.",
		})
		return
	}

	// Create reset token
	resetToken, err := a.authService.CreatePasswordResetToken(
		user.UserID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// For development/testing only.
	// Later this token should be sent through email.
	c.JSON(http.StatusOK, gin.H{
		"message":     "If the email exists, a password reset link has been created.",
		"reset_token": resetToken.Token,
		"expires_at":  resetToken.ExpiresAt,
	})
}

// --------------------------------------------------
// Reset Password
// --------------------------------------------------

func (a *AuthController) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Find reset token
	resetToken, err := a.authService.FindPasswordResetToken(req.Token)

	if err != nil || resetToken == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid or expired reset token",
		})
		return
	}

	// Check expiration
	if time.Now().After(resetToken.ExpiresAt) {
		_ = a.authService.DeletePasswordResetToken(req.Token)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid or expired reset token",
		})
		return
	}

	// Hash new password
	passwordHash, err := services.HashPassword(req.NewPassword)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	// Update password
	if err := a.userRepo.UpdatePassword(
		resetToken.UserID,
		passwordHash,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update password",
		})
		return
	}

	// Delete reset token
	if err := a.authService.DeletePasswordResetToken(req.Token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete reset token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "password reset successful",
	})
}

// --------------------------------------------------
// Google OAuth - Start login (redirect to Google)
// --------------------------------------------------

// GoogleLogin godoc
// @Summary Start Google OAuth login
// @Description Redirects the browser to Google's consent screen
// @Tags Auth
// @Router /auth/google [get]

func (a *AuthController) GoogleLogin(c *gin.Context) {

	state := generateRandomState()

	// เก็บ state ไว้ใน cookie สั้นๆ เพื่อตรวจสอบตอน callback กลับมา
	// (ป้องกัน CSRF พื้นฐาน)
	c.SetCookie(
		"google_oauth_state",
		state,
		300, // 5 นาที
		"/",
		"",
		false,
		true,
	)

	authURL := config.GoogleOAuthConfig.AuthCodeURL(
		state,
		oauth2.AccessTypeOnline,
	)

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// --------------------------------------------------
// Google OAuth - Callback
// --------------------------------------------------

// GoogleCallback godoc
// @Summary Google OAuth callback
// @Description Handles the redirect back from Google, issues a JWT
// @Tags Auth
// @Router /auth/google/callback [get]

func (a *AuthController) GoogleCallback(c *gin.Context) {

	// ------------------------------------------------
	// ตรวจสอบ state (CSRF check)
	// ------------------------------------------------

	expectedState, err := c.Cookie("google_oauth_state")

	if err != nil || expectedState == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing oauth state, please try logging in again",
		})
		return
	}

	c.SetCookie("google_oauth_state", "", -1, "/", "", false, true)

	if c.Query("state") != expectedState {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid oauth state",
		})
		return
	}

	code := c.Query("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing authorization code",
		})
		return
	}

	// ------------------------------------------------
	// แลก code เป็น access token กับ Google
	// ------------------------------------------------

	token, err := config.GoogleOAuthConfig.Exchange(
		c.Request.Context(),
		code,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to exchange authorization code",
		})
		return
	}

	// ------------------------------------------------
	// ดึงข้อมูลโปรไฟล์จาก Google
	// ------------------------------------------------

	client := config.GoogleOAuthConfig.Client(
		c.Request.Context(),
		token,
	)

	resp, err := client.Get(
		"https://www.googleapis.com/oauth2/v2/userinfo",
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch google profile",
		})
		return
	}

	defer resp.Body.Close()

	var profile struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to parse google profile",
		})
		return
	}

	if profile.Email == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "google account has no email",
		})
		return
	}

	profile.Email = strings.ToLower(strings.TrimSpace(profile.Email))

	// ------------------------------------------------
	// หา user เดิม หรือสร้างใหม่
	// ------------------------------------------------

	user, err := a.userRepo.FindByGoogleID(profile.ID)

	if err != nil {

		// ยังไม่เคย login ด้วย Google มาก่อน
		// เช็คว่า email นี้เคยสมัครด้วย password ไว้แล้วหรือไม่
		existingByEmail, emailErr := a.userRepo.FindByEmail(
			profile.Email,
		)

		if emailErr == nil && existingByEmail != nil {

			// ผูก Google เข้ากับบัญชีเดิม
			if err := a.userRepo.LinkGoogleID(
				existingByEmail.UserID,
				profile.ID,
			); err != nil {

				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "failed to link google account",
				})
				return
			}

			googleID := profile.ID
			existingByEmail.GoogleID = &googleID
			user = existingByEmail

		} else {

			// user ใหม่เอี่ยม ไม่มี password เลย
			googleID := profile.ID

			newUser := models.User{
				UserID:   uuid.NewString(),
				Name:     profile.Name,
				Email:    profile.Email,
				Role:     "user",
				GoogleID: &googleID,
			}

			if profile.Picture != "" {
				picture := profile.Picture
				newUser.ProfileImage = &picture
			}

			if err := a.userRepo.CreateUser(newUser); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "failed to create user",
				})
				return
			}

			user = &newUser
		}
	}

	// ------------------------------------------------
	// ออก JWT เหมือน login ปกติ
	// ------------------------------------------------

	appToken, err := a.authService.GenerateToken(
		user.UserID,
		user.Role,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return
	}

	// ------------------------------------------------
	// Redirect กลับไปหน้า frontend พร้อม token
	// ------------------------------------------------

	frontendURL := os.Getenv("FRONTEND_URL")

	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}

	redirectURL := frontendURL +
		"/oauth/callback?token=" +
		url.QueryEscape(appToken)

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// --------------------------------------------------
// Helper: random state string for CSRF protection
// --------------------------------------------------

func generateRandomState() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
