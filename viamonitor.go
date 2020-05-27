package kramer

import (
	"bufio"
	"context"
	"encoding/xml"
	"net"
)

type Message struct {
	EventType string
	Action    string
	User      string
	State     string
}

type PersistentViaConnection struct {
	Conn   net.Conn
	Reader *bufio.Reader
}

// Implementing ping functionality here for via-control
// Do I really need to gather the output of this considering I am using this
// to just keep the connection alive -- you decide!
func (v *Via) Ping(conn *PersistentViaConnection) error {
	var cmd command
	cmd.Username = "su"
	cmd.Command = "IpInfo"
	m, err := xml.Marshal(cmd)
	if err != nil {
		return err
	}

	_, err = conn.Conn.Write(m)
	if err != nil {
		return err
	}
	v.Infof("Sending Ping to %s", v.Address)
	return nil
}

//Create a persistent connection in order to catch actions and events that are printed
//out on console. This includes login, logoff, media presentation, and sharing events
func (v *Via) PersistConnection(ctx context.Context) (*PersistentViaConnection, error) {
	// get the connection
	v.Infof("Opening persistent telnet connection for reading events from %s", v.Address)
	gconn, err := getConnection(v.Address)
	if err != nil {
		return nil, err
	}

	// login
	err = v.login(ctx, gconn)
	if err != nil {
		v.Debugf("Houston, we have a problem logging in. The login failed")
		return nil, err
	}

	return &PersistentViaConnection{
		Conn:   gconn,
		Reader: bufio.NewReader(gconn),
	}, nil
}

// This part is actually part of the VIA-Controller Microservice.
// It is an example of how the messenger works (In the VIA-Controller, it's called in readPump)
/*
func (c *PersistantViaConnection) NextMessage() (Message, error) {
	var msg Message

	buf, err := c.reader.ReadBytes('\x0d')
	if err != nil {
		return msg, err
	}

	// parse buf
	buf
	// set message fields
	msg.EventType = "current-user-count"

	return msg, nil
}
*/
// example
//func main() {
//	via := &kramer.Via{
//		Address:
//		Unaem:
//	}
//
//	conn, err := via.PersistantViaConnection()
//
//	for {
//		msg, err := conn.NextMessage()
//		if err != nil {
//
//		}
//
//		// send event from message
//	}
//}
