package kramer

import (
	"context"
	"fmt"
	"strings"

	"github.com/byuoitav/common/log"
)

func (vs *Kramer4x4) SetFrontLock(ctx context.Context, state bool) error {
	log.L.Debugf("Setting front button lock status to %v", state)

	var num int8
	if state {
		num = 1
	}

	cmd := []byte(fmt.Sprintf("#LOCK-FP %v\r\n", num))

	resp, err := vs.SendCommand(ctx, cmd)
	switch {
	case err != nil:
		return fmt.Errorf("unable to send command: %w", err)
	case !strings.Contains(string(resp), "OK"):
		return fmt.Errorf("Incorrect response for command (%s). (Response: %s)", cmd, resp)
	}

	return nil
}
