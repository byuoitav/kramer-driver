package via

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/byuoitav/common/log"
	//"github.com/byuoitav/common/structs"
)

//VIA Struct that defines general parameters needed for any VIA
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
func (v *VIA) GetMutedByBlock(ctx context.Context, block string) (bool, error) {
	return false, errors.New(fmt.Sprintf("Error in getting mute status of VIA, Feature not supported"))
}

// SetVolumeByBlock: Connect and set the volume on the VIA
func (v *VIA) SetVolumeByBlock(ctx context.Context, block string, volume int) error {
	log.L.Infof("Setting volume for %v", v.Address)

	volumec := strconv.Itoa(volume)

	_, err := v.SetVolume(ctx, v.Address, volumec)
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

// GetInfo: needed by the DSP drivers implementation.  Will get hardware information
func (v *VIA) GetInfo(ctx context.Context) (interface{}, error) {
	var info interface{}
	return info, fmt.Errorf("GetInfo has not been implemented in this version of the driver")
}

// RebootVIA: Reboot a VIA using the API
func (v *VIA) RebootVIA(ctx context.Context) error {
	log.L.Infof("Rebooting %v", v.Address)
	err := v.Reboot(ctx)
	if err != nil {
		log.L.Debugf("Failed to reboot the VIA: %v", err)
		return errors.New(fmt.Sprintf("Error Rebooting VIA: %v", err))
	}
	return nil
}

// ResetVIA: Reset a VIA sessions - Causes VIAAdmin to log out and log back in
func (v *VIA) ResetVIA(ctx context.Context) error {
	log.L.Infof("Reseting %v", v.Address)
	err := v.Reset(ctx)
	if err != nil {
		log.L.Debugf("Failed to reboot the VIA due to %v", err)
		return errors.New(fmt.Sprintf("Error Resetting VIA: %v", err))
	}
	return nil
}

func (v *VIA) VIARoomCode(ctx context.Context, address string) error {
	log.L.Infof("Getting Room Code for %v", v.Address)
	resp, err := v.GetRoomCode(ctx, address)
	if err != nil {
		log.L.Debugf("Failed to get room code: %v", err)
		return errors.New(fmt.Sprintf("Error getting room code: %v", err))
	}
	log.L.Infof("VIA room code: %v", resp)
	return nil
}

// VolumeParse parser to pull out the volume level from the VIA API returned string
func VolumeParse(vollevel string) (int, error) {
	re := regexp.MustCompile("[0-9]+")
	vol := re.FindString(vollevel)
	vfin, err := strconv.Atoi(vol)
	if err != nil {
		err = fmt.Errorf("Error converting response: %s", err.Error())
		log.L.Infof("%s", err.Error())
		return 0, err
	}
	return vfin, nil
}
