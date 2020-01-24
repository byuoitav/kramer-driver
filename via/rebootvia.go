package via

import (
	"net/http"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/kramer-microservice/via"
	"github.com/labstack/echo"
)

func RebootVia(context echo.Context) error {
	address := context.Param("address")

	err := via.Reboot(address)
	if err != nil {
		log.L.Debugf("There was a problem: %v", err.Error())
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	log.L.Debugf("Success.")

	return context.JSON(http.StatusOK, "Success")
}
