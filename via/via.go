package via

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/byuoitav/common/log"
)

// VIA: Struct that defines general parameters needed for any VIA
type VIA struct {
	Address  string
	Username string
	Password string
}

// Command: Struct used to build the commands that need to be sent to the VIA
type Command struct {
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

type message struct {
	EventType string
	Action    string
	User      string
	State     string
}

// Device - a representation of a device involved in a TEC Pi system.
type Device struct {
	ID          string                 `json:"_id"`
	Name        string                 `json:"name"`
	Address     string                 `json:"address"`
	Description string                 `json:"description"`
	DisplayName string                 `json:"display_name"`
	Type        DeviceType             `json:"type,omitempty"`
	Roles       []Role                 `json:"roles"`
	Ports       []Port                 `json:"ports"`
	Tags        []string               `json:"tags,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`

	// Proxy is a map of regex (matching command id's) to the host:port of the proxy
	Proxy map[string]string `json:"proxy,omitempty"`
}

// GetVolumeByBlock: opening a connection with the VIAs and then return the volume for the device
func (v *VIA) GetVolumeByBlock(ctx context.Context, block string) (int, error) {
	log.L.Infof("Getting volume for %v", v.Address)
	viaResp, err := v.GetVolume(ctx, v.Address)
	if err != nil {
		return 0, err
	}
	log.L.Infof("%v", viaResp)

	return viaResp, nil

}

// GetMutedByBlock: Return error because VIAs do not support a mute function
func (v *VIA) GetMutedByBlock(ctx context.Context) (bool, error) {
	return _, errors.New(fmt.Sprintf("Error in getting mute status of VIA, Feature not supported"))
}

// SetVolumeByBlock: Connect and set the volume on the VIA
func (v *VIA) SetVolumeByBlock(ctx context.Context, block string, volume int) error {
	log.L.Infof("Setting volume for %v", v.Address)
	viaResp, err := v.SetVolume(ctx, v.Address, volume)
	if err != nil {
		log.L.Debugf("Failed to set VIA Volume for %v", v.Address)
		return errors.New(fmt.Sprintf("Error setting volume for %v", v.Address))
	}
	return nil
}

// SetMutedByBlock: Return error because VIAs do not support mute
func (v *VIA) SetMutedByBlock(ctx context.Context, block string, muted bool) error {
	return errors.New(fmt.Sprintf("Error setting mute status of VIA, Feature not supported"))
}

// RebootVIA: Reboot a VIA using the API
func (v *VIA) RebootVIA(ctx context.Context) {
	log.L.Infof("Rebooting %v", v.Address)
	viaResp, err := v.Reboot(ctx)
	if err != nil {
		log.L.Debugf("Failed to reboot the VIA: %v", err)
		return errors.New(fmt.Sprintf("Error Rebooting VIA: %v", err))
	}
}

// ResetVIA: Reset a VIA sessions - Causes VIAAdmin to log out and log back in
func (v *VIA) ResetVIA(ctx context.Context) {
	log.L.Infof("Reseting %v", v.Address)
	viaResp, err := v.Reset(ctx)
	if err != nil {
		log.L.Debugf("Failed to reboot the VIA due to %v", err)
		return errors.New(fmt.Sprintf("Error Resetting VIA: %v", err))
	}
}
