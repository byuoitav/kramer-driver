package kramer

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/byuoitav/connpool"
)

type VideoSwitcher struct {
	Address string
	Log     Logger

	pool *connpool.Pool
}

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
		conn, err := d.DialContext(ctx, "tcp", vs.Address+":5000")
		if err != nil {
			return nil, fmt.Errorf("unable to open connection: %w", err)
		}

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context done before welcome message is sent: %w", ctx.Err())
		case <-time.After(500 * time.Millisecond):
		}

		return conn, nil
	}

	return vs
}

// SendCommand sends the byte array to the desired address of projector
func (vs *VideoSwitcher) SendCommand(ctx context.Context, cmd []byte) ([]byte, error) {
	var resp []byte

	err := vs.pool.Do(ctx, func(conn connpool.Conn) error {
		_ = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

		n, err := conn.Write(cmd)
		switch {
		case err != nil:
			return err
		case n != len(cmd):
			return fmt.Errorf("wrote %v/%v bytes of command 0x%x", n, len(cmd), cmd)
		}

		resp, err = conn.ReadUntil(LINE_FEED, 3*time.Second)
		if err != nil {
			return fmt.Errorf("unable to read response: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
