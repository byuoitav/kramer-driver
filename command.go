package kramer

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"time"

	"github.com/fatih/color"

	"github.com/byuoitav/common/log"
)

// Takes a command and sends it to the address, and returns the devices response to that command
func SendCommand(conn *net.TCPConn, address, command string) (resp string, err error) {
	defer color.Unset()

	resp, err = writeCommand(conn, command)
	if err != nil {
		return "", err
	}

	color.Set(color.FgBlue)
	log.L.Infof("Response from device: %s", resp)
	return resp, nil
}

func writeCommand(conn *net.TCPConn, command string) (string, error) {
	command = strings.Replace(command, " ", string(SPACE), -1)
	color.Set(color.FgMagenta)
	log.L.Infof("Sending command %s", command)
	color.Unset()
	command += string(CARRIAGE_RETURN) + string(LINE_FEED)
	conn.Write([]byte(command))

	// get response
	resp, err := readUntil(LINE_FEED, conn, 5)
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

func readUntil(delimeter byte, conn *net.TCPConn, timeoutInSeconds int) ([]byte, error) {
	conn.SetReadDeadline(time.Now().Add(time.Duration(int64(timeoutInSeconds)) * time.Second))

	buffer := make([]byte, 128)
	message := []byte{}

	for !charInBuffer(delimeter, buffer) {
		_, err := conn.Read(buffer)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error reading response: %s", err.Error()))
			log.L.Infof("%s", err.Error())
			return message, err
		}

		message = append(message, buffer...)
	}

	return removeNil(message), nil
}

func readAll(conn *net.TCPConn, timeoutInSeconds int) ([]byte, error) {
	conn.SetReadDeadline(time.Now().Add(time.Duration(int64(timeoutInSeconds)) * time.Second))

	bytes, err := ioutil.ReadAll(conn)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error reading response: %s", err.Error()))
		return []byte{}, err
	}

	return removeNil(bytes), nil
}

func removeNil(b []byte) (ret []byte) {
	for _, c := range b {
		switch c {
		case '\x00':
			break
		default:
			ret = append(ret, c)
		}
	}
	return ret
}

func charInBuffer(toCheck byte, buffer []byte) bool {
	for _, b := range buffer {
		if toCheck == b {
			return true
		}
	}

	return false
}

func logError(e string) {
	color.Set(color.FgRed)
	log.L.Infof("%s", e)
	color.Unset()
}
