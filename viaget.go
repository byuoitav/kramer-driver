package kramer

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// User status constants
const (
	Inactive = "0"
	Active   = "1"
	Waiting  = "2"
)

func (v *Via) GetVolume(ctx context.Context) (int, error) {
	var cmd command
	cmd.Command = "Vol"
	cmd.Param1 = "Get"

	v.Infof("Sending command to get VIA Volume to %s", v.Address)
	// Note: Volume Get command in VIA API doesn't have any error handling so it only returns Vol|Get|XX or nothing
	// Checking for errors during execution of command
	vollevel, err := v.sendCommand(ctx, cmd)
	if err != nil {
		v.Debugf("Failed to get volume from %s: %s", v.Address, err)
		return 0, err
	}

	return v.volumeParse(vollevel)
}

// volumeParse parser to pull out the volume level from the VIA API returned string
func (v *Via) volumeParse(vollevel string) (int, error) {
	re := regexp.MustCompile("[0-9]+")
	vol := re.FindString(vollevel)

	vfin, err := strconv.Atoi(vol)
	if err != nil {
		err = fmt.Errorf("Error converting response: %s", err.Error())
		v.Infof("%s", err.Error())
		return 0, err
	}

	return vfin, nil
}

// isConnected checks the status of the VIA connection
func (v *Via) isConnected(ctx context.Context) bool {
	connected := false

	v.Infof("Getting connected status of %s", v.Address)

	var cmd command
	cmd.Command = "GetSerialNo"

	resp, err := v.sendCommand(ctx, cmd)
	if err == nil && strings.Contains(resp, "GetSerialNo") {
		connected = true
	}

	return connected
}

