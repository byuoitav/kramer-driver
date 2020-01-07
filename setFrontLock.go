package kramer

import (
	"context"
	"fmt"
	"strings"

	"github.com/byuoitav/common/log"

	"github.com/fatih/color"
)

func (vs *VideoSwitcher) SetFrontLock(ctx context.Context, state bool) error {
	defer color.Unset()
	readWelcome := true

	color.Set(color.FgYellow)
	log.L.Debugf("Setting front button lock status to %v", state)

	var num int8
	if state {
		num = 1
	}
	command := fmt.Sprintf("#LOCK-FP %v", num)

	vs.pool.Do(func(ctx, conn connpool.Conn) error {
		resp, err := SendCommand(conn, command.Address, command.Command)
	})

	if err != nil {
		return err
	}

	if strings.Contains(resp, "OK") {
		return nil
	}

	err = fmt.Errorf("Incorrect response for command (%s). (Response: %s)", command, resp)
	if err != nil {
		return fmt.Errorf("error setting front lock: %s", err)
	}

	color.Set(color.FgGreen, color.Bold)
	log.L.Debugf("Success")
	return nil
}

func SetFrontLockHelper(address string, state, readWelcome bool) error {
	var num int8
	if state {
		num = 1
	}
	command := fmt.Sprintf("#LOCK-FP %v", num)

	resp, err := SendCommand(conn, command.Address, command.Command)

	if err != nil {
		return err
	}

	if strings.Contains(resp, "OK") {
		return nil
	}

	return fmt.Errorf("Incorrect response for command (%s). (Response: %s)", command, resp)
}
