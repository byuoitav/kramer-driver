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
	fmt.Println(resps)
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

func convertToDB(level int) int {
	volumeToDB := float64(volumeRange) / float64(100)
	convertedValue := int(math.Round(float64(volumeToDB * float64(level))))
	convertedValue = convertedValue + minVolume

	return convertedValue
}

func convertBackToVolume(level int) int {

	dbToVolume := float64(100) / float64(volumeRange)
	level = level - minVolume
	convertedValue := int(math.Round(float64(float64(level) * dbToVolume)))
	return convertedValue
}

// GetVolume returns the volume Level for the given input
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
	fmt.Println(resps)
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
