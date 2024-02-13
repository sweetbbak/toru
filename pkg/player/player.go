package player

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// GenericPlayer represents most players. The stream URL will be appended to the arguments.
type GenericPlayer struct {
	Name string
	Args []string
}

// Player opens a stream URL in a video player.
type Player interface {
	Open(url string) error
}

var genericPlayers = []GenericPlayer{
	{Name: "mpv", Args: []string{"mpv"}},
	{Name: "vlc", Args: []string{"vlc"}},
	{Name: "mplayer", Args: []string{"mplayer"}},
}

// Open the given stream in a GenericPlayer.
func (p GenericPlayer) Open(url string) (*os.Process, error) {
	command := []string{}
	if runtime.GOOS == "darwin" {
		command = []string{"open", "-a"}
	}

	command = append(command, p.Args...)
	command = append(command, url)

	// It is the user's responsibility to pass the correct arguments to open the url.
	cmd := exec.Command(command[0], command[1:]...)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	return cmd.Process, nil
}

// openPlayer opens a stream using the specified player and port based on user specs.
func OpenPlayer(player string, url string) (*os.Process, error) {
	command := strings.Split(player, " ")
	if len(command) < 1 {
		return nil, fmt.Errorf("command cannot be empty")
	}

	command = append(command, url)

	cmd := exec.Command(command[0], command[1:]...)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	return cmd.Process, nil
}

// get a generic player for playing media
func NewPlayer() (GenericPlayer, error) {
	for _, p := range genericPlayers {
		_, err := exec.LookPath(p.Name)
		if err != nil {
			continue
		} else {
			return p, nil
		}
	}
	return GenericPlayer{}, fmt.Errorf("No supported player found")
}

// joinPlayerNames returns a list of supported video players ready for display.
func joinPlayerNames() string {
	names := make([]string, len(genericPlayers))
	for i := range genericPlayers {
		names[i] = genericPlayers[i].Name
	}
	return strings.Join(names, ", ")
}
