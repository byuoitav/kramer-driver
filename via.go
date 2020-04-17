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

// comman: Struct used to build the commands that need to be sent to the VIA
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

func (v *Via) GetVolume(ctx context.Context) (int, error) {
	var cmd command
	cmd.Command = "Vol"
	cmd.Param1 = "Get"

	log.L.Infof("Sending command to get VIA Volume to %s", v.Address)
	// Note: Volume Get command in VIA API doesn't have any error handling so it only returns Vol|Get|XX or nothing
	// I am still checking for errors just in case something else fails during execution
	vollevel, err := v.SendCommand(ctx, command)
	if err != nil {

	}

	return v.volumeParse(vollevel)
}

// GetInfo: needed by the DSP drivers implementation.  Will get hardware information
func (v *Via) GetInfo(ctx context.Context) (interface{}, error) {
	var info interface{}
	return info, fmt.Errorf("GetInfo has not been implemented in this version of the driver")
}

// volumeParse parser to pull out the volume level from the VIA API returned string
func (v *Via) volumeParse(vollevel string) (int, error) {
	re := regexp.MustCompile("[0-9]+")
	vol := re.FindString(vollevel)

	vfin, err := strconv.Atoi(vol)
	if err != nil {
		err = fmt.Errorf("Error converting response: %s", err.Error())
		log.L.Infof("%s", err.Error())
		return 0, err
	}

	return vfin, nil
}

// Reboot: Reboot a VIA using the API
func (v *Via) Reboot(ctx context.Context) error {
	var command command
	command.Command = viaReboot

	log.L.Infof("Sending command %s to %s", viaReboot, v.Address)

	_, err := v.SendCommand(ctx, command)
	if err != nil {
		return err
	}

	return nil
}

// Reset: Reset a VIA sessions - Causes VIAAdmin to log out and log back in
func (v *Via) Reset(ctx context.Context) error {
	var command command
	command.Command = viaReset

	log.L.Infof("Sending command %s to %s", viaReset, v.Address)

	resp, err := v.SendCommand(ctx, command)
	if err != nil {
		return err
	}

	if strings.Contains(resp, viaReset) && strings.Contains(resp, "1") {
		return nil
	}

	return errors.New(fmt.Sprintf("Incorrect response for command. (Response: %s)", resp))
}

// Get the Room Code and return the current room code as a string
func (v *VIA) RoomCode(ctx context.Context) (string, error) {
	var command Command
	command.Command = "RCode"
	command.Param1 = "Get"
	command.Param2 = "Code"

	log.L.Infof("Sending command to get current room code to %s", v.Address)
	// Note: RCode Get command in VIA API doesn't have any error handling so it only returns RCode|Get|Code|XXXX or nothing
	resp, err := v.sendCommand(ctx, command)
	if err != nil {
		return "", err
	}
	split := strings.Split(resp, "|")
	if len(split) != 4 {
		return "", fmt.Errorf("Unknown response %v", resp)
	}

	roomcode := strings.TrimSpace(split[3])

	return roomcode, nil
}

// SetAlert - Send an alert to the VIA
func (v *VIA) SetAlert(ctx context.Context, message string) error {
	log.L.Infof("Sending Alert to %v", v.Address)
	resp, err := v.Alert(ctx, alertMessage)
	if err != nil {
		log.L.Debugf("Failed to send alrt message to %s", v.Address)
		return fmt.Errorf("Error sending alert message: %v", err)
	}
	log.L.Infof("Alert: %s - Sent", resp)
	return nil
}

func (v *Via) SetVolume(ctx context.Context, volume int) (string, error) {
	var command command
	command.Command = "Vol"
	command.Param1 = "Set"
	command.Param2 = strconv.Itoa(volume)

	log.L.Infof("Sending volume set command to %s", v.Address)

	resp, err := v.SendCommand(ctx, command)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error in setting volume on %s", v.Address))
	}

	return resp, nil

}

func (v *VIA) Alert(ctx context.Context, AlertMessage string) (string, error) {
	var command Command
	command.Command = "IAlert"
	command.Param1 = AlertMessage
	command.Param2 = "0"
	command.Param3 = "5"

	log.L.Infof("Sending an alert message -%s- to %s", AlertMessage, v.Address)

	resp, err := v.SendCommand(ctx, command)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error in setting volume on %s", v.Address))
	}

	return resp, nil
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
		command.addAuth(v.Username, v.Password, false)
		command.writeCommand(conn)
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

	cmd.addAuth(v.Username, v.Password, true)
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
	err = cmd.writeCommand(conn)
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

// TODO remove
func (c *Command) writeCommand(conn *net.TCPConn) error {
	b, err := xml.Marshal(c)
	if err != nil {
		return err
	}

	return conn.Write(b)
}

// AddAuth adds auth onto the command
// TODO remove
// changed: Made function Public
func (c *Command) addAuth(viaUser string, viaPass string, password bool) {
	c.Username = viaUser
	if password {
		c.Password = viaPass
	}
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
