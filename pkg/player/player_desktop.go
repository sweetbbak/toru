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

// get a generic player for playing media
// pass an optional string with either a player + args or a player name
// can be empty. use "{{.URL}}" as a placeholder for the url if needed
func NewPlayer(player string) (GenericPlayer, error) {
	var command []string
	if len(command) > 0 {
		command = strings.Split(player, " ")
	}

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
