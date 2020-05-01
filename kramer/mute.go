package kramer

import (
	"context"
	"fmt"
	"strings"
)

// GetMuted returns the Mute Status current input
func (dsp *Dsp) GetMutedByBlock(ctx context.Context, block string) (bool, error) {

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
func (dsp *Dsp) SetMutedByBlock(ctx context.Context, block string, muted bool) error {

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
func (vsdsp *VideoSwitcherDsp) GetMutedByBlock(ctx context.Context, block string) (bool, error) {

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
func (vsdsp *VideoSwitcherDsp) SetMutedByBlock(ctx context.Context, block string, muted bool) error {

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
