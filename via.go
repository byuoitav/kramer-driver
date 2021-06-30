package kramer

import (
	"bufio"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net"
	"regexp"
	//"strconv"
	"time"
)

const (
	viaReboot = "Reboot"
	viaReset  = "Reset"
)

// VIA Struct that defines general parameters needed for any VIA
type Via struct {
	Address  string
	Username string
	Password string
	Logger   Logger
}

// These functions fulfill the DSP driver requirements
// GetVolumes: opening a connection with the VIAs and then return the volume for the device
func (v *Via) Volumes(ctx context.Context, block []string) (map[string]int, error) {
	toReturn := make(map[string]int)
	resp, err := v.GetVolume(ctx)
	if err != nil {
		v.Debugf("Error retrieving volume: s%", err)
		return toReturn, err
	}

	toReturn[""] = resp
	return toReturn, nil
}

// SetVolume: Connect and set the volume on the VIA
func (v *Via) SetVolume(ctx context.Context, block string, volume int) error {
	_, err := v.SetViaVolume(ctx, volume)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to set volume for %v: %v", v.Address, err.Error()))
	}

	return nil
}

// GetInfo: needed by the DSP drivers implementation.  Will get hardware information
func (v *Via) GetInfo(ctx context.Context) (interface{}, error) {
	info, err := v.GetHardwareInfo(ctx)
	if err != nil {
		return info, fmt.Errorf("Failed to get hardware information: %s", err)
	}
	return info, nil
}

// End of DSP Driver requirements
func getConnection(address string) (*net.TCPConn, error) {
	radder, err := net.ResolveTCPAddr("tcp", address+":9982")
	if err != nil {
		err = fmt.Errorf("error resolving address : %s", err.Error())
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, radder)
	if err != nil {
		err = fmt.Errorf("error dialing address : %s", err.Error())
		return nil, err
	}

	return conn, nil
}

// SendCommand opens a connection with <addr> and sends the <command> to the via, returning the response from the via, or an error if one occured.
func (v *Via) sendCommand(ctx context.Context, cmd command) (string, error) {
	// get the connection
	v.Infof("Opening telnet connection with %s", v.Address)
	conn, err := getConnection(v.Address)
	if err != nil {
		return "", err
	}

	timeoutDuration := 7 * time.Second

	// Set Read Connection Duration
	conn.SetReadDeadline(time.Now().Add(timeoutDuration))

	// login
	err = v.login(ctx, conn)
	if err != nil {
		v.Debugf("Houston, we have a problem logging in. The login failed: %s", err)
		return "", err
	}

	// write command
	if len(cmd.Command) > 0 {
		cmd.Username = v.Username
		b, err := xml.Marshal(cmd)
		if err != nil {
			return "", err
		}

		_, err = conn.Write(b)
		if err != nil {
			return "", err
		}
	}

	reader := bufio.NewReader(conn)
	resp, err := reader.ReadBytes('\n')
	if err != nil {
		err = fmt.Errorf("error reading from system: %s", err.Error())
		v.Errorf(err.Error())
		return "", err
	}

	if len(string(resp)) > 0 {
		v.Infof("Response from device: %s", resp)
	}

	return string(resp), nil
}

func (v *Via) login(ctx context.Context, conn *net.TCPConn) error {
	var cmd command

	cmd.Username = v.Username
	cmd.Password = v.Password
	cmd.Command = "Login"

	// read welcome message (Only Important when we first open a connection and login)
	reader := bufio.NewReader(conn)
	_, err := reader.ReadBytes('\n')
	if err != nil {
		err = fmt.Errorf("error reading from system: %s", err.Error())
		v.Errorf(err.Error())
		return err
	}

	v.Debugf("Logging in...")
	v.Debugf("Username: %s", v.Username)

	b, err := xml.Marshal(cmd)
	if err != nil {
		return err
	}

	_, err = conn.Write(b)
	if err != nil {
		return err
	}

	resp, err := reader.ReadBytes('\n')
	if err != nil {
		err = fmt.Errorf("error reading from system: %s", err.Error())
		v.Errorf(err.Error())
		return err
	}

	s := string(resp)

	errRx := regexp.MustCompile(`Error`)
	SuccessRx := regexp.MustCompile(`Successful`)
	respRx := errRx.MatchString(s)
	SuccessResp := SuccessRx.MatchString(s)

	if respRx == true {
		v.Debugf("Response from device: %s", s)
		err := fmt.Errorf("Unable to login due to an error: %s", s)
		return err
	}

	if SuccessResp == true {
		v.Debugf("Connection is successful, We are connected: %s", s)
	}

	return nil
}
