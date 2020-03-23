package via

import (
	"context"
	"fmt"
	"strings"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/common/structs"
)

// User status constants
const (
	Inactive = "0"
	Active   = "1"
	Waiting  = "2"
)

// IsConnected checks the status of the VIA connection
func (v *VIA) IsConnected(ctx context.Context) bool {
	connected := false

	log.L.Infof("Getting connected status of %s", v.Address)

	var command Command
	command.Command = "GetSerialNo"

	resp, err := v.SendCommand(ctx, command)
	if err == nil && strings.Contains(resp, "GetSerialNo") {
		connected = true
	}

	return connected
}

// Get the Room Code and return the current room code as a string
func (v *VIA) GetRoomCode(ctx context.Context) (string, error) {
	var command Command
	command.Command = "RCode"
	command.Param1 = "Get"
	command.Param2 = "Code"

	log.L.Infof("Sending command to get current room code to %s", v.Address)
	// Note: RCode Get command in VIA API doesn't have any error handling so it only returns RCode|Get|Code|XXXX or nothing
	resp, err := v.SendCommand(ctx, command)
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
func (v *VIA) GetPresenterCount(ctx context.Context) (int, error) {
	var command Command
	command.Command = "PList"
	command.Param1 = "all"
	command.Param2 = "1"

	log.L.Infof("Sending command to get VIA Presentation count to %s", v.Address)
	// Note: Volume Get command in VIA API doesn't have any error handling so it only returns Vol|Get|XX or nothing
	// I am still checking for errors just in case something else fails during execution
	resp, err := v.SendCommand(ctx, command)
	if err != nil {
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

// GetVolume for a VIA device
func (v *VIA) GetVolume(ctx context.Context) (int, error) {

	var command Command
	command.Command = "Vol"
	command.Param1 = "Get"

	log.L.Infof("Sending command to get VIA Volume to %s", v.Address)
	// Note: Volume Get command in VIA API doesn't have any error handling so it only returns Vol|Get|XX or nothing
	// I am still checking for errors just in case something else fails during execution
	vollevel, _ := v.SendCommand(ctx, command)

	return VolumeParse(vollevel)
}

// GetHardwareInfo for a VIA device
func (v *VIA) GetHardwareInfo(ctx context.Context) (structs.HardwareInfo, error) {
	log.L.Infof("Getting hardware info of %s", v.Address)

	var toReturn structs.HardwareInfo
	var command Command

	// get serial number
	command.Command = "GetSerialNo"

	serial, err := v.SendCommand(ctx, command)
	if err != nil {
		return toReturn, nerr.Translate(err).Addf("failed to get serial number from %s", v.Address)
	}

	toReturn.SerialNumber = parseResponse(serial, "|")

	// get firmware version
	command.Command = "GetVersion"

	version, err := v.SendCommand(ctx, command)
	if err != nil {
		return toReturn, nerr.Translate(err).Addf("failed to get the firmware version of %s", v.Address)
	}

	toReturn.FirmwareVersion = parseResponse(version, "|")

	// get MAC address
	command.Command = "GetMacAdd"

	macAddr, err := v.SendCommand(ctx, command)
	if err != nil {
		return toReturn, nerr.Translate(err).Addf("failed to get the MAC address of %s", v.Address)
	}

	// get IP information
	command.Command = "IpInfo"

	ipInfo, err := v.SendCommand(ctx, command)
	if err != nil {
		return toReturn, nerr.Translate(err).Addf("failed to get the IP information from %s", v.Address)
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

func parseIPInfo(ip string) (hostname string, network structs.NetworkInfo) {
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
func (v *VIA) GetActiveSignal(ctx context.Context) (structs.ActiveSignal, error) {
	signal := structs.ActiveSignal{Active: false}

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
func (v *VIA) GetStatusOfUsers(ctx context.Context) (structs.VIAUsers, error) {
	var toReturn structs.VIAUsers
	toReturn.InactiveUsers = []string{}
	toReturn.ActiveUsers = []string{}
	toReturn.UsersWaiting = []string{}

	var command Command
	command.Command = "PList"
	command.Param1 = "all"
	command.Param2 = "4"

	log.L.Infof("Sending command to get VIA users info to %s", v.Address)

	response, err := v.SendCommand(ctx, command)
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
