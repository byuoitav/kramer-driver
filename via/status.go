package via

import (
	"net/http"
	"strconv"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/status"
	"github.com/byuoitav/kramer-microservice/via"
	"github.com/labstack/echo"
)

// GetViaActiveSignal returns the status of users that are logged in to the VIA
func GetViaActiveSignal(context echo.Context) error {
	signal, err := via.GetActiveSignal(context.Param("address"))
	if err != nil {
		log.L.Errorf("Failed to retrieve VIA active signal: %s", err.Error())
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusOK, signal)
}

// GetViaRoomCode - Get the room code of a VIA and return it per request
func GetViaRoomCode(context echo.Context) error {
	code, err := via.GetRoomCode(context.Param("address"))
	if err != nil {
		log.L.Errorf("Failed to retrieve VIA room code: %s", err.Error())
		return context.JSON(http.StatusInternalServerError, err)
	}
	return context.JSON(http.StatusOK, code)
}

// Get a list of all connected users to the VIA
func GetConnectedUsers(context echo.Context) error {
	userlist, err := via.GetStatusOfUsers(context.Param("address"))
	if err != nil {
		log.L.Errorf("Failed to retrieve current user list: %s", err.Error())
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusOK, userlist)
}
