package kramer

import (
	"bufio"
	"context"
	"net"

	"github.com/byuoitav/common/log"
)

type Message struct {
	EventType string
	Action    string
	User      string
	State     string
}

type PersistantViaConnection struct {
	conn   net.Conn
	reader *bufio.Reader
}

//Create a persistent connection in order to catch actions and events that are printed
//out on console. This includes login, logoff, media presentation, and sharing events
func (v *Via) PersistConnection(ctx context.Context) (*PersistantViaConnection, error) {
	// get the connection
	log.L.Infof("Opening persistent telnet connection for reading events from %s", v.Address)
	conn, err := getConnection(v.Address)
	if err != nil {
		return nil, err
	}

	// login
	err = v.login(ctx, conn)
	if err != nil {
		log.L.Debugf("Houston, we have a problem logging in. The login failed")
		return nil, err
	}

	return &PersistantViaConnection{
		conn:   conn,
		reader: bufio.NewReader(conn),
	}, nil
}

func (c *PersistantViaConnection) NextMessage() (Message, error) {
	var msg Message

	buf, err := c.reader.ReadBytes('\x0d')
	if err != nil {
		return msg, err
	}

	// parse buf

	// set message fields
	msg.EventType = "current-user-count"

	return msg, nil
}

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
