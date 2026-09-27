package services

import (
	"context"
	"encoding/json"
	"log"

	"cinema-ticket-booking/config"
)

type BookingSuccessEvent struct {
	Event       string   `json:"event"`
	BookingID   string   `json:"booking_id"`
	UserID      string   `json:"user_id"`
	ShowtimeID  string   `json:"showtime_id"`
	SeatNumbers []string `json:"seat_numbers"`
}

type MessageQueueService struct {
}

func NewMessageQueueService() *MessageQueueService {
	return &MessageQueueService{}
}

// PublishBookingSuccess sends a booking success event
// through Redis Pub/Sub.
func (s *MessageQueueService) PublishBookingSuccess(
	ctx context.Context,
	event BookingSuccessEvent,
) error {

	if config.RedisClient == nil {
		return nil
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = config.RedisClient.Publish(
		ctx,
		"booking.success",
		data,
	).Err()

	if err != nil {
		log.Println("❌ Redis Publish failed:", err)
		return err
	}

	log.Println("📤 Redis Pub/Sub message published to booking.success")

	return nil

}

// StartConsumer listens for booking success events.
func (s *MessageQueueService) StartConsumer(
	ctx context.Context,
) {

	if config.RedisClient == nil {
		log.Println("⚠️ Redis client is not initialized")
		return
	}

	pubsub := config.RedisClient.Subscribe(
		ctx,
		"booking.success",
	)

	log.Println("📡 Message Queue Consumer started")
	log.Println("📡 Listening on channel: booking.success")

	go func() {

		defer pubsub.Close()

		for {

			msg, err := pubsub.ReceiveMessage(ctx)

			if err != nil {

				if ctx.Err() != nil {
					return
				}

				log.Println("❌ Message Queue error:", err)
				continue
			}

			log.Println("📨 Booking Success Event received:")
			log.Println(msg.Payload)

			// Mock notification / async processing
			log.Println("🔔 Mock Notification: Booking success processed")
		}
	}()
}
