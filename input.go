package kramer

import "context"

// GetInput returns the current input
func (vs *VideoSwitcher) getInputByOutput(ctx context.Context, output string) (string, error) {
	return "", nil
}

// SwitchInput changes the input on the given output to input
func (vs *VideoSwitcher) setInputByOutput(ctx context.Context, output, input string) error {
	return nil
}
