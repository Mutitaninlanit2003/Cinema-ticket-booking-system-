package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"cinema-ticket-booking/services"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type SeatMapController struct {
	service *services.SeatMapService
}

func NewSeatMapController(
	service *services.SeatMapService,
) *SeatMapController {
	return &SeatMapController{
		service: service,
	}
}

var seatMapUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type SeatMapRequest struct {
	Action string `json:"action"`
	SeatID string `json:"seat_id"`
	UserID string `json:"user_id"`
}

func (c *SeatMapController) WebSocket(ctx *gin.Context) {

	showtimeID := ctx.Param("showtime_id")

	conn, err := seatMapUpgrader.Upgrade(
		ctx.Writer,
		ctx.Request,
		nil,
	)

	if err != nil {
		return
	}

	defer conn.Close()

	client := &services.SeatMapClient{
		Conn: conn,
	}

	c.service.AddClient(showtimeID, client)

	defer c.service.RemoveClient(
		showtimeID,
		conn,
	)

	// แจ้งว่าเชื่อมต่อสำเร็จ
	_ = conn.WriteJSON(
		services.SeatMapMessage{
			Type:       "connected",
			ShowtimeID: showtimeID,
			Message:    "Seat Map WebSocket connected",
		},
	)

	for {

		_, data, err := conn.ReadMessage()

		if err != nil {
			break
		}

		var request SeatMapRequest

		if err := json.Unmarshal(data, &request); err != nil {
			_ = conn.WriteJSON(
				services.SeatMapMessage{
					Type:    "error",
					Message: "invalid JSON message",
				},
			)

			continue
		}

		client.UserID = request.UserID

		switch request.Action {

		case "lock":

			err := c.service.LockSeat(
				context.Background(),
				showtimeID,
				request.SeatID,
				request.UserID,
			)

			if err != nil {
				_ = conn.WriteJSON(
					services.SeatMapMessage{
						Type:    "error",
						SeatID:  request.SeatID,
						Message: err.Error(),
					},
				)
			}

		case "unlock":

			err := c.service.UnlockSeat(
				context.Background(),
				showtimeID,
				request.SeatID,
				request.UserID,
			)

			if err != nil {
				_ = conn.WriteJSON(
					services.SeatMapMessage{
						Type:    "error",
						SeatID:  request.SeatID,
						Message: err.Error(),
					},
				)
			}

		default:

			_ = conn.WriteJSON(
				services.SeatMapMessage{
					Type:    "error",
					Message: "unknown action",
				},
			)
		}
	}
}
