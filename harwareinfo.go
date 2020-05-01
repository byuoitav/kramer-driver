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

func (vs *Kramer4x4) hardwareCommand(ctx context.Context, commandType, param string) (string, error) {
	var cmd []byte

	if len(param) > 0 {
		num, _ := strconv.Atoi(param)
		cmd = []byte(fmt.Sprintf("#%s? %d\r\n", commandType, num))
	} else {
		cmd = []byte(fmt.Sprintf("#%s?\r\n", commandType))
	}

	resp, err := vs.SendCommand(ctx, cmd)

	if err != nil {
		return "", fmt.Errorf("unable to send command: %w", err)
	}
	resps := string(resp)
	resps = strings.Split(resps, fmt.Sprintf("%s", commandType))[1]
	resps = strings.Trim(resps, "\r\n")
	resps = strings.TrimSpace(resps)

	return resps, nil
}

func (vs *Kramer4x4) GetHardwareInfo(ctx context.Context) (structs.HardwareInfo, error) {
	var toReturn structs.HardwareInfo
	// get the hostname
	addr, e := net.LookupAddr(vs.Address)
	if e != nil {
		toReturn.Hostname = vs.Address
	} else {
		toReturn.Hostname = strings.Trim(addr[0], ".")
	}

	// get build date
	buildDate, err := vs.hardwareCommand(ctx, BuildDate, "")
	if err != nil {
		return toReturn, fmt.Errorf("failed to get build date from %s", vs.Address)
	}

	toReturn.BuildDate = buildDate

	// get device model
	model, err := vs.hardwareCommand(ctx, Model, "")
	if err != nil {
		return toReturn, fmt.Errorf("failed to get model number from %s", vs.Address)
	}

	toReturn.ModelName = model

	// get device protocol version
	protocol, err := vs.hardwareCommand(ctx, ProtocolVersion, "")
	if err != nil {
		return toReturn, fmt.Errorf("failed to get protocol version from %s", vs.Address)
	}

	toReturn.ProtocolVersion = strings.Trim(protocol, "3000:")

	// get firmware version
	firmware, err := vs.hardwareCommand(ctx, FirmwareVersion, "")
	if err != nil {
		return toReturn, fmt.Errorf("failed to get firmware version from %s", vs.Address)
	}

	toReturn.FirmwareVersion = firmware

	// get serial number
	serial, err := vs.hardwareCommand(ctx, SerialNumber, "")
	if err != nil {
		return toReturn, fmt.Errorf("failed to get serial number from %s", vs.Address)
	}

	toReturn.SerialNumber = serial

	// get IP address
	ipAddress, err := vs.hardwareCommand(ctx, IPAddress, "")
	if err != nil {
		return toReturn, fmt.Errorf("failed to get IP address from %s... ironic...", vs.Address)
	}

	// get gateway
	gateway, err := vs.hardwareCommand(ctx, Gateway, "")
	if err != nil {
		return toReturn, fmt.Errorf("failed to get the gateway address from %s", vs.Address)
	}

	// get MAC address
	mac, err := vs.hardwareCommand(ctx, MACAddress, "")
	if err != nil {
		return toReturn, fmt.Errorf("failed to get the MAC address from %s", vs.Address)
	}

	// set network information
	toReturn.NetworkInfo = structs.NetworkInfo{
		IPAddress:  ipAddress,
		MACAddress: mac,
		Gateway:    gateway,
	}

	return toReturn, nil
}

func (dsp *KramerAFM20DSP) hardwareCommand(ctx context.Context, commandType, param string) (string, error) {
	var cmd []byte

	if len(param) > 0 {
		num, _ := strconv.Atoi(param)
		cmd = []byte(fmt.Sprintf("#%s? %d\r\n", commandType, num))
	} else {
		cmd = []byte(fmt.Sprintf("#%s?\r\n", commandType))
	}

	resp, err := dsp.SendCommand(ctx, cmd)

	if err != nil {
		return "", fmt.Errorf("unable to send command: %w", err)
	}
	resps := string(resp)
	resps = strings.Split(resps, fmt.Sprintf("%s", commandType))[1]
	resps = strings.Trim(resps, "\r\n")
	resps = strings.TrimSpace(resps)

	return resps, nil
}

