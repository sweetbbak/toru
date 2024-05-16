package player

import (
	"os"
	"os/exec"
)

// Open the given stream in a GenericPlayer.
func (p GenericPlayer) Open(media MediaEntry) (*os.Process, error) {
	args := []string{"-a", p.Name}
	args = append(args, p.expandArgs(media)...)
	cmd := exec.Command("open", args...)

	// It is the user's responsibility to pass the correct arguments to open the url.
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	return cmd.Process, nil
}
