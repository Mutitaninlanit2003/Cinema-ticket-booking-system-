package services

import (
	"context"
	"errors"

	"cinema-ticket-booking/config"
	"cinema-ticket-booking/models"
)

const (
	SeatAvailable = "AVAILABLE"
	SeatLocked    = "LOCKED"
	SeatBooked    = "BOOKED"
)

type SeatStatusService struct {
	seatLock *SeatLockService
}

func NewSeatStatusService(
	seatLock *SeatLockService,
) *SeatStatusService {
	return &SeatStatusService{
		seatLock: seatLock,
	}
}

func (s *SeatStatusService) GetSeatStatus(
	ctx context.Context,
	showtimeID string,
	seat models.Seat,
) (string, error) {

	if config.RedisClient == nil {
		return "", errors.New("redis client is not initialized")
	}

	key := s.seatLock.lockKey(showtimeID, seat.SeatID)

	// ตรวจว่ามีคน Lock อยู่หรือไม่
	lockedUserID, err := config.RedisClient.Get(ctx, key).Result()

	if err == nil && lockedUserID != "" {
		return SeatLocked, nil
	}

	// ถ้าไม่มี Lock ตอนนี้ให้ถือว่า AVAILABLE
	// ส่วน BOOKED จะตรวจจาก Booking ในขั้นถัดไป
	return SeatAvailable, nil
}
