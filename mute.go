package kramer

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

// GetMuted returns the Mute Status current input
// The blocks are going to be a number between 1-20, determined by its configuration
func (dsp *KramerAFM20DSP) Mutes(ctx context.Context, blocks []string) (map[string]bool, error) {

	toReturn := make(map[string]bool)

	for _, block := range blocks {
		dsp.Log.Infof("sending get muteStatus command", zap.String("block", block))
		cmd := []byte(fmt.Sprintf("#X-MUTE? OUT.ANALOG_AUDIO.%s.AUDIO.1\r\n", block))
		resp, err := dsp.SendCommand(ctx, cmd)
		if err != nil {
			dsp.Log.Errorf("error sending command: %s", err.Error())
			return toReturn, fmt.Errorf("error sending command: %w", err)
		}

		resps := string(resp)
		if strings.Contains(resps, "ERR") {
			return toReturn, fmt.Errorf("an error occured: (command: %s) response: %s)", cmd, resps)
		}
		resps = strings.TrimSpace(resps)

		parts := strings.Split(resps, ",")

		if parts[1] == "OFF" {
			dsp.Log.Infof("successfully got mute status", zap.String("block", block), zap.Bool("status", false))
			toReturn[block] = false
		} else {
			dsp.Log.Infof("successfully got mute status", zap.String("block", block), zap.Bool("status", true))
			toReturn[block] = true
		}
	}

	return toReturn, nil
}

// SetMuted changes the input on the given output to input
// The blocks are going to be a number between 1-20, determined by its configuration
func (dsp *KramerAFM20DSP) SetMute(ctx context.Context, block string, mute bool) error {
	dsp.Log.Infof("sending set muteStatus command", zap.String("block", block), zap.Bool("status", mute))

	var cmd []byte
	if mute {
		cmd = []byte(fmt.Sprintf("#X-MUTE OUT.ANALOG_AUDIO.%s.AUDIO.1, ON\r", block))
	} else {
		cmd = []byte(fmt.Sprintf("#X-MUTE OUT.ANALOG_AUDIO.%s.AUDIO.1, OFF\r", block))
	}
	resp, err := dsp.SendCommand(ctx, cmd)
	if err != nil {
		dsp.Log.Errorf("error sending command: %s", err.Error())
		return fmt.Errorf("error sending command: %w", err)
	}

	resps := string(resp)
	if strings.Contains(resps, "ERR") {
		return fmt.Errorf("an error occured: (command: %s) response: %s)", cmd, resps)
	}

	dsp.Log.Infof("successfully set mute status", zap.String("block", block), zap.Bool("status", mute))

	return nil
}

// GetMuted returns the Mute Status current input
// Audio inputs are formatted 0:0 - 4:2, and audio level is between 0-100.
// for more information on Audio Inputs reference https://cdn.kramerav.com/web/downloads/manuals/vp-558_rev_4.pdf (pg. 64)
func (vsdsp *KramerVP558) Mutes(ctx context.Context, blocks []string) (map[string]bool, error) {
	toReturn := make(map[string]bool)

	for _, block := range blocks {
		vsdsp.Log.Infof("sending get mute status command", zap.String("block", block))
		cmd := []byte(fmt.Sprintf("#MUTE? %s\r\n", block))
		resp, err := vsdsp.SendCommand(ctx, cmd, false)
		if err != nil {
			vsdsp.Log.Errorf("error sending command: %s", err.Error())
			return toReturn, fmt.Errorf("error sending command: %w", err)
		}

		resps := string(resp)
		if strings.Contains(resps, "ERR") {
			return toReturn, fmt.Errorf("an error occured: (command: %s) response: %s)", cmd, resps)
		}
		resps = strings.TrimSpace(resps)

		parts := strings.Split(resps, ",")

		if len(parts) < 2 {
			return toReturn, fmt.Errorf("unexpected response, unable to parse: %s", resps)
		}

		if parts[1] == "0" {
			vsdsp.Log.Infof("successfully got mute status", zap.String("block", block), zap.Bool("status", false))
			toReturn[block] = false
		} else {
			vsdsp.Log.Infof("successfully got mute status", zap.String("block", block), zap.Bool("status", true))
			toReturn[block] = true
		}
	}

	return toReturn, nil
}

// setMuted changes the input on the given output to input
// Audio inputs are formatted 0:0 - 4:2, and audio level is between 0-100.
// for more information on Audio Inputs reference https://cdn.kramerav.com/web/downloads/manuals/vp-558_rev_4.pdf (pg. 64)
func (vsdsp *KramerVP558) SetMute(ctx context.Context, block string, muted bool) error {
	vsdsp.Log.Infof("sending set muteStatus command", zap.String("block", block), zap.Bool("status", muted))

	//cheack to see if the mute status is going to be changing
	currentStatus, err := vsdsp.Mutes(ctx, []string{block})
	if err != nil {
		vsdsp.Log.Errorf("error sending command: %s", err.Error())
		return fmt.Errorf("error sending command: %w", err)
	}
	//if there is a change, two responses will be sent and both need to be read
	readAgain := false
	if currentStatus[block] != muted {
		readAgain = true
	}

	var cmd []byte
	if muted {
		cmd = []byte(fmt.Sprintf("#MUTE %s,1\r", block))
	} else {
		cmd = []byte(fmt.Sprintf("#MUTE %s,0\r", block))
	}
	resp, err := vsdsp.SendCommand(ctx, cmd, readAgain)
	if err != nil {
		vsdsp.Log.Errorf("error sending command: %s", err.Error())
		return fmt.Errorf("error sending command: %w", err)
	}

	resps := string(resp)
	if strings.Contains(resps, "ERR") {
		return fmt.Errorf("an error occured: (command: %s) response: %s)", cmd, resps)
	}

	vsdsp.Log.Infof("successfully set mute status", zap.String("block", block), zap.Bool("status", muted))

	return nil
}
