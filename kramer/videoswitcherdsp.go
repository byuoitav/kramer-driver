package kramer

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/byuoitav/connpool"
)

type KramerVP558 struct {
	Address string
	Log     Logger

	pool *connpool.Pool
}

// var (
// 	_defaultTTL   = 30 * time.Second
// 	_defaultDelay = 500 * time.Millisecond
// )

type KramerVP558options struct {
	ttl    time.Duration
	delay  time.Duration
	logger Logger
}

type KramerVP558Option interface {
	apply(*KramerVP558options)
}

type KramerVP558optionFunc func(*KramerVP558options)

func NewVideoSwitcherDsp(addr string, opts ...KramerVP558Option) *KramerVP558 {
	options := KramerVP558options{
		ttl:   _defaultTTL,
		delay: _defaultDelay,
	}

	for _, o := range opts {
		o.apply(&options)
	}

	vsdsp := &KramerVP558{
		Address: addr,
		pool: &connpool.Pool{
			TTL:    options.ttl,
			Delay:  options.delay,
			Logger: options.logger,
		},
	}

	vsdsp.pool.NewConnection = func(ctx context.Context) (net.Conn, error) {
		d := net.Dialer{}
		conn, err := d.DialContext(ctx, "tcp", vsdsp.Address+":5000")
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

	return vsdsp
}

// SendCommand sends the byte array to the desired address of projector
func (vsdsp *KramerVP558) SendCommand(ctx context.Context, cmd []byte) ([]byte, error) {
	var resp []byte

	err := vsdsp.pool.Do(ctx, func(conn connpool.Conn) error {
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
