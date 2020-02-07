package via

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/kramer-driver/via"
)

const REBOOT = "Reboot"
const RESET = "Reset"

func (v *VIA) Reboot(ctx context.Context) error {
	var command Command
	command.Command = REBOOT

	log.L.Infof("Sending command %s to %s", REBOOT, v.Address)

	_, err := SendCommand(command, v.Address)
	if err != nil {
		return err
	}

	return nil
}

func (v *VIA) Reset(ctx context.Context) error {
	var command Command
	command.Command = RESET

	log.L.Infof("Sending command %s to %s", RESET, v.Address)

	resp, err := SendCommand(command, v.Address)
	if err != nil {
		return err
	}

	if strings.Contains(resp, RESET) && strings.Contains(resp, "1") {
		return nil
	}

	return errors.New(fmt.Sprintf("Incorrect response for command. (Response: %s)", resp))
}

func (v *VIA) SetVolume(ctx context.Context, address string, volumec string) (string, error) {
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