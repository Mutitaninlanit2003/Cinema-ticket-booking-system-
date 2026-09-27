package services

import (
	"context"
	"errors"
	"time"

	"cinema-ticket-booking/config"

	"github.com/redis/go-redis/v9"
)

var (
	ErrSeatLocked      = errors.New("seat is already locked")
	ErrSeatLockOwner   = errors.New("you do not own this seat lock")
	ErrSeatLockExpired = errors.New("seat lock has expired")
)

type SeatLockService struct {
	ttl time.Duration
}

func NewSeatLockService() *SeatLockService {
	return &SeatLockService{
		ttl: 15 * time.Minute,
	}
}

func (s *SeatLockService) lockKey(showtimeID, seatID string) string {
	return "seat_lock:" + showtimeID + ":" + seatID
}

func (s *SeatLockService) LockSeat(
	ctx context.Context,
	showtimeID string,
	seatID string,
	userID string,
) error {

	if config.RedisClient == nil {
		return errors.New("redis client is not initialized")
	}

	key := s.lockKey(showtimeID, seatID)

	locked, err := config.RedisClient.SetNX(
		ctx,
		key,
		userID,
		s.ttl,
	).Result()

	if err != nil {
		return err
	}

	if !locked {
		return ErrSeatLocked
	}

	return nil
}

func (s *SeatLockService) UnlockSeat(
	ctx context.Context,
	showtimeID string,
	seatID string,
	userID string,
) error {

	if config.RedisClient == nil {
		return errors.New("redis client is not initialized")
	}

	key := s.lockKey(showtimeID, seatID)

	currentUserID, err := config.RedisClient.Get(
		ctx,
		key,
	).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrSeatLockExpired
		}

		return err
	}

	if currentUserID != userID {
		return ErrSeatLockOwner
	}

	return config.RedisClient.Del(ctx, key).Err()
}

func (s *SeatLockService) VerifySeatLock(
	ctx context.Context,
	showtimeID string,
	seatID string,
	userID string,
) error {

	if config.RedisClient == nil {
		return errors.New("redis client is not initialized")
	}

	key := s.lockKey(showtimeID, seatID)

	currentUserID, err := config.RedisClient.Get(
		ctx,
		key,
	).Result()

	if err != nil {
		if err == redis.Nil {
			return ErrSeatLockExpired
		}

		return err
	}

	if currentUserID != userID {
		return ErrSeatLockOwner
	}

	return nil
}
