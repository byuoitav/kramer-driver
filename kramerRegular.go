package kramer

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/byuoitav/common/structs"
)

const (
	BuildDate       = "BUILD-DATE"
	Model           = "MODEL"
	SerialNumber    = "SN"
	FirmwareVersion = "VERSION"
	ProtocolVersion = "PROT-VER"
	Temperature     = "HW-TEMP"
	PowerSave       = "POWER-SAVE"
	IPAddress       = "NET-IP"
	Gateway         = "NET-GATE"
	MACAddress      = "NET-MAC"
	NetDNS          = "NET-DNS"
	Signal          = "SIGNAL"
)

const (
	CARRIAGE_RETURN           = 0x0D
	LINE_FEED                 = 0x0A
	SPACE                     = 0x20
	DELAY_BETWEEN_CONNECTIONS = time.Second * 10
)

type Response struct {
	Response string
	Err      error
}

type CommandInfo struct {
	ResponseChannel chan Response
	Address         string
	Command         string
	ReadWelcome     bool
}

var StartChannel = make(chan CommandInfo, 1000)

//var connMap = make(map[string]chan CommandInfo)

func hardwareCommand(commandType, param, address string, readWelcome bool) (string, error) {
	var command string

	if len(param) > 0 {
		num, _ := strconv.Atoi(param)
		command = fmt.Sprintf("#%s? %d", commandType, num)
	} else {
		command = fmt.Sprintf("#%s?", commandType)
	}

	respChan := make(chan Response)

	c := CommandInfo{respChan, address, command, readWelcome}

	StartChannel <- c

	re := <-respChan

	resp := re.Response
	err := re.Err

	if err != nil {
		return resp, err
	}

	resp = strings.Split(resp, fmt.Sprintf("%s", commandType))[1]
	resp = strings.Trim(resp, "\r\n")
	resp = strings.TrimSpace(resp)

	return resp, nil
}

// GetInput returns the current input
func (k *Kramer) getInputByOutput(ctx context.Context, output string) (string, error) {

	return "", nil
}

// SwitchInput changes the input on the given output to input
func (k *Kramer) setInputByOutput(ctx context.Context, output, input string) error {
	return nil
}

// GetHardwareInfo returns a hardware info struct
func (k *Kramer) getSwitcherHardwareInfo(ctx context.Context) (structs.HardwareInfo, error) {
	var toReturn structs.HardwareInfo
	readWelcome := true
	// get the hostname
	addr, e := net.LookupAddr(k.Address)
	if e != nil {
		toReturn.Hostname = k.Address
	} else {
		toReturn.Hostname = strings.Trim(addr[0], ".")
	}

	// get build date
	buildDate, err := hardwareCommand(BuildDate, "", k.Address, readWelcome)
	if err != nil {
		return toReturn, fmt.Errorf("failed to get build date from %s", k.Address)
	}

	toReturn.BuildDate = buildDate

	// get device model
	model, err := hardwareCommand(Model, "", k.Address, readWelcome)
	if err != nil {
		return toReturn, fmt.Errorf("failed to get model number from %s", k.Address)
	}

	toReturn.ModelName = model

	// get device protocol version
	protocol, err := hardwareCommand(ProtocolVersion, "", k.Address, readWelcome)
	if err != nil {
		return toReturn, fmt.Errorf("failed to get protocol version from %s", k.Address)
	}

	toReturn.ProtocolVersion = strings.Trim(protocol, "3000:")

	// get firmware version
	firmware, err := hardwareCommand(FirmwareVersion, "", k.Address, readWelcome)
	if err != nil {
		return toReturn, fmt.Errorf("failed to get firmware version from %s", k.Address)
	}

	toReturn.FirmwareVersion = firmware

	// get serial number
	serial, err := hardwareCommand(SerialNumber, "", k.Address, readWelcome)
	if err != nil {
		return toReturn, fmt.Errorf("failed to get serial number from %s", k.Address)
	}

	toReturn.SerialNumber = serial

	// get IP address
	ipAddress, err := hardwareCommand(IPAddress, "", k.Address, readWelcome)
	if err != nil {
		return toReturn, fmt.Errorf("failed to get IP address from %s... ironic...", k.Address)
	}

	// get gateway
	gateway, err := hardwareCommand(Gateway, "", k.Address, readWelcome)
	if err != nil {
		return toReturn, fmt.Errorf("failed to get the gateway address from %s", k.Address)
	}

	// get MAC address
	mac, err := hardwareCommand(MACAddress, "", k.Address, readWelcome)
	if err != nil {
		return toReturn, fmt.Errorf("failed to get the MAC address from %s", k.Address)
	}

	// set network information
	toReturn.NetworkInfo = structs.NetworkInfo{
		IPAddress:  ipAddress,
		MACAddress: mac,
		Gateway:    gateway,
	}

	return toReturn, nil
}

//GetInfo .
func (k *Kramer) getInfo(ctx context.Context) (interface{}, error) {
	var info interface{}
	return info, fmt.Errorf("not currently implemented")
}
