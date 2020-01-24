package via

import (
	"net/http"

	"github.com/byuoitav/common/log"
	"github.com/labstack/echo"
)

func ResetVia(context echo.Context) error {

	address := context.Param("address")

	err := Reset(address)
	if err != nil {
		log.L.Debugf("There was a problem: %v", err.Error())
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	log.L.Debugf("Success.")

	return context.JSON(http.StatusOK, "Success")
}
