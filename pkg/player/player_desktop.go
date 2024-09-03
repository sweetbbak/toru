//go:build (linux || darwin || windows) && !android

package player

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/sweetbbak/toru/pkg/player/mpv"
)

var desktopPlayers = []GenericPlayer{
	{Name: "mpv", Args: []string{"--script={{GetHelperScriptPath}}", "{{.URL}}"}, GetHelperScriptPath: mpv.GetScriptPath},
	{Name: "vlc", Args: []string{"{{.URL}}"}},
	{Name: "mplayer", Args: []string{"{{.URL}}"}},
	{Name: "iina", Args: []string{"--no-stdin", "--keep-running", "{{.URL}}"}},
	{Name: "catt", Args: []string{"cast", "{{.URL}}"}},
}

func findPlayer() (GenericPlayer, error) {
	for _, p := range desktopPlayers {
		_, err := exec.LookPath(p.Name)
		if err != nil {
			continue
		} else {
			return p, nil
		}
	}

	return GenericPlayer{}, fmt.Errorf("No supported player found")
}

// get a generic player for playing media
// pass an optional string with either a player + args or a player name
// can be empty. use "{{.URL}}" as a placeholder for the url if needed
func NewPlayer(player string) (GenericPlayer, error) {
	var command []string
	if len(player) > 0 {
		command = strings.Split(player, " ")
	}

	switch len(command) {
	case 0:
		// no player provided
		// look for one or return error
		return findPlayer()
	case 1:
		// the "defualt" players on Android
		switch strings.ToLower(command[0]) {
		case "mpv":
			return desktopPlayers[0], nil
		case "vlc":
			return desktopPlayers[1], nil
		default:
			return GenericPlayer{Name: command[0], Args: command[1:]}, nil
		}
	default:
		// else we assume the user knows what they want to do here and just create a player based on what they want
		// this allows non-standard players, like a browser or terminal
		return GenericPlayer{Name: command[0], Args: command[1:]}, nil
	}
}
