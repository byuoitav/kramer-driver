package kramer

import (
	"context"
	"fmt"
)

//GetInfo .
func (vs *VideoSwitcher) GetInfo(ctx context.Context) (interface{}, error) {
	return nil, fmt.Errorf("not currently implemented")
}

//GetInfo .
func (vsdsp *VideoSwitcherDsp) GetInfo(ctx context.Context) (interface{}, error) {
	return nil, fmt.Errorf("not currently implemented")
}

//GetInfo .
func (dsp *Dsp) GetInfo(ctx context.Context) (interface{}, error) {
	return nil, fmt.Errorf("not currently implemented")
}
