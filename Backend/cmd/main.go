package main

import (
	"context"
	"log"
	"os"

	"cinema-ticket-booking/config"
	"cinema-ticket-booking/controllers"
	_ "cinema-ticket-booking/docs"
	"cinema-ticket-booking/models"
	"cinema-ticket-booking/repository"
	"cinema-ticket-booking/routes"
	"cinema-ticket-booking/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Cinema Ticket Booking API
// @version 1.0
// @description Cinema Ticket Booking System API
// @host localhost:8080
// @BasePath /
func main() {

	// ==========================
	// Connect Database
	// ==========================

	config.ConnectDB()

	// ==========================
	// Connect Redis
	// ==========================

	config.ConnectRedis()

	// ==========================
	// Auto Migration
	// ==========================

	err := config.DB.AutoMigrate(
		&models.User{},
		&models.Movie{},
		&models.Theater{},
		&models.Screen{},
		&models.Seat{},
		&models.Showtime{},
		&models.Booking{},
		&models.BookingSeat{},
		&models.Payment{},
		&models.AuditLog{},
		&models.PasswordResetToken{},
	)

	if err != nil {
		log.Fatal("❌ AutoMigration failed: ", err)
	}

	log.Println("✅ AutoMigration completed")

	// ==========================
	// Seed Admin User
	// ==========================
	// Runs once at startup, before the router is set up. Idempotent —
	// safe to run on every restart, since it checks for an existing
	// admin by email before creating one. This is the ONLY way an
	// admin account can be created: /auth/register always forces
	// role="user", so there is no public endpoint that can produce
	// a privileged account.

	userRepository := repository.NewUserRepository()

	if err := seedAdminUser(userRepository); err != nil {
		log.Fatal("❌ Admin seed failed: ", err)
	}

	// ==========================
	// Uploads directory
	// ==========================
	// Needed by the avatar-upload endpoint (POST /users/me/avatar) —
	// c.SaveUploadedFile fails if this doesn't already exist on disk.
	// MkdirAll is a no-op if the folder is already there, so this is
	// safe to run on every restart.

	if err := os.MkdirAll("./uploads/avatars", os.ModePerm); err != nil {
		log.Fatal("❌ Failed to create uploads/avatars directory: ", err)
	}

	// ==========================
	// Gin Router
	// ==========================

	router := gin.Default()

	router.Static(
		"/uploads",
		"./uploads",
	)

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			// "http://localhost:5174",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"PATCH",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		},
	}))

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Cinema Ticket Booking API",
		})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	// ==========================
	// User
	// ==========================

	routes.UserRoutes(router)

	// ==========================
	// Google OAuth
	// ==========================

	config.InitGoogleOAuth()

	// ==========================
	// Movie
	// ==========================

	routes.MovieRoutes(router)

	// ==========================
	// Theater
	// ==========================

	routes.TheaterRoutes(router)

	// ==========================
	// Showtime
	// ==========================

	showtimeRepository := repository.NewShowtimeRepository(
		config.DB,
	)

	showtimeService := services.NewShowtimeService(
		showtimeRepository,
	)

	showtimeController := controllers.NewShowtimeController(
		showtimeService,
	)

	routes.ShowtimeRoutes(
		router,
		showtimeController,
	)

	// ==========================
	// Shared Seat Lock Service
	// ==========================

	seatLockService := services.NewSeatLockService()

	messageQueueService := services.NewMessageQueueService()

	messageQueueService.StartConsumer(
		context.Background(),
	)

	// ==========================
	// Booking
	// ==========================

	bookingRepository := repository.NewBookingRepository(
		config.DB,
	)

	bookingService := services.NewBookingService(
		bookingRepository,
		seatLockService,
		messageQueueService,
	)

	bookingController := controllers.NewBookingController(
		bookingService,
	)

	routes.BookingRoutes(
		router,
		bookingController,
	)

	// ==========================
	// Admin
	// ==========================

	adminRepository := repository.NewAdminRepository(
		config.DB,
	)

	adminService := services.NewAdminService(
		adminRepository,
	)

	adminController := controllers.NewAdminController(
		adminService,
	)

	routes.AdminRoutes(
		router,
		adminController,
	)

	// ==========================
	// Audit Log
	// ==========================

	auditLogRepository := repository.NewAuditLogRepository(
		config.DB,
	)

	auditLogService := services.NewAuditLogService(
		auditLogRepository,
	)

	auditLogController := controllers.NewAuditLogController(
		auditLogService,
	)

	routes.AuditLogRoutes(
		router,
		auditLogController,
	)

	// ==========================
	// Payment
	// ==========================

	paymentRepository := repository.NewPaymentRepository(
		config.DB,
	)

	paymentService := services.NewPaymentService(
		paymentRepository,
		bookingService,
	)

	paymentController := controllers.NewPaymentController(
		paymentService,
	)

	routes.PaymentRoutes(
		router,
		paymentController,
	)

	// ==========================
	// Authentication
	// ==========================

	authService, err := services.NewAuthService()

	if err != nil {
		log.Fatal("❌ Auth Service initialization failed: ", err)
	}

	authController := controllers.NewAuthController(
		userRepository,
		authService,
	)

	routes.AuthRoutes(
		router,
		authController,
	)

	// ==========================
	// Screen
	// ==========================

	routes.ScreenRoutes(router)

	// ==========================
	// Seat
	// ==========================

	routes.SeatRoutes(router)

	// ==========================
	// Seat Map
	// ==========================

	seatMapService := services.NewSeatMapService(
		seatLockService,
	)

	seatMapController := controllers.NewSeatMapController(
		seatMapService,
	)

	routes.SeatMapRoutes(
		router,
		seatMapController,
	)

	// ==========================
	// Seat Lock
	// ==========================

	seatLockController := controllers.NewSeatLockController(
		seatLockService,
		seatMapService,
	)

	routes.SeatLockRoutes(
		router,
		seatLockController,
	)

	// ==========================
	// Swagger
	// ==========================

	router.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	// ==========================
	// Start Server
	// ==========================

	log.Println("Server running on port 8080")

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

// seedAdminUser ensures exactly one admin account exists, created
// from ADMIN_EMAIL / ADMIN_PASSWORD in .env. If an account with that
// email already exists, this is a no-op — safe to call on every
// server start.
func seedAdminUser(userRepo *repository.UserRepository) error {
	adminName := "Cinema Admin"
	adminEmail := "admin@cinema.com"
	adminPassword := "admin123456"
	adminRole := "admin"

	// Check whether admin already exists
	existingUser, err := userRepo.FindByEmail(adminEmail)

	if err == nil && existingUser != nil {
		log.Println("ℹ️ Admin user already exists:", adminEmail)
		return nil
	}

	// Hash admin password
	passwordHash, err := services.HashPassword(adminPassword)
	if err != nil {
		return err
	}

	adminUser := models.User{
		UserID:       uuid.NewString(),
		Name:         adminName,
		Email:        adminEmail,
		PasswordHash: &passwordHash,
		Role:         adminRole,
	}

	if err := userRepo.CreateUser(adminUser); err != nil {
		return err
	}

	log.Println("✅ Admin user seeded:", adminEmail)

	return nil
}
