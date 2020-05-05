package kramer

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/status"
	"github.com/fatih/color"
)

func logError(e string) {
	color.Set(color.FgRed)
	log.L.Infof("%s", e)
	color.Unset()
}

// GetInput returns the current input
func (vs *Kramer4x4) GetInputByOutput(ctx context.Context, output string) (string, error) {

	p, err := ToIndexOne(output)
	if err != nil || LessThanZero(output) {
		return "", fmt.Errorf("Error! Port parameter must be zero or greater")
	}

	log.L.Debugf("Getting input for output port %s", output)
	log.L.Debugf("Changing to 1-based indexing... (+1 to each port number)")

	cmd := []byte(fmt.Sprintf("#VID? %s\r\n", p))
	resp, err := vs.SendCommand(ctx, cmd)
	if err != nil {
		logError(err.Error())
		return "", fmt.Errorf("error sending command: %w", err)
	}

	resps := string(resp)
	if !strings.Contains(resps, "VID") {
		return "", fmt.Errorf("Incorrect response for command (%s). (Response: %s)", cmd, resps)
	}

	parts := strings.Split(resps, "VID")
	resps = strings.TrimSpace(parts[1])

	parts = strings.Split(resps, ">")

	var i status.Input
	i.Input = parts[0]

	log.L.Debugf("Changing to 0-based indexing... (-1 to each port number)")
	i.Input, err = ToIndexZero(i.Input)
	if err != nil {
		return "", fmt.Errorf("unable to switch to index zero: %w", err)
	}

	i.Input = fmt.Sprintf("%v:%v", i.Input, output)
	color.Set(color.FgGreen, color.Bold)
	log.L.Debugf("Input for output port %s is %v", output, i.Input)
	return i.Input, nil
}

// SwitchInput changes the input on the given output to input
func (vs *Kramer4x4) SetInputByOutput(ctx context.Context, output, input string) error {
	i, err := ToIndexOne(input)
	if err != nil || LessThanZero(input) {
		return fmt.Errorf("Error! Input parameter %s is not valid!", input)
	}

	o, err := ToIndexOne(output)
	if err != nil || LessThanZero(output) {
		return fmt.Errorf("Error! Output parameter must be zero or greater")
	}

	log.L.Debugf("Routing %v to %v on %v", input, output, vs.Address)
	log.L.Debugf("Changing to 1-based indexing... (+1 to each port number)")

	cmd := []byte(fmt.Sprintf("#VID %s>%s\r\n", i, o))

	resp, err := vs.SendCommand(ctx, cmd)
	if err != nil {
		logError(err.Error())
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
func (vsdsp *KramerVP558) GetInputByOutput(ctx context.Context, output string) (string, error) {

	p, err := ToIndexOne(output)
	if err != nil || LessThanZero(output) {
		return "", fmt.Errorf("Error! Port parameter must be zero or greater")
	}

	log.L.Debugf("Getting input for output port %s", output)
	log.L.Debugf("Changing to 1-based indexing... (+1 to each port number)")

	cmd := []byte(fmt.Sprintf("#ROUTE? 1,%s\r\n", p))
	resp, err := vsdsp.SendCommand(ctx, cmd)
	if err != nil {
		logError(err.Error())
		return "", fmt.Errorf("error sending command: %w", err)
	}

	resps := string(resp)
	if strings.Contains(resps, "ERR") {
		return "", fmt.Errorf("Incorrect response for command (%s). (Response: %s)", cmd, resps)
	}
	resps = strings.TrimSpace(resps)

	parts := strings.Split(resps, ",")
	if len(parts) != 3 {
		return "", fmt.Errorf("Incorrect response for command (%s). (Response: %s)", cmd, resps)
	}

	var i status.Input
	i.Input = parts[2]

	log.L.Debugf("Changing to 0-based indexing... (-1 to each port number)")
	i.Input, err = ToIndexZero(i.Input)
	if err != nil {
		return "", fmt.Errorf("unable to switch to index zero: %w", err)
	}

	log.L.Debugf("Input for output port %s is %v", output, i.Input)
	return i.Input, nil
}

// SwitchInput changes the input on the given output to input
// The API is zero indexed, so outputs 0-3 correspond with outputs 1-4 on device
// inputs 0-10 correspond with outputs 1-11 see https://cdn.kramerav.com/web/downloads/manuals/vp-558_rev_4.pdf (page 66)
func (vsdsp *KramerVP558) SetInputByOutput(ctx context.Context, output, input string) error {
	i, err := ToIndexOne(input)
	if err != nil || LessThanZero(input) {
		return fmt.Errorf("Error! Input parameter %s is not valid!", input)
	}

	o, err := ToIndexOne(output)
	if err != nil || LessThanZero(output) {
		return fmt.Errorf("Error! Output parameter must be zero or greater")
	}

	log.L.Debugf("Routing %v to %v on %v", input, output, vsdsp.Address)
	log.L.Debugf("Changing to 1-based indexing... (+1 to each port number)")

	cmd := []byte(fmt.Sprintf("#ROUTE 1,%s,%s\r\n", o, i))

	resp, err := vsdsp.SendCommand(ctx, cmd)
	if err != nil {
		logError(err.Error())
		return fmt.Errorf("unable to send command: %w", err)
	}

	resps := string(resp)
	if strings.Contains(resps, "ERR") {
		return fmt.Errorf("Incorrect response for command (%s). (Response: %s)", cmd, resp)
	}

	return nil
}
