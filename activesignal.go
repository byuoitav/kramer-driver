package kramer

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/byuoitav/common/structs"
	"github.com/fatih/color"
)

func (vs *VideoSwitcher) GetActiveSignal(ctx context.Context, port string) (error, structs.ActiveSignal) {
	rW := true
	var signal structs.ActiveSignal
	i, err := ToIndexOne(port)
	if err != nil || LessThanZero(port) {
		return fmt.Errorf("Error: %s", err), signal
	}

	signal, ne := GetActiveSignalByPort(vs.Address, i, rW)
	if ne != nil {
		return fmt.Errorf("Error: %s", err), signal
	}

	return nil, signal
}

// This function converts a number (in a string) to index-based 1.
func ToIndexOne(numString string) (string, error) {
	num, err := strconv.Atoi(numString)
	if err != nil {
		return "", err
	}

	// add one to make it match pulse eight.
	// we are going to use 0 based indexing on video matrixing,
	// and the kramer uses 1-based indexing.
	num++

	return strconv.Itoa(num), nil
}

// Returns if a given number (in a string) is less than zero.
func LessThanZero(numString string) bool {
	defer color.Unset()
	num, err := strconv.Atoi(numString)
	if err != nil {
		color.Set(color.FgRed)
		return false
	}

	return num < 0
}

func GetActiveSignalByPort(address, port string, readWelcome bool) (structs.ActiveSignal, error) {
	var signal structs.ActiveSignal

	signal.Active = false

	signalResponse, err := hardwareCommand(Signal, port, address, readWelcome)
	if err != nil {
		return signal, fmt.Errorf("failed to get the signal for %s on %s", port, address)
	}

	signalStatus := strings.Split(signalResponse, ",")[1]

	if signalStatus == "1" {
		signal.Active = true
	}

	return signal, nil
}
