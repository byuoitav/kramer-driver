package kramer

import (
	"encoding/xml"
	"time"
)

// Options Struct is used for logging options
type options struct {
	ttl   time.Duration
	delay time.Duration
}

func (f optionFunc) apply(o *options) {
	f(o)
}

type Option interface {
	apply(*options)
}

// command: Struct used to build the XML commands that need to be sent to the VIA
type command struct {
	XMLName  xml.Name `xml:"P"`
	Username string   `xml:"UN"`
	Password string   `xml:"Pwd"`
	Command  string   `xml:"Cmd"`
	Param1   string   `xml:"P1,omitempty"`
	Param2   string   `xml:"P2,omitempty"`
	Param3   string   `xml:"P3,omitempty"`
	Param4   string   `xml:"P4,omitempty"`
	Param5   string   `xml:"P5,omitempty"`
	Param6   string   `xml:"P6,omitempty"`
	Param7   string   `xml:"P7,omitempty"`
	Param8   string   `xml:"P8,omitempty"`
	Param9   string   `xml:"P9,omitempty"`
	Param10  string   `xml:"P10,omitempty"`
}

// HardwareInfo contains the common information for device hardware information
type HardwareInfo struct {
	Hostname              string           `json:"hostname,omitempty"`
	ModelName             string           `json:"model_name,omitempty"`
	SerialNumber          string           `json:"serial_number,omitempty"`
	BuildDate             string           `json:"build_date,omitempty"`
	FirmwareVersion       string           `json:"firmware_version,omitempty"`
	ProtocolVersion       string           `json:"protocol_version,omitempty"`
	NetworkInfo           NetworkInfo      `json:"network_information,omitempty"`
	FilterStatus          string           `json:"filter_status,omitempty"`
	WarningStatus         []string         `json:"warning_status,omitempty"`
	ErrorStatus           []string         `json:"error_status,omitempty"`
	PowerStatus           string           `json:"power_status,omitempty"`
	PowerSavingModeStatus string           `json:"power_saving_mode_status,omitempty"`
	TimerInfo             []map[string]int `json:"timer_info,omitempty"`
	Temperature           string           `json:"temperature,omitempty"`
}

// NetworkInfo contains the network information for the device
type NetworkInfo struct {
	IPAddress  string   `json:"ip_address,omitempty"`
	MACAddress string   `json:"mac_address,omitempty"`
	Gateway    string   `json:"gateway,omitempty"`
	DNS        []string `json:"dns,omitempty"`
}

// VIAUsers contains the counts of the users logged in to the VIA and their status
type VIAUsers struct {
	InactiveUsers []string `json:"inactive_users"`
	ActiveUsers   []string `json:"active_users"`
	UsersWaiting  []string `json:"users_waiting"`
}

type ActiveSignal struct {
	// TODO should probably change this to a bool
	Active bool `json:"active"`
}
