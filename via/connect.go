package via

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/byuoitav/common/log"

	"github.com/fatih/color"
)

// SendCommand opens a connection with <addr> and sends the <command> to the via, returning the response from the via, or an error if one occured.
func SendCommand(command Command, addr string) (string, error) {
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
	login(conn)

	// write command
	if len(command.Command) > 0 {
		command.addAuth(false)
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
		color.Set(color.FgBlue)
		log.L.Infof("Response from device: %s", resp)
	}

	return string(resp), nil
}

func login(conn *net.TCPConn) error {
	var cmd Command
	cmd.addAuth(true)
	cmd.Command = "Login"

	log.L.Infof("Logging in...")

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
		color.Set(color.FgMagenta)
	}

	conn.Write(b)
	return nil
}

// AddAuth adds auth onto the command
// changed: Made function Public
func (c *Command) addAuth(password bool) {
	c.Username = via.Username
	if password {
		c.Password = via.Password
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

// Create a presistent connection in order to catch actions and events that are printed
// out on console. This includes login, logoff, media presentation, and sharing events
func PersistConnection(addr string) (*net.TCPConn, error) {
	defer color.Unset()
	color.Set(color.FgCyan)

	// get the connection
	log.L.Infof("Opening persistent telnet connection for reading events from %s", addr)
	pconn, err := getConnection(addr)
	if err != nil {
		return nil, err
	}

	// login
	login(pconn)

	return pconn, nil
}

// VolumeParse parser to pull out the volume level from the VIA API returned string
func VolumeParse(vollevel string) (int, error) {
	re := regexp.MustCompile("[0-9]+")
	vol := re.FindString(vollevel)
	vfin, err := strconv.Atoi(vol)
	if err != nil {
		err = fmt.Errorf("Error converting response: %s", err.Error())
		color.Set(color.FgRed)
		log.L.Infof("%s", err.Error())
		color.Unset()
		return 0, err
	}
	return vfin, nil
}
