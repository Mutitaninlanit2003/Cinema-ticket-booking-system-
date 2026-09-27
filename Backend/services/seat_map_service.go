package services

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type SeatMapClient struct {
	Conn   *websocket.Conn
	UserID string
}

type SeatMapMessage struct {
	Type       string `json:"type"`
	ShowtimeID string `json:"showtime_id"`
	SeatID     string `json:"seat_id,omitempty"`
	Status     string `json:"status,omitempty"`
	UserID     string `json:"user_id,omitempty"`
	Message    string `json:"message,omitempty"`
}

type SeatMapService struct {
	mu sync.RWMutex

	// showtimeID -> clients
	clients map[string]map[*websocket.Conn]*SeatMapClient

	seatLock *SeatLockService
}

func NewSeatMapService(seatLock *SeatLockService) *SeatMapService {
	return &SeatMapService{
		clients:  make(map[string]map[*websocket.Conn]*SeatMapClient),
		seatLock: seatLock,
	}
}

// AddClient registers a WebSocket client.
func (s *SeatMapService) AddClient(
	showtimeID string,
	client *SeatMapClient,
) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.clients[showtimeID] == nil {
		s.clients[showtimeID] = make(map[*websocket.Conn]*SeatMapClient)
	}

	s.clients[showtimeID][client.Conn] = client
}

// RemoveClient removes a WebSocket client.
func (s *SeatMapService) RemoveClient(
	showtimeID string,
	conn *websocket.Conn,
) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if clients, ok := s.clients[showtimeID]; ok {
		delete(clients, conn)

		if len(clients) == 0 {
			delete(s.clients, showtimeID)
		}
	}
}

// Broadcast sends a message to every client watching the same showtime.
func (s *SeatMapService) Broadcast(
	showtimeID string,
	message SeatMapMessage,
) {
	s.mu.RLock()
	clients := make([]*SeatMapClient, 0)

	for _, client := range s.clients[showtimeID] {
		clients = append(clients, client)
	}

	s.mu.RUnlock()

	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	for _, client := range clients {
		if err := client.Conn.WriteMessage(
			websocket.TextMessage,
			data,
		); err != nil {
			s.RemoveClient(showtimeID, client.Conn)
		}
	}
}

// LockSeat handles a seat lock request.
func (s *SeatMapService) LockSeat(
	ctx context.Context,
	showtimeID string,
	seatID string,
	userID string,
) error {

	err := s.seatLock.LockSeat(
		ctx,
		showtimeID,
		seatID,
		userID,
	)

	if err != nil {
		return err
	}

	s.Broadcast(
		showtimeID,
		SeatMapMessage{
			Type:       "seat_status",
			ShowtimeID: showtimeID,
			SeatID:     seatID,
			Status:     "LOCKED",
			UserID:     userID,
			Message:    "Seat locked",
		},
	)

	return nil
}

// UnlockSeat handles a seat unlock request.
func (s *SeatMapService) UnlockSeat(
	ctx context.Context,
	showtimeID string,
	seatID string,
	userID string,
) error {

	err := s.seatLock.UnlockSeat(
		ctx,
		showtimeID,
		seatID,
		userID,
	)

	if err != nil {
		return err
	}

	s.Broadcast(
		showtimeID,
		SeatMapMessage{
			Type:       "seat_status",
			ShowtimeID: showtimeID,
			SeatID:     seatID,
			Status:     "AVAILABLE",
			UserID:     userID,
			Message:    "Seat available",
		},
	)

	return nil
}
