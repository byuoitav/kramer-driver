package kramer

import (
	"context"
	"fmt"

	"github.com/byuoitav/common/structs"
)

// GetInput returns the current input
func (k *Kramer) getInputByOutput2000(ctx context.Context, output string) (string, error) {

	return "", nil
}

// GetHardwareInfo returns a hardware info struct
func (k *Kramer) getHardwareInfo2000(ctx context.Context) (structs.HardwareInfo, error) {
	var hwinfo structs.HardwareInfo

	return hwinfo, nil
}

// SwitchInput changes the input on the given output to input
func (k *Kramer) setInputByOutput2000(ctx context.Context, output, input string) error {
	return nil
}

//GetInfo .
func (k *Kramer) getInfo2000(ctx context.Context) (interface{}, error) {
	var info interface{}
	return info, fmt.Errorf("not currently implemented")
}
