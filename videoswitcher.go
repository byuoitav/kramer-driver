package kramer

import (
	"context"
	"net"
	"time"

	"github.com/byuoitav/connpool"
)

type VideoSwitcher struct {
	Address string
	Log     Logger

	pool *connpool.Pool
}

// type Kramer struct {
// 	Address     string
// 	Type        DeviceType
// 	Room_System string
// 	System_ID   string
// }

var (
	_defaultTTL   = 30 * time.Second
	_defaultDelay = 500 * time.Millisecond
)

type options struct {
	ttl    time.Duration
	delay  time.Duration
	logger Logger
}

type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func NewVideoSwitcher(addr string, opts ...Option) *VideoSwitcher {
	options := options{
		ttl:   _defaultTTL,
		delay: _defaultDelay,
	}

	for _, o := range opts {
		o.apply(&options)
	}

	vs := &VideoSwitcher{
		Address: addr,
		pool: &connpool.Pool{
			TTL:    options.ttl,
			Delay:  options.delay,
			Logger: options.logger,
		},
	}

	vs.pool.NewConnection = func(ctx context.Context) (net.Conn, error) {
		d := net.Dialer{}
		return d.DialContext(ctx, "tcp", vs.Address+":5000")

		// color.Set(color.FgMagenta)
		// log.L.Infof("Opening telnet connection with %s", vs.Address)
		// color.Unset()

		// addr, err := net.ResolveTCPAddr("tcp", vs.Address+":5000")
		// if err != nil {
		// 	return nil, err
		// }

		// conn, err := net.DialTCP("tcp", nil, addr)
		// if err != nil {
		// 	return nil, err
		// }

		// if readWelcome {
		// 	color.Set(color.FgMagenta)
		// 	log.L.Infof("Reading welcome message")
		// 	color.Unset()
		// 	_, err := readUntil(CARRIAGE_RETURN, conn, 3)
		// 	if err != nil {
		// 		return conn, err
		// 	}
		// }

		// return conn, err
	}

	return vs
}
