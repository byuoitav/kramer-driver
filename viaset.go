package kramer

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/byuoitav/common/log"
)

// Reboot: Reboot a VIA using the API
func (v *Via) Reboot(ctx context.Context) error {
	var cmd command
	cmd.Command = viaReboot

	log.L.Infof("Sending command %s to %s", viaReboot, v.Address)

	_, err := v.sendCommand(ctx, cmd)
	if err != nil {
		return err
	}

	return nil
}

// Reset: Reset a VIA sessions - Causes VIAAdmin to log out and log back in which can help with some lock up issues.
func (v *Via) Reset(ctx context.Context) error {
	var cmd command
	cmd.Command = viaReset

	log.L.Infof("Sending command %s to %s", viaReset, v.Address)

	resp, err := v.sendCommand(ctx, cmd)
	if err != nil {
		return err
	}

	if strings.Contains(resp, viaReset) && strings.Contains(resp, "1") {
		return nil
	}

	return errors.New(fmt.Sprintf("Incorrect response for command. (Response: %s)", resp))
}

// SetAlert - Send an alert to the VIA
func (v *Via) SetAlert(ctx context.Context, AlertMessage string) error {
	var cmd command
	cmd.Command = "IAlert"
	cmd.Param1 = AlertMessage
	cmd.Param2 = "0"
	cmd.Param3 = "5"

	log.L.Infof("Sending Alert to %v", v.Address)

	log.L.Debugf("Sending an alert message -%s- to %s", AlertMessage, v.Address)

	resp, err := v.sendCommand(ctx, cmd)
	if err != nil {
		return fmt.Errorf("Error in setting alert on %s", v.Address)
	}
	sp := strings.Split(resp, "|")
	s := sp[1]
	sint := strconv.Atoi(s)
	if sint != 1 {
		return fmt.Errorf("Alert was not successfully sent, checking settings and try again")
	}

	return nil
}

// SetVolume - Used to set the volume on a VIA (Used by both VIA-Control and DSP Driver sets)
func (v *Via) SetVolume(ctx context.Context, volume int) (string, error) {
	var cmd command
	cmd.Command = "Vol"
	cmd.Param1 = "Set"
	cmd.Param2 = strconv.Itoa(volume)

	log.L.Infof("Sending volume set command to %s", v.Address)

	resp, err := v.sendCommand(ctx, cmd)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error in setting volume on %s", v.Address))
	}

	return resp, nil

}
