package kramer

import (
	"bufio"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/byuoitav/common/log"
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
/*
// command: Struct used to build the XML commands that need to be sent to the VIA
type command struct {
	XMLName  xml.Name `xml:"P"`
	Username string   `xml:"UN"`
	Password string   `xml:"Pwd"`
	Command  string   `xml:"Cmd"`
	Param1   string   `xml:"P1,omitempty"`
	Param2   string   `xml:"P2,omitempty"`
	Param3   string   `xml:"P3,omitempty"`
	Param4   string   `xml:"P4,omitempty"`
	Param5   string   `xml:"P5,omitempty"`
	Param6   string   `xml:"P6,omitempty"`
	Param7   string   `xml:"P7,omitempty"`
	Param8   string   `xml:"P8,omitempty"`
	Param9   string   `xml:"P9,omitempty"`
	Param10  string   `xml:"P10,omitempty"`
}
*/

// These functions fulfill the DSP driver requirements 
// GetVolumeByBlock: opening a connection with the VIAs and then return the volume for the device
func (v *Via) GetVolumeByBlock(ctx context.Context, block string) (int, error) {
	return v.GetVolume(ctx)
}

// SetVolumeByBlock: Connect and set the volume on the VIA
func (v *Via) SetVolumeByBlock(ctx context.Context, block string, volume int) error {
	return v.SetVolume(ctx, volume)
}

// GetMutedByBlock: Return error because VIAs do not support a mute function
func (v *Via) GetMutedByBlock(ctx context.Context, block string) (bool, error) {
	return false, errors.New("Error in getting mute status of VIA, Feature not supported")
}

// SetMutedByBlock: Return error because VIAs do not support mute
func (v *Via) SetMutedByBlock(ctx context.Context, block string, muted bool) error {
	return errors.New("Error setting mute status of VIA, Feature not supported")
}

// GetInfo: needed by the DSP drivers implementation.  Will get hardware information
func (v *Via) GetInfo(ctx context.Context) (interface{}, error) {
	var info interface{}
	return info, fmt.Errorf("GetInfo has not been implemented in this version of the driver")
}
// End of DSP Driver requirements

func getConnection(address string) (*net.TCPConn, error) {
	radder, err := net.ResolveTCPAddr("tcp", address+":9982")
	if err != nil {
		err = fmt.Errorf("error resolving address : %s", err.Error())
		log.L.Infof(err.Error())
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, radder)
	if err != nil {
		err = fmt.Errorf("error dialing address : %s", err.Error())
		log.L.Infof(err.Error())
		return nil, err
	}

	return conn, nil
}

// SendCommand opens a connection with <addr> and sends the <command> to the via, returning the response from the via, or an error if one occured.
func (v *Via) sendCommand(ctx context.Context, command command) (string, error) {
	//Username, Password := v.importUser()
	// get the connection
	log.L.Infof("Opening telnet connection with %s", v.Address)
	conn, err := getConnection(v.Address)
	if err != nil {
		return "", err
	}

	timeoutDuration := 7 * time.Second

	// Set Read Connection Duration
	conn.SetReadDeadline(time.Now().Add(timeoutDuration))

	// login
	//login(conn, Username, Password)
	err = v.login(ctx, conn)
	if err != nil {
		log.L.Debugf("Houston, we have a problem logging in. The login failed")
		return "", err
	}

	// write command
	if len(command.Command) > 0 {
		command.Username = v.Username
		b, err := xml.Marshal(command)
		if err != nil {
			return err
		}

		err = conn.Write(b)
		if err != nil {
			return err
		}
	}

	reader := bufio.NewReader(conn)
	resp, err := reader.ReadBytes('\n')
	if err != nil {
		err = fmt.Errorf("error reading from system: %s", err.Error())
		log.L.Error(err.Error())
		return "", err
	}

	if len(string(resp)) > 0 {
		log.L.Infof("Response from device: %s", resp)
	}

	return string(resp), nil
}

func (v *Via) login(ctx context.Context, conn *net.TCPConn) error {
	var cmd Command

	//cmd.addAuth(v.Username, v.Password, true)
	cmd.Username = v.Username
	cmd.Password = v.Password
	cmd.Command = "Login"

	// read welcome message (Only Important when we first open a connection and login)
	reader := bufio.NewReader(conn)
	_, err := reader.ReadBytes('\n')
	if err != nil {
		err = fmt.Errorf("error reading from system: %s", err.Error())
		log.L.Error(err.Error())
		return err
	}

	log.L.Infof("Logging in...")
	log.L.Debugf("Username: %s", v.Username)
	b, err := xml.Marshal(cmd)
	if err != nil {
		return err
	}
	err = conn.Write(b)
	//err = cmd.writeCommand(conn)
	if err != nil {
		return err
	}
	//reader := bufio.NewReader(conn)
	resp, err := reader.ReadBytes('\n')
	if err != nil {
		err = fmt.Errorf("error reading from system: %s", err.Error())
		log.L.Error(err.Error())
		return err
	}

	s := string(resp)

	errRx := regexp.MustCompile(`Error`)
	SuccessRx := regexp.MustCompile(`Successful`)
	respRx := errRx.MatchString(s)
	SuccessResp := SuccessRx.MatchString(s)

	if respRx == true {
		log.L.Infof("Response from device: %s", s)
		err := fmt.Errorf("Unable to login due to an error: %s", s)
		return err
	}

	if SuccessResp == true {
		log.L.Debugf("Connection is successful, We are connected: %s", s)
	}

	//log.L.Infof("Login successful")

	return nil
}

func getConnection(address string) (*net.TCPConn, error) {
	radder, err := net.ResolveTCPAddr("tcp", address+":9982")
	if err != nil {
		err = fmt.Errorf("error resolving address : %s", err.Error())
		log.L.Infof(err.Error())
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, radder)
	if err != nil {
		err = fmt.Errorf("error dialing address : %s", err.Error())
		log.L.Infof(err.Error())
		return nil, err
	}

	return conn, nil
}
