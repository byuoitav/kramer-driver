package via

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"net"

	"time"

	"github.com/byuoitav/common/log"
)

var ViaUser string
var ViaPass string

func (v *VIA) importUser() (ViaUser, ViaPass string) {
	ViaUser = v.Username
	ViaPass = v.Password
	return ViaUser, ViaPass
}

// SendCommand opens a connection with <addr> and sends the <command> to the via, returning the response from the via, or an error if one occured.
func (v *VIA) SendCommand(command Command, addr string) (string, error) {
	Username, Password := v.importUser()
	// get the connection
	log.L.Infof("Opening telnet connection with %s", addr)
	conn, err := getConnection(addr)
	if err != nil {
		return "", err
	}

	timeoutDuration := 7 * time.Second

	// Set Read Connection Duration
	conn.SetReadDeadline(time.Now().Add(timeoutDuration))

	// login
	login(conn, Username, Password)

	// write command
	if len(command.Command) > 0 {
		command.addAuth(Username, Password, false)
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

func login(conn *net.TCPConn, ViaUser, ViaPass string) error {
	var cmd Command
	cmd.addAuth(ViaUser, ViaPass, true)
	cmd.Command = "Login"

	log.L.Infof("Logging in...")
	log.L.Infof("Username: %s", ViaUser)
	err := cmd.writeCommand(conn)
	if err != nil {
		return err
	}

	log.L.Infof("Login successful")

	return nil
}

func (c *Command) writeCommand(conn *net.TCPConn) error {

	// read welcome message
	reader := bufio.NewReader(conn)
	_, err := reader.ReadBytes('\n')
	if err != nil {
		err = fmt.Errorf("error reading from system: %s", err.Error())
		log.L.Error(err.Error())
		return err
	}

	b, err := xml.Marshal(c)
	if err != nil {
		return err
	}

	if len(c.Password) == 0 {
		log.L.Infof("Sending command: %s", b)
	}

	conn.Write(b)
	return nil
}

// AddAuth adds auth onto the command
// changed: Made function Public
func (c *Command) addAuth(ViaUser string, ViaPass string, password bool) {
	c.Username = ViaUser
	if password {
		c.Password = ViaPass
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

//Create a persistent connection in order to catch actions and events that are printed
//out on console. This includes login, logoff, media presentation, and sharing events
func PersistConnection(addr string) (*net.TCPConn, error) {
	// get the connection
	log.L.Infof("Opening persistent telnet connection for reading events from %s", addr)
	pconn, err := getConnection(addr)
	if err != nil {
		return nil, err
	}

	// login
	login(pconn, ViaUser, ViaPass)

	return pconn, nil
}
