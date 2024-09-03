//go:build linux && !android
// +build linux,!android

package player

import (
	"os"
	"os/exec"
)

// Open the given stream in a GenericPlayer.
func (p GenericPlayer) Open(media MediaEntry) (*os.Process, error) {
	cmd := exec.Command(p.Name, p.expandArgs(media)...)

	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	return cmd.Process, nil
}
