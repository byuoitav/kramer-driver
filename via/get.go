package via

import (
	"fmt"
	"strings"

	"github.com/byuoitav/common/log"

	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/common/structs"
	"github.com/fatih/color"
)

// User status constants
const (
	Inactive = "0"
	Active   = "1"
	Waiting  = "2"
)

// IsConnected checks the status of the VIA connection
func IsConnected(address string) bool {
	defer color.Unset()
	color.Set(color.FgYellow)
	connected := false

	log.L.Infof("Getting connected status of %s", address)

	var command Command
	resp, err := SendCommand(command, address)
	if err == nil && strings.Contains(resp, "Successful") {
		connected = true
	}

	return connected
}

// Get the Room Code and return the current room code as a string
func GetRoomCode(address string) (string, error) {
	var command Command
	command.Command = "RCode"
	command.Param1 = "Get"
	command.Param2 = "Code"

	log.L.Infof("Sending command to get current room code to %s", address)
	// Note: RCode Get command in VIA API doesn't have any error handling so it only returns RCode|Get|Code|XXXX or nothing
	resp, err := SendCommand(command, address)
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
func GetPresenterCount(address string) (int, error) {
	var command Command
	command.Command = "PList"
	command.Param1 = "all"
	command.Param2 = "1"

	log.L.Infof("Sending command to get VIA Presentation count to %s", address)
	// Note: Volume Get command in VIA API doesn't have any error handling so it only returns Vol|Get|XX or nothing
	// I am still checking for errors just in case something else fails during execution
	resp, err := SendCommand(command, address)
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
func GetVolume(address string) (int, error) {

	defer color.Unset()
	color.Set(color.FgYellow)

	var command Command
	command.Command = "Vol"
	command.Param1 = "Get"

	log.L.Infof("Sending command to get VIA Volume to %s", address)
	// Note: Volume Get command in VIA API doesn't have any error handling so it only returns Vol|Get|XX or nothing
	// I am still checking for errors just in case something else fails during execution
	vollevel, _ := SendCommand(command, address)

	return VolumeParse(vollevel)
}

// GetHardwareInfo for a VIA device
func GetHardwareInfo(address string) (structs.HardwareInfo, *nerr.E) {
	defer color.Unset()
	color.Set(color.FgYellow)

	log.L.Infof("Getting hardware info of %s", address)

	var toReturn structs.HardwareInfo
	var command Command

	// get serial number
	command.Command = "GetSerialNo"

	serial, err := SendCommand(command, address)
	if err != nil {
		return toReturn, nerr.Translate(err).Addf("failed to get serial number from %s", address)
	}

	toReturn.SerialNumber = parseResponse(serial, "|")

	// get firmware version
	command.Command = "GetVersion"

	version, err := SendCommand(command, address)
	if err != nil {
		return toReturn, nerr.Translate(err).Addf("failed to get the firmware version of %s", address)
	}

	toReturn.FirmwareVersion = parseResponse(version, "|")

	// get MAC address
	command.Command = "GetMacAdd"

	macAddr, err := SendCommand(command, address)
	if err != nil {
		return toReturn, nerr.Translate(err).Addf("failed to get the MAC address of %s", address)
	}

	// get IP information
	command.Command = "IpInfo"

	ipInfo, err := SendCommand(command, address)
	if err != nil {
		return toReturn, nerr.Translate(err).Addf("failed to get the IP information from %s", address)
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
func GetActiveSignal(address string) (structs.ActiveSignal, *nerr.E) {
	signal := structs.ActiveSignal{Active: false}

	count, err := GetPresenterCount(address)
	if err != nil {
		return signal, nerr.Translate(err).Add("failed to get the status of users")
	}

	if count > 0 {
		signal.Active = true
	}

	return signal, nil
}

// getStatusOfUsers returns the status of users that are logged in to the VIA
func GetStatusOfUsers(address string) (structs.VIAUsers, *nerr.E) {
	var toReturn structs.VIAUsers
	toReturn.InactiveUsers = []string{}
	toReturn.ActiveUsers = []string{}
	toReturn.UsersWaiting = []string{}

	defer color.Unset()
	color.Set(color.FgYellow)

	var command Command
	command.Command = "PList"
	command.Param1 = "all"
	command.Param2 = "4"

	log.L.Infof("Sending command to get VIA users info to %s", address)

	response, err := SendCommand(command, address)
	if err != nil {
		return toReturn, nerr.Translate(err).Addf("failed to get user information from %s", address)
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
