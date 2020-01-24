package via

import (
	"net/http"
	"strconv"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/status"
	"github.com/byuoitav/kramer-microservice/via"
	"github.com/labstack/echo"
)

func SetViaVolume(context echo.Context) error {
	address := context.Param("address")
	value := context.Param("volvalue")
	log.L.Debugf("Value passed by SetViaVolume is %v", value)

	volume, err := strconv.Atoi(value)
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	} else if volume > 100 || volume < 1 {
		log.L.Debugf("Volume command error - volume value %s is outside the bounds of 1-100", value)
		return context.JSON(http.StatusBadRequest, "Error: volume must be a value from 1 to 100!")
	}

	volumec := strconv.Itoa(volume)
	log.L.Debugf("Setting volume for %s to %v...", address, volume)

	response, err := via.SetVolume(address, volumec)

	if err != nil {
		log.L.Debugf("An Error Occured: %s", err)
		return context.JSON(http.StatusBadRequest, "An error has occured while setting volume")
	}
	log.L.Debugf("Success: %s", response)

	return context.JSON(http.StatusOK, status.Volume{Volume: volume})
}

func GetViaVolume(context echo.Context) error {

	address := context.Param("address")

	ViaVolume, err := via.GetVolume(address)

	if err != nil {
		log.L.Debugf("Failed to retreive VIA volume")
		return context.JSON(http.StatusBadRequest, "Failed to retreive VIA volume")
	} else {
		log.L.Debugf("VIA volume is currently set to %v", strconv.Itoa(ViaVolume))
		return context.JSON(http.StatusOK, status.Volume{ViaVolume})
	}

}
