package kramer

import (
	"context"
	"fmt"

	"github.com/byuoitav/common/structs"
)

type DeviceType int

const (
	Undefined DeviceType = iota
	kramer
	kramer2000
)

type Kramer struct {
	Address     string
	Type        DeviceType
	Room_System string
	System_ID   string
}

func (k *Kramer) SetInputByOutput(ctx context.Context, output, input string) error {
	switch k.Type {
	case kramer:
		return k.setInputByOutput(ctx, output, input)
	case kramer2000:
		return k.setInputByOutput2000(ctx, output, input)
	default:
		return fmt.Errorf("unknown device type")
	}
}

func (k *Kramer) SetFrontLock(ctx context.Context) error {
	return nil
}

func (k *Kramer) GetInputByOutput(ctx context.Context, output, input string) (string, error) {
	switch k.Type {
	case kramer:
		return k.getInputByOutput(ctx, output)
	case kramer2000:
		return k.getInputByOutput2000(ctx, output)

	default:
		return "", fmt.Errorf("unknown device type")
	}
}

func (k *Kramer) GetHardwareInfo(ctx context.Context) (structs.HardwareInfo, error) {
	var resp structs.HardwareInfo
	switch k.Type {
	case kramer:
		return k.getSwitcherHardwareInfo(ctx)
	case kramer2000:
		return k.getSwitcherHardwareInfo(ctx)
	default:
		return resp, fmt.Errorf("unknown device type")
	}
}

func (k *Kramer) GetActiveSignal(ctx context.Context) error {
	return nil
}
