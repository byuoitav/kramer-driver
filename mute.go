package kramer

import (
	"context"
	"fmt"
	"strings"
)

// GetMuted returns the Mute Status current input
// The blocks are going to be a number between 1-20, determined by its configuration
func (dsp *KramerAFM20DSP) GetMutedByBlock(ctx context.Context, block string) (bool, error) {

	fmt.Println(block)
	cmd := []byte(fmt.Sprintf("#X-MUTE? OUT.ANALOG_AUDIO.%s.AUDIO.1\r\n", block))
	resp, err := dsp.SendCommand(ctx, cmd)
	if err != nil {
		logError(err.Error())
		return false, fmt.Errorf("error sending command: %w", err)
	}

	resps := string(resp)
	if strings.Contains(resps, "ERR") {
		return false, fmt.Errorf("an error occured: (command: %s) response: %s)", cmd, resps)
	}

	parts := strings.Split(resps, ",")
	resps = strings.TrimSpace(parts[1])

	if resps == "OFF" {
		return false, nil
	} else {
		return true, nil
	}
}

// SetMuted changes the input on the given output to input
// The blocks are going to be a number between 1-20, determined by its configuration
func (dsp *KramerAFM20DSP) SetMutedByBlock(ctx context.Context, block string, muted bool) error {

	var cmd []byte
	if muted {
		cmd = []byte(fmt.Sprintf("#X-MUTE OUT.ANALOG_AUDIO.%s.AUDIO.1, ON\r", block))
	} else {
		cmd = []byte(fmt.Sprintf("#X-MUTE OUT.ANALOG_AUDIO.%s.AUDIO.1, OFF\r", block))
	}
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

// GetMuted returns the Mute Status current input
// Audio inputs are formatted 0:0 - 4:2, and audio level is between 0-100.
// for more information on Audio Inputs reference https://cdn.kramerav.com/web/downloads/manuals/vp-558_rev_4.pdf (pg. 64)
func (vsdsp *KramerVP558) GetMutedByBlock(ctx context.Context, block string) (bool, error) {

	fmt.Println(block)
	cmd := []byte(fmt.Sprintf("#MUTE? %s\r\n", block))
	resp, err := vsdsp.SendCommand(ctx, cmd)
	if err != nil {
		logError(err.Error())
		return false, fmt.Errorf("error sending command: %w", err)
	}

	resps := string(resp)
	if strings.Contains(resps, "ERR") {
		return false, fmt.Errorf("an error occured: (command: %s) response: %s)", cmd, resps)
	}
	parts := strings.Split(resps, ",")
	resps = strings.TrimSpace(parts[1])

	if resps == "0" {
		return false, nil
	} else {
		return true, nil
	}
}

// setMuted changes the input on the given output to input
// Audio inputs are formatted 0:0 - 4:2, and audio level is between 0-100.
// for more information on Audio Inputs reference https://cdn.kramerav.com/web/downloads/manuals/vp-558_rev_4.pdf (pg. 64)
func (vsdsp *KramerVP558) SetMutedByBlock(ctx context.Context, block string, muted bool) error {

	var cmd []byte
	if muted {
		cmd = []byte(fmt.Sprintf("#MUTE %s,1\r", block))
	} else {
		cmd = []byte(fmt.Sprintf("#MUTE %s,0\r", block))
	}
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
