package player

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ganked from ani-cli
var androidPlayers = []GenericPlayer{
	{Name: "mpv", Args: []string{"am", "start", "--user", "0", "-a", "android.intent.action.VIEW", "-d", "{{.URL}}", "-n", "is.xyz.mpv/.MPVActivity"}},
	{Name: "vlc", Args: []string{"am", "start", "--user", "0", "-a", "android.intent.action.VIEW", "-d", "{{.URL}}", "-n", "org.videolan.vlc/org.videolan.vlc.gui.video.VideoPlayerActivity"}},
}

// Open the given stream in a GenericPlayer.
func (p GenericPlayer) Open(media MediaEntry) (*os.Process, error) {
	cmd := exec.Command(p.Args[0], p.expandArgs(media)[1:]...)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	return cmd.Process, nil
}

// get a generic player for playing media
// pass an optional string with either a player + args or a player name
// can be empty. use "{{.URL}}" as a placeholder for the url if needed
func NewPlayer(player string) (GenericPlayer, error) {
	var command []string
	if len(command) > 0 {
		command = strings.Split(player, " ")
	}

	if len(command) > 0 {
		switch strings.ToLower(command[0]) {
		case "mpv":
			return androidPlayers[0], nil
		case "vlc":
			return androidPlayers[1], nil
		default:
			return GenericPlayer{Name: command[0], Args: command[1:]}, nil
		}
	}

	return GenericPlayer{}, fmt.Errorf("Unknown or empty player specified by user")
}
