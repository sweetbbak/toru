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
	{Name: "iina", Args: []string{"iina", "--no-stdin", "--keep-running"}},
	{Name: "catt", Args: []string{"catt", "cast"}},
}

// ganked from ani-cli
var androidPlayers = []GenericPlayer{
	{Name: "mpv", Args: []string{"am", "start", "--user", "0", "-a", "android.intent.action.VIEW", "-d", "{URL}", "-n", "is.xyz.mpv/.MPVActivity"}},
	{Name: "vlc", Args: []string{"am", "start", "--user", "0", "-a", "android.intent.action.VIEW", "-d", "{URL}", "-n", "org.videolan.vlc/org.videolan.vlc.gui.video.VideoPlayerActivity"}},
}

// Open the given stream in a GenericPlayer.
func (p GenericPlayer) Open(url string) (*os.Process, error) {
	command := []string{}
	var isAndroid bool

	switch runtime.GOOS {
	case "darwin":
		command = []string{"open", "-a"}
	case "android":
		isAndroid = true
		for _, u := range p.Args {
			switch u {
			case "url", "URL", "{url}", "{URL}":
				u = url
			}
		}
	}

	command = append(command, p.Args...)

	if isAndroid {
		return openAndroid(command)
	}

	command = append(command, url)
	// It is the user's responsibility to pass the correct arguments to open the url.
	cmd := exec.Command(command[0], command[1:]...)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	return cmd.Process, nil
}

func openAndroid(command []string) (*os.Process, error) {
	cmd := exec.Command(command[0], command[1:]...)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	return cmd.Process, nil
}

// openPlayer opens a stream using the specified player and port based on user specs.
func PlayerCMD(player string, url string) (*os.Process, error) {
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
// pass an optional string with either a player + args or a player name
// can be empty. use "{url}" as a placeholder for the url if needed
func NewPlayer(player string) (GenericPlayer, error) {
	var command []string
	if len(command) > 0 {
		command = strings.Split(player, " ")
	}

	os := runtime.GOOS
	switch os {
	case "android":
		if len(command) > 0 {
			switch strings.ToLower(command[0]) {
			case "mpv":
				return androidPlayers[0], nil
			case "vlc":
				return androidPlayers[1], nil
			default:
				return GenericPlayer{Name: command[0], Args: command[1:]}, nil
			}
		} else {
			return GenericPlayer{}, fmt.Errorf("Unknown or empty player specified by user")
		}
	}

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
