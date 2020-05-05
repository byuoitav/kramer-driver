package kramer

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	maxVolume   = 15
	minVolume   = -100
	volumeRange = maxVolume - minVolume
)

// GetVolume returns the volume Level for the given input
// The blocks are going to be a number between 1-20, determined by its configuration
func (dsp *KramerAFM20DSP) GetVolumeByBlock(ctx context.Context, block string) (int, error) {

	cmd := []byte(fmt.Sprintf("#X-AUD-LVL? OUT.ANALOG_AUDIO.%s.AUDIO.1\r\n", block))
	resp, err := dsp.SendCommand(ctx, cmd)
	if err != nil {
		logError(err.Error())
		return 0, fmt.Errorf("error sending command: %w", err)
	}

	resps := string(resp)
	if strings.Contains(resps, "ERR") {
		return 0, fmt.Errorf("an error occured: (command: %s) response: %s)", cmd, resps)
	}
	resps = strings.TrimSpace(resps)

	parts := strings.Split(resps, ",")

	dbParts := strings.Split(parts[1], ".")
	currentDB, err := strconv.Atoi(dbParts[0])
	if err != nil {
		return 0, err
	}

	volume := convertBackToVolume(currentDB)

	return volume, nil
}

// SetVolume changes the volume level on the given block to the level parameter
// The blocks are going to be a number between 1-20, determined by its configuration
func (dsp *KramerAFM20DSP) SetVolumeByBlock(ctx context.Context, block string, level int) error {
	volumeLevel := convertToDB(level)
	var cmd []byte
	cmd = []byte(fmt.Sprintf("#X-AUD-LVL OUT.ANALOG_AUDIO.%s.AUDIO.1, %v\r", block, volumeLevel))

	resp, err := dsp.SendCommand(ctx, cmd)
	if err != nil {
		logError(err.Error())
		return fmt.Errorf("error sending command: %w", err)
	}

	resps := string(resp)
	if strings.Contains(resps, "ERR") {
		return fmt.Errorf("an error occured: (command: %s) response: %s)", cmd, resps)
	}

	return nil
}

//converts a volume level 0-100 to the db range between -100 and 15 db
func convertToDB(level int) int {
	volumeToDB := float64(volumeRange) / float64(100)
	convertedValue := int(math.Round(float64(volumeToDB * float64(level))))
	convertedValue = convertedValue + minVolume

	return convertedValue
}

//converts the db level of the device to 0-100 volume range
func convertBackToVolume(level int) int {

	dbToVolume := float64(100) / float64(volumeRange)
	level = level - minVolume
	convertedValue := int(math.Round(float64(float64(level) * dbToVolume)))
	return convertedValue
}

// GetVolume returns the volume Level for the given input
// Audio inputs are formatted 0:0 - 4:2, and audio level is between 0-100.
// for more information on Audio Inputs reference https://cdn.kramerav.com/web/downloads/manuals/vp-558_rev_4.pdf (pg. 64)
func (vsdsp *KramerVP558) GetVolumeByBlock(ctx context.Context, block string) (int, error) {

	cmd := []byte(fmt.Sprintf("#AUD-LVL? 1,%s\r\n", block))
	resp, err := vsdsp.SendCommand(ctx, cmd)
	if err != nil {
		logError(err.Error())
		return 0, fmt.Errorf("error sending command: %w", err)
	}

	resps := string(resp)
	if strings.Contains(resps, "ERR") {
		return 0, fmt.Errorf("an error occured: (command: %s) response: %s)", cmd, resps)
	}
	resps = strings.TrimSpace(resps)
	parts := strings.Split(resps, ",")

	volume, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, err
	}

	return volume, nil
}

// SetVolume changes the volume level on the given block to the level parameter
// Audio inputs are formatted 0:0 - 4:2, and audio level is between 0-100.
// for more information on Audio Inputs reference https://cdn.kramerav.com/web/downloads/manuals/vp-558_rev_4.pdf (pg. 64)
func (vsdsp *KramerVP558) SetVolumeByBlock(ctx context.Context, block string, level int) error {
	var cmd []byte
	cmd = []byte(fmt.Sprintf("#AUD-LVL 1,%s,%v\r", block, level))

	resp, err := vsdsp.SendCommand(ctx, cmd)
	if err != nil {
		logError(err.Error())
		return fmt.Errorf("error sending command: %w", err)
	}

	resps := string(resp)
	if strings.Contains(resps, "ERR") {
		return fmt.Errorf("an error occured: (command: %s) response: %s)", cmd, resps)
	}

	return nil
}
