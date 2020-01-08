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

// GetInput returns the current input
func (vs *VideoSwitcher) getInputByOutput(ctx context.Context, output string) (string, error) {

	p, err := ToIndexOne(output)
	if err != nil || LessThanZero(output) {
		return "", fmt.Errorf("Error! Port parameter must be zero or greater")
	}

	log.L.Debugf("Getting input for output port %s", output)
	log.L.Debugf("Changing to 1-based indexing... (+1 to each port number)")

	cmd := []byte(fmt.Sprintf("#VID? %s", p))
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
func (vs *VideoSwitcher) setInputByOutput(ctx context.Context, output, input string) error {
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

	cmd := []byte(fmt.Sprintf("#VID %s>%s", i, o))

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