func (dsp *KramerAFM20DSP) GetHardwareInfo(ctx context.Context) (structs.HardwareInfo, error) {
	var toReturn structs.HardwareInfo
	// get the hostname
	addr, e := net.LookupAddr(dsp.Address)
	if e != nil {
		toReturn.Hostname = dsp.Address
	} else {
		toReturn.Hostname = strings.Trim(addr[0], ".")
	}

	// get build date
	buildDate, err := dsp.hardwareCommand(ctx, BuildDate, "")
	if err != nil {
		return toReturn, fmt.Errorf("failed to get build date from %s", dsp.Address)
	}

	toReturn.BuildDate = buildDate

	// get device model
	model, err := dsp.hardwareCommand(ctx, Model, "")
	if err != nil {
		return toReturn, fmt.Errorf("failed to get model number from %s", dsp.Address)
	}

	toReturn.ModelName = model

	// get device protocol version
	protocol, err := dsp.hardwareCommand(ctx, ProtocolVersion, "")
	if err != nil {
		return toReturn, fmt.Errorf("failed to get protocol version from %s", dsp.Address)
	}

	toReturn.ProtocolVersion = strings.Trim(protocol, "3000:")

	// get firmware version
	firmware, err := dsp.hardwareCommand(ctx, FirmwareVersion, "")
	if err != nil {
		return toReturn, fmt.Errorf("failed to get firmware version from %s", dsp.Address)
	}

	toReturn.FirmwareVersion = firmware

	// get serial number
	serial, err := dsp.hardwareCommand(ctx, SerialNumber, "")
	if err != nil {
		return toReturn, fmt.Errorf("failed to get serial number from %s", dsp.Address)
	}

	toReturn.SerialNumber = serial

	// get IP address
	ipAddress, err := dsp.hardwareCommand(ctx, IPAddress, "")
	if err != nil {
		return toReturn, fmt.Errorf("failed to get IP address from %s... ironic...", dsp.Address)
	}

	// get gateway
	gateway, err := dsp.hardwareCommand(ctx, Gateway, "")
	if err != nil {
		return toReturn, fmt.Errorf("failed to get the gateway address from %s", dsp.Address)
	}

	// get MAC address
	mac, err := dsp.hardwareCommand(ctx, MACAddress, "")
	if err != nil {
		return toReturn, fmt.Errorf("failed to get the MAC address from %s", dsp.Address)
	}

	// set network information
	toReturn.NetworkInfo = structs.NetworkInfo{
		IPAddress:  ipAddress,
		MACAddress: mac,
		Gateway:    gateway,
	}

	return toReturn, nil
}

func (dsp *KramerVP558) hardwareCommand(ctx context.Context, commandType, param string) (string, error) {
	var cmd []byte

	if len(param) > 0 {
		num, _ := strconv.Atoi(param)
		cmd = []byte(fmt.Sprintf("#%s? %d\r\n", commandType, num))
	} else {
		cmd = []byte(fmt.Sprintf("#%s?\r\n", commandType))
	}

	resp, err := dsp.SendCommand(ctx, cmd)

	if err != nil {
		return "", fmt.Errorf("unable to send command: %w", err)
	}
	resps := string(resp)
	resps = strings.Split(resps, fmt.Sprintf("%s", commandType))[1]
	resps = strings.Trim(resps, "\r\n")
	resps = strings.TrimSpace(resps)

	return resps, nil
}

func (dsp *KramerVP558) GetHardwareInfo(ctx context.Context) (structs.HardwareInfo, error) {
	var toReturn structs.HardwareInfo
	// get the hostname
	addr, e := net.LookupAddr(dsp.Address)
	if e != nil {
		toReturn.Hostname = dsp.Address
	} else {
		toReturn.Hostname = strings.Trim(addr[0], ".")
	}

	// get build date
	buildDate, err := dsp.hardwareCommand(ctx, BuildDate, "")
	if err != nil {
		return toReturn, fmt.Errorf("failed to get build date from %s", dsp.Address)
	}

	toReturn.BuildDate = buildDate

	// get device model
	model, err := dsp.hardwareCommand(ctx, Model, "")
	if err != nil {
		return toReturn, fmt.Errorf("failed to get model number from %s", dsp.Address)
	}

	toReturn.ModelName = model

	// get device protocol version
	protocol, err := dsp.hardwareCommand(ctx, ProtocolVersion, "")
	if err != nil {
		return toReturn, fmt.Errorf("failed to get protocol version from %s", dsp.Address)
	}

	toReturn.ProtocolVersion = strings.Trim(protocol, "3000:")

	// get firmware version
	firmware, err := dsp.hardwareCommand(ctx, FirmwareVersion, "")
	if err != nil {
		return toReturn, fmt.Errorf("failed to get firmware version from %s", dsp.Address)
	}

	toReturn.FirmwareVersion = firmware

	// get serial number
	serial, err := dsp.hardwareCommand(ctx, SerialNumber, "")
	if err != nil {
		return toReturn, fmt.Errorf("failed to get serial number from %s", dsp.Address)
	}

	toReturn.SerialNumber = serial

	// get IP address
	ipAddress, err := dsp.hardwareCommand(ctx, IPAddress, "")
	if err != nil {
		return toReturn, fmt.Errorf("failed to get IP address from %s... ironic...", dsp.Address)
	}

	// get gateway
	gateway, err := dsp.hardwareCommand(ctx, Gateway, "")
	if err != nil {
		return toReturn, fmt.Errorf("failed to get the gateway address from %s", dsp.Address)
	}

	// get MAC address
	mac, err := dsp.hardwareCommand(ctx, MACAddress, "")
	if err != nil {
		return toReturn, fmt.Errorf("failed to get the MAC address from %s", dsp.Address)
	}

	// set network information
	toReturn.NetworkInfo = structs.NetworkInfo{
		IPAddress:  ipAddress,
		MACAddress: mac,
		Gateway:    gateway,
	}

	return toReturn, nil
}
