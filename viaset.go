package kramer

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Reboot: Reboot a VIA using the API
func (v *Via) Reboot(ctx context.Context) error {
	var cmd command
	cmd.Command = viaReboot

	v.Infof("Sending command %s to %s", viaReboot, v.Address)

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

	v.Infof("Sending command %s to %s", viaReset, v.Address)

	resp, err := v.sendCommand(ctx, cmd)
	if err != nil {
		return err
	}

	if strings.Contains(resp, viaReset) && strings.Contains(resp, "1") {
		return nil
	}

	return errors.New(fmt.Sprintf("Incorrect response for command. (Response: %s)", resp))
}

// SetAlert - Send an alert popup on the screen of the VIA
func (v *Via) SetAlert(ctx context.Context, AlertMessage string) error {
	v.Debugf("AlertMessage to pass: %s", AlertMessage)

	var cmd command
	cmd.Command = "IAlert"
	cmd.Param1 = AlertMessage
	cmd.Param2 = "0"
	cmd.Param3 = "5"

	v.Infof("Sending Alert to %v", v.Address)

	v.Debugf("Sending an alert message -%s- to %s", AlertMessage, v.Address)

	resp, err := v.sendCommand(ctx, cmd)
	if err != nil {
		v.Errorf("Error in setting alert on %s", v.Address)
		return fmt.Errorf("Error in setting alert on %s", v.Address)
	}
	clean := strings.TrimRight(resp, "\r\n")
	sp := strings.Split(clean, "|")
	s := sp[1]
	sint, err := strconv.Atoi(s)
	if sint != 1 {
		v.Errorf("Alert was not successfully sent, check settings and try again")
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

	v.Infof("Sending volume set command to %s", v.Address)

	resp, err := v.sendCommand(ctx, cmd)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error in setting volume on %s", v.Address))
	}

	return resp, nil

}