// Get the Room Code and return the current room code as a string
func (v *Via) GetRoomCode(ctx context.Context) (string, error) {
	var cmd command
	cmd.Command = "RCode"
	cmd.Param1 = "Get"
	cmd.Param2 = "Code"

	v.Infof("Sending command to get current room code to %s", v.Address)
	// Note: RCode Get command in VIA API doesn't have any error handling so it only returns RCode|Get|Code|XXXX or nothing
	resp, err := v.sendCommand(ctx, cmd)
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

//GetPresenterCount .
func (v *Via) GetPresenterCount(ctx context.Context) (int, error) {
	var cmd command
	cmd.Command = "PList"
	cmd.Param1 = "all"
	cmd.Param2 = "1"

	v.Logger.Infof("Sending command to get VIA Presentation count to %s", v.Address)
	// Check for errors on execution of the command was sent correctly and no error occured during send
	resp, err := v.sendCommand(ctx, cmd)
	if err != nil {
		v.Debugf("Error in retrieving volume settings on ")
		return 0, err
	}

	firstsplit := strings.Split(resp, "|")
	//check to assert that first split is len 4
	if len(firstsplit) != 4 {
		return 0, fmt.Errorf("Unknown response %v", resp)
	}

	if strings.Contains(strings.ToLower(firstsplit[3]), "error14") {
		return 0, nil
	}

	//otherwise we go through and split on #, then count the number of
	secondSplit := strings.Split(firstsplit[3], "$")

	return len(secondSplit), nil
}

// GetHardwareInfo for a VIA device
func (v *Via) GetHardwareInfo(ctx context.Context) (HardwareInfo, error) {
	v.Infof("Getting hardware info of %s", v.Address)

	var toReturn HardwareInfo
	var cmd command

	// get serial number
	cmd.Command = "GetSerialNo"

	serial, err := v.sendCommand(ctx, cmd)
	if err != nil {
		return toReturn, fmt.Errorf("failed to get serial number from %s", v.Address)
	}

	toReturn.SerialNumber = parseResponse(serial, "|")

	// get firmware version
	cmd.Command = "GetVersion"

	version, err := v.sendCommand(ctx, cmd)
	if err != nil {
		return toReturn, fmt.Errorf("failed to get the firmware version of %s", v.Address)
	}

	toReturn.FirmwareVersion = parseResponse(version, "|")

	// get MAC address
	cmd.Command = "GetMacAdd"

	macAddr, err := v.sendCommand(ctx, cmd)
	if err != nil {
		return toReturn, fmt.Errorf("failed to get the MAC address of %s", v.Address)
	}

	// get IP information
	cmd.Command = "IpInfo"

	ipInfo, err := v.sendCommand(ctx, cmd)
	if err != nil {
		return toReturn, fmt.Errorf("failed to get the IP information from %s", v.Address)
	}

	hostname, network := parseIPInfo(ipInfo)

	toReturn.Hostname = hostname
	network.MACAddress = parseResponse(macAddr, "|")
	toReturn.NetworkInfo = network

	return toReturn, nil
}

func parseResponse(resp string, delimiter string) string {
	pieces := strings.Split(resp, delimiter)

	var msg string

	if len(pieces) < 2 {
		msg = pieces[0]
	} else {
		msg = pieces[1]
	}

	return strings.Trim(msg, "\r\n")
}

func parseIPInfo(ip string) (hostname string, network NetworkInfo) {
	ipList := strings.Split(ip, "|")

	for _, item := range ipList {
		if strings.Contains(item, "IP") {
			network.IPAddress = strings.Split(item, ":")[1]
		}
		if strings.Contains(item, "GAT") {
			network.Gateway = strings.Split(item, ":")[1]
		}
		if strings.Contains(item, "DNS") {
			network.DNS = []string{strings.Split(item, ":")[1]}
		}
		if strings.Contains(item, "Host") {
			hostname = strings.Trim(strings.Split(item, ":")[1], "\r\n")
		}
	}

	return hostname, network
}

// GetActiveSignal determines the active signal of the VIA by getting the user count
func (v *Via) GetActiveSignal(ctx context.Context) (ActiveSignal, error) {
	signal := ActiveSignal{Active: false}

	count, err := v.GetPresenterCount(ctx)
	if err != nil {
		return signal, err
	}

	if count > 0 {
		signal.Active = true
	}

	return signal, nil
}

// getStatusOfUsers returns the status of users that are logged in to the VIA
func (v *Via) GetStatusOfUsers(ctx context.Context) (VIAUsers, error) {
	var toReturn VIAUsers
	toReturn.InactiveUsers = []string{}
	toReturn.ActiveUsers = []string{}
	toReturn.UsersWaiting = []string{}

	var cmd command
	cmd.Command = "PList"
	cmd.Param1 = "all"
	cmd.Param2 = "4"

	v.Infof("Sending command to get VIA users info to %s", v.Address)

	response, err := v.sendCommand(ctx, cmd)
	if err != nil {
		return toReturn, err
	}

	fullList := strings.Split(response, "|")

	userList := strings.Split(fullList[3], "#")

	for _, user := range userList {
		if len(user) == 0 {
			continue
		}

		userSplit := strings.Split(user, "_")

		if len(userSplit) < 2 {
			continue
		}

		nickname := userSplit[0]
		state := userSplit[1]

		switch state {
		case Inactive:
			toReturn.InactiveUsers = append(toReturn.InactiveUsers, nickname)
			break
		case Active:
			toReturn.ActiveUsers = append(toReturn.ActiveUsers, nickname)
			break
		case Waiting:
			toReturn.UsersWaiting = append(toReturn.UsersWaiting, nickname)
			break
		}
	}

	return toReturn, nil
}

// Get the Room Code and return the current room code as a string
func (v *Via) RoomCode(ctx context.Context) (string, error) {
	var cmd command
	cmd.Command = "RCode"
	cmd.Param1 = "Get"
	cmd.Param2 = "Code"

	v.Infof("Sending command to get current room code to %s", v.Address)
	// Note: RCode Get command in VIA API doesn't have any error handling so it only returns RCode|Get|Code|XXXX or nothing
	resp, err := v.sendCommand(ctx, cmd)
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
