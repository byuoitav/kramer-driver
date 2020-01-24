package via

import (
	"errors"
	"fmt"
	"strings"

	"github.com/byuoitav/common/log"

	"github.com/fatih/color"
)

const REBOOT = "Reboot"
const RESET = "Reset"

func Reboot(address string) error {
	defer color.Unset()
	color.Set(color.FgYellow)

	var command Command
	command.Command = REBOOT

	log.L.Infof("Sending command %s to %s", REBOOT, address)

	_, err := SendCommand(command, address)
	if err != nil {
		return err
	}

	return nil
}

func Reset(address string) error {
	defer color.Unset()
	color.Set(color.FgYellow)

	var command Command
	command.Command = RESET

	log.L.Infof("Sending command %s to %s", RESET, address)

	resp, err := SendCommand(command, address)
	if err != nil {
		return err
	}

	if strings.Contains(resp, RESET) && strings.Contains(resp, "1") {
		return nil
	}

	return errors.New(fmt.Sprintf("Incorrect response for command. (Response: %s)", resp))
}

func SetVolume(address string, volumec string) (string, error) {
	defer color.Unset()
	color.Set(color.FgYellow)

	var command Command
	command.Command = "Vol"
	command.Param1 = "Set"
	command.Param2 = volumec

	log.L.Infof("Sending volume set command to %s", address)

	resp, err := SendCommand(command, address)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error in setting volume on %s", address))
	}

	return resp, nil

}
