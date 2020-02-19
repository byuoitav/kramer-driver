package via

import (
	"context"
	"net/http"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/kramer-microservice/via"
)

// GetViaActiveSignal returns the status of users that are logged in to the VIA
func (v *VIA) GetViaActiveSignal(ctx context.Context) error {
	signal, err := via.GetActiveSignal(v.Address)
	if err != nil {
		log.L.Errorf("Failed to retrieve VIA active signal: %s", err.Error())
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusOK, signal)
}

// GetViaRoomCode - Get the room code of a VIA and return it per request
func (v *VIA) GetViaRoomCode(ctx context.Context) error {
	code, err := via.GetRoomCode(v.Address)
	if err != nil {
		log.L.Errorf("Failed to retrieve VIA room code: %s", err.Error())
		return context.JSON(http.StatusInternalServerError, err)
	}
	return context.JSON(http.StatusOK, code)
}

// Get a list of all connected users to the VIA
func (v *VIA) GetConnectedUsers(ctx context.Context) error {
	userlist, err := via.GetStatusOfUsers(v.Address)
	if err != nil {
		log.L.Errorf("Failed to retrieve current user list: %s", err.Error())
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusOK, userlist)
}
