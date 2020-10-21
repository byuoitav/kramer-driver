package kramer

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/byuoitav/common/status"
	"github.com/fatih/color"
)

// GetInput returns the current input
func (vs *Kramer4x4) GetAudioVideoInputs(ctx context.Context) (map[string]string, error) {
	toReturn := make(map[string]string)

	for x := 0; x < 4; x++ {
		vs.Log.Debugf("Getting input for output port %d", x)
		vs.Log.Debugf("Changing to 1-based indexing... (+1 to each port number)")

		cmd := []byte(fmt.Sprintf("#VID? %d\r", x+1))
		vs.Log.Debugf("Command: %s", cmd)
		resp, err := vs.SendCommand(ctx, cmd)
		if err != nil {
			vs.Log.Errorf("error sending command: %s", err.Error())
			return toReturn, fmt.Errorf("error sending command: %w", err)
		}

		resps := string(resp)
		if !strings.Contains(resps, "VID") {
			return toReturn, fmt.Errorf("Incorrect response for command (%s). (Response: %s)", cmd, resps)
		}

		parts := strings.Split(resps, "VID")
		resps = strings.TrimSpace(parts[1])

		parts = strings.Split(resps, ">")

		var i status.Input
		i.Input = parts[0]

		vs.Log.Debugf("Changing to 0-based indexing... (-1 to each port number)")
		i.Input, err = ToIndexZero(i.Input)
		if err != nil {
			return toReturn, fmt.Errorf("unable to switch to index zero: %w", err)
		}

		i.Input = fmt.Sprintf("%v:%v", i.Input, x)
		color.Set(color.FgGreen, color.Bold)
		vs.Log.Debugf("Input for output port %s is %v", i, i.Input)

		toReturn[strconv.Itoa(x)] = i.Input
	}
	return toReturn, nil
}

// func (vs *Kramer4x4) GetAudioVideoInputs(ctx context.Context) (map[string]string, error) {
// 	inputs := make(map[string]string)
// 	cmd := []byte(fmt.Sprintf("#VID? *\n"))
// 	resp, err := vs.SendCommand(ctx, cmd)
// 	if err != nil {
// 		return inputs, fmt.Errorf("error sending command: %w", err)
// 	}
// 	split := strings.Split(string(resp), "VID")
// 	if len(split) != 2 {
// 		// TODO weird response
// 		return nil, fmt.Errorf("WEIRD response: %v", split)
// 	}
// 	for _, str := range strings.Split(split[1], ",") {
// 		split := strings.Split(strings.TrimSpace(str), ">")
// 		if len(split) != 2 {
// 			// TODO weird response
// 			return nil, fmt.Errorf("WERID Response: %v", split)
// 		}
// 		inputs[split[0]] = split[1]
// 	}
// 	return inputs, nil
// }

// SwitchInput changes the input on the given output to input
func (vs *Kramer4x4) SetAudioVideoInput(ctx context.Context, output, input string) error {
	i, err := ToIndexOne(input)
	if err != nil || LessThanZero(input) {
		return fmt.Errorf("error! Input parameter %s is not valid", input)
	}

	o, err := ToIndexOne(output)
	if err != nil || LessThanZero(output) {
		return fmt.Errorf("error! Output parameter must be zero or greater")
	}

	vs.Log.Debugf("Routing %v to %v on %v", input, output, vs.Address)
	vs.Log.Debugf("Changing to 1-based indexing... (+1 to each port number)")

	cmd := []byte(fmt.Sprintf("#VID %s>%s\r\n", i, o))

	resp, err := vs.SendCommand(ctx, cmd)
	if err != nil {
		vs.Log.Errorf("unable to send command: %s", err.Error())
		return fmt.Errorf("unable to send command: %w", err)
	}

	resps := string(resp)
	if !strings.Contains(resps, "VID") {
		return fmt.Errorf("Incorrect response for command (%s). (Response: %s)", cmd, resp)
	}

	return nil
}

// This function converts a number (in a string) to index-base 0.
func ToIndexZero(numString string) (string, error) {
	num, err := strconv.Atoi(numString)
	if err != nil {
		return "", err
	}

	num--

	return strconv.Itoa(num), nil
}

// GetInput returns the current input
// The API is zero indexed, so outputs 0-3 correspond with outputs 1-4 on device
// inputs 0-10 correspond with outputs 1-11 see https://cdn.kramerav.com/web/downloads/manuals/vp-558_rev_4.pdf (page 66)
func (vsdsp *KramerVP558) GetAudioVideoInputs(ctx context.Context) (map[string]string, error) {
	toReturn := make(map[string]string)

	for x := 0; x < 4; x++ {
		vsdsp.Log.Debugf("Getting input for output port %s", x)

		cmd := []byte(fmt.Sprintf("#ROUTE? 1,%d\r\n", x))
		resp, err := vsdsp.SendCommand(ctx, cmd, false)
		if err != nil {
			vsdsp.Log.Errorf("error sending command: %s", err.Error())
			return toReturn, fmt.Errorf("error sending command: %w", err)
		}

		resps := string(resp)
		if strings.Contains(resps, "ERR") {
			return toReturn, fmt.Errorf("Incorrect response for command (%s). (Response: %s)", cmd, resps)
		}
		resps = strings.TrimSpace(resps)

		parts := strings.Split(resps, ",")
		if len(parts) != 3 {
			return toReturn, fmt.Errorf("Incorrect response for command (%s). (Response: %s)", cmd, resps)
		}

		var i status.Input
		i.Input = parts[2]

		// vsdsp.Log.Infof("successfully got input", zap.String("output", output), zap.String("input", i.Input))
		toReturn[strconv.Itoa(x)] = i.Input
	}
	return toReturn, nil
}

// SwitchInput changes the input on the given output to input
// outputs 1-4 on device inputs 1-11 see https://cdn.kramerav.com/web/downloads/manuals/vp-558_rev_4.pdf (page 66)
func (vsdsp *KramerVP558) SetAudioVideoInput(ctx context.Context, output, input string) error {

	vsdsp.Log.Debugf("Routing %v to %v on %v", input, output, vsdsp.Address)
	// vsdsp.Log.Infof("sending setInput command", zap.String("output", output), zap.String("input", input))

	cmd := []byte(fmt.Sprintf("#ROUTE 1,%s,%s\r\n", output, input))

	//cheack to see if the current input is going to be changing
	currentInputs, err := vsdsp.GetAudioVideoInputs(ctx)
	if err != nil {
		vsdsp.Log.Errorf("error sending command: %s", err.Error())
		return fmt.Errorf("error sending command: %w", err)
	}
	//if there is a change, two responses will be sent and both need to be read
	readAgain := false
	if currentInputs[output] != input {
		readAgain = true
	}

	resp, err := vsdsp.SendCommand(ctx, cmd, readAgain)
	if err != nil {
		vsdsp.Log.Errorf("unable to send command: %s", err.Error())
		return fmt.Errorf("unable to send command: %w", err)
	}

	resps := string(resp)
	if strings.Contains(resps, "ERR") {
		return fmt.Errorf("Incorrect response for command (%s). (Response: %s)", cmd, resp)
	}

	// vsdsp.Log.Infof("successfully sent setInput command", zap.String("output", output), zap.String("input", input))

	return nil
}
