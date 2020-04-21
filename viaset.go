package kramer

import "fmt"

// Reboot: Reboot a VIA using the API
func (v *Via) Reboot(ctx context.Context) error {
	var command command
	command.Command = viaReboot

	log.L.Infof("Sending command %s to %s", viaReboot, v.Address)

	_, err := v.SendCommand(ctx, command)
	if err != nil {
		return err
	}

	return nil
}

// Reset: Reset a VIA sessions - Causes VIAAdmin to log out and log back in which can help with some lock up issues.
func (v *Via) Reset(ctx context.Context) error {
	var command command
	command.Command = viaReset

	log.L.Infof("Sending command %s to %s", viaReset, v.Address)

	resp, err := v.SendCommand(ctx, command)
	if err != nil {
		return err
	}

	if strings.Contains(resp, viaReset) && strings.Contains(resp, "1") {
		return nil
	}

	return errors.New(fmt.Sprintf("Incorrect response for command. (Response: %s)", resp))
}

// Get the Room Code and return the current room code as a string
func (v *VIA) RoomCode(ctx context.Context) (string, error) {
	var command Command
	command.Command = "RCode"
	command.Param1 = "Get"
	command.Param2 = "Code"

	log.L.Infof("Sending command to get current room code to %s", v.Address)
	// Note: RCode Get command in VIA API doesn't have any error handling so it only returns RCode|Get|Code|XXXX or nothing
	resp, err := v.sendCommand(ctx, command)
	if err != nil {
		return "", err
	}
	split := strings.Split(resp, "|")
	if len(split) != 4 {
		return "", fmt.Errorf("Unknown response %v", resp)
	}

	roomcode := strings.TrimSpace(split[3])

	return roomcode, nil
}

// SetAlert - Send an alert to the VIA
func (v *Via) SetAlert(ctx context.Context, message string) error {
	var command Command
	command.Command = "IAlert"
	command.Param1 = AlertMessage
	command.Param2 = "0"
	command.Param3 = "5"

	log.L.Infof("Sending Alert to %v", v.Address)

	log.L.Debugf("Sending an alert message -%s- to %s", AlertMessage, v.Address)

	resp, err := v.SendCommand(ctx, command)
	if err != nil {
		return fmt.Errorf("Error in setting alert on %s", v.Address)
	}
	sp := strings.Split(resp, "|")
	s := sp[1]
	sint := stings.Atoi(s)
	if sint != 1 {
		return fmt.Errorf("Alert was not successfully sent, checking settings and try again")
	}

	return nil
}

// SetVolume - Used to set the volume on a VIA (Used by both VIA-Control and DSP Driver sets)
func (v *Via) SetVolume(ctx context.Context, volume int) (string, error) {
	var command command
	command.Command = "Vol"
	command.Param1 = "Set"
	command.Param2 = strconv.Itoa(volume)

	log.L.Infof("Sending volume set command to %s", v.Address)

	resp, err := v.SendCommand(ctx, command)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error in setting volume on %s", v.Address))
	}

	return resp, nil

}
