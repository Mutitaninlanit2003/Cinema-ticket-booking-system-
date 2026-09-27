package services

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"cinema-ticket-booking/models"
	"cinema-ticket-booking/repository"
)

type SeatService struct {
	Repo *repository.SeatRepository
}

func NewSeatService() *SeatService {
	return &SeatService{
		Repo: repository.NewSeatRepository(),
	}
}

// GenerateSeatsRequest represents the request for generating multiple seats.
type GenerateSeatsRequest struct {
	ScreenID    string `json:"screen_id" binding:"required,uuid"`
	StartRow    string `json:"start_row" binding:"required"`
	EndRow      string `json:"end_row" binding:"required"`
	SeatsPerRow int    `json:"seats_per_row" binding:"required,min=1,max=99"`
	SeatType    string `json:"seat_type" binding:"required,oneof=standard premium vip"`
	Status      string `json:"status"`
}

// Create Seat
func (s *SeatService) CreateSeat(
	seat models.Seat,
) error {

	screen, err := s.Repo.GetScreenByUUID(
		seat.ScreenID,
	)

	if err != nil {
		return err
	}

	seat.ScreenRefID = screen.ID

	now := time.Now()

	seat.CreatedAt = now
	seat.UpdatedAt = now

	return s.Repo.CreateSeat(seat)
}

// Get Seats
func (s *SeatService) GetSeats() ([]models.Seat, error) {
	return s.Repo.GetSeats()
}

// UpdateSeatRequest represents the fields that can be updated for a seat.
type UpdateSeatRequest struct {
	ScreenID   string `json:"screen_id"`
	SeatNumber string `json:"seat_number" binding:"required"`
	SeatType   string `json:"seat_type" binding:"required,oneof=standard premium vip"`
}

// Update Seat
func (s *SeatService) UpdateSeat(
	seatID string,
	req *UpdateSeatRequest,
) error {

	// Normalize values
	req.SeatNumber = strings.ToUpper(
		strings.TrimSpace(req.SeatNumber),
	)

	req.SeatType = strings.ToLower(
		strings.TrimSpace(req.SeatType),
	)

	if req.SeatNumber == "" {
		return errors.New("seat_number is required")
	}

	if req.SeatType != "standard" &&
		req.SeatType != "premium" &&
		req.SeatType != "vip" {
		return errors.New("invalid seat_type")
	}

	// If screen_id is provided, make sure the screen exists.
	if req.ScreenID != "" {
		_, err := s.Repo.GetScreenByUUID(req.ScreenID)
		if err != nil {
			return err
		}
	}

	now := time.Now()

	return s.Repo.UpdateSeat(
		seatID,
		req.SeatNumber,
		req.SeatType,
		now,
	)
}

// Delete Seat
func (s *SeatService) DeleteSeat(
	seatID string,
) error {

	return s.Repo.DeleteSeat(seatID)
}

// Generate multiple seats
func (s *SeatService) GenerateSeats(
	req *GenerateSeatsRequest,
) (int, error) {

	// ---------------------------------------
	// Validate screen
	// ---------------------------------------

	screen, err := s.Repo.GetScreenByUUID(
		req.ScreenID,
	)

	if err != nil {
		return 0, err
	}

	// ---------------------------------------
	// Normalize row input
	// ---------------------------------------

	startRow := strings.ToUpper(
		strings.TrimSpace(req.StartRow),
	)

	endRow := strings.ToUpper(
		strings.TrimSpace(req.EndRow),
	)

	// ---------------------------------------
	// Validate row
	// ---------------------------------------

	if len(startRow) != 1 || len(endRow) != 1 {
		return 0, errors.New(
			"start_row and end_row must be a single letter",
		)
	}

	if startRow[0] < 'A' ||
		startRow[0] > 'Z' ||
		endRow[0] < 'A' ||
		endRow[0] > 'Z' {

		return 0, errors.New(
			"row must be between A and Z",
		)
	}

	if startRow[0] > endRow[0] {
		return 0, errors.New(
			"start_row must not be after end_row",
		)
	}

	// ---------------------------------------
	// Status
	// ---------------------------------------

	status := strings.ToUpper(
		strings.TrimSpace(req.Status),
	)

	if status == "" {
		status = "AVAILABLE"
	}

	// ---------------------------------------
	// Generate seats
	// ---------------------------------------

	seats := make([]*models.Seat, 0)

	now := time.Now()

	for row := startRow[0]; row <= endRow[0]; row++ {

		for number := 1; number <= req.SeatsPerRow; number++ {

			seatNumber := string(row) +
				strconv.Itoa(number)

			seat := &models.Seat{

				ScreenID:    screen.ScreenID,
				ScreenRefID: screen.ID,

				SeatNumber: seatNumber,
				SeatType:   req.SeatType,
				Status:     status,

				CreatedAt: now,
				UpdatedAt: now,
			}

			seats = append(
				seats,
				seat,
			)
		}
	}

	// ---------------------------------------
	// Save all seats
	// ---------------------------------------

	if err := s.Repo.CreateSeatsBulk(
		seats,
	); err != nil {
		return 0, err
	}

	return len(seats), nil
}

// // Delete all seats for a screen
// func (s *SeatService) DeleteSeatsByScreenID(screenID string) error {
// 	return s.Repo.DeleteSeatsByScreenID(screenID)
// }
