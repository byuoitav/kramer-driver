package via

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/byuoitav/common/log"
)

type VIA struct {
	Address  string
	Username string
	Password string
}

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

func (v *VIA) GetVolumeByBlock(ctx context.Context, block string) (int, error) {
	// opening a connection with the VIAs and then return the volume for the device
	log.L.Infof("Getting volume for %v", v.Address)
	viaResp, err := v.GetVolume(ctx, v.Address)
	if err != nil {
		return 0, err
	}
	log.L.Infof("%v", viaResp)

	return viaResp, nil

}

func (v *VIA) GetMutedByBlock(ctx context.Context) (bool, error) {
	// Return error because VIAs do not support a mute function
	return errors.New(fmt.Sprintf("Error in getting mute status of VIA, Feature not supported"))
}

func (v *VIA) SetVolumeByBlock(ctx context.Context, block string, volume int) error {
	// Connect and set the volume on the VIA
	log.L.Infof("Setting volume for %v", v.Address)
	viaResp, err := v.SetVolume(ctx, v.Address, volume)
	if err != nil {
		log.L.Debugf("Failed to set VIA Volume for %v", v.Address)
		return errors.New(fmt.Sprintf("Error setting volume for %v", v.Address))
	}
	return nil
}

func (v *VIA) SetMutedByBlock(ctx context.Context, block string, muted bool) error {
	// Return error because VIAs do not support mute
	return errors.New(fmt.Sprintf("Error setting mute status of VIA, Feature not supported"))
}

func (v *VIA) RebootVIA(ctx context.Context) {
	// Reboot a VIA using the API
	log.L.Infof("Rebooting %v", v.Address)
	viaResp, err := v.Reboot(ctx)
	if err != nil {
		log.L.Debugf("Failed to reboot the VIA: %v", err)
		return errors.New(fmt.Sprintf("Error Rebooting VIA: %v", err))
	}
}

func (v *VIA) ResetVIA(ctx context.Context) {
	// Reset a VIA sessions - Causes VIAAdmin to log out and log back in
	log.L.Infof("Reseting %v", v.Address)
	viaResp, err := v.Reset(ctx)
	if err != nil {
		log.L.Debugf("Failed to reboot the VIA due to %v", err)
		return errors.New(fmt.Sprintf("Error Resetting VIA: %v", err))
	}
}
