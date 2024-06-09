package player

import (
	"log"
	"testing"
)

var playerTests = []struct {
	Player    string
	doesExist bool // if the player exists in the database
	expected  GenericPlayer
}{
	{"mpv", true, GenericPlayer{Name: "mpv", Args: []string{""}, GetHelperScriptPath: nil}},
	{"vlc", true, GenericPlayer{Name: "vlc", Args: []string{""}, GetHelperScriptPath: nil}},
	{"mpv --vo=kitty", true, GenericPlayer{Name: "mpv", Args: []string{"--vo=kitty"}, GetHelperScriptPath: nil}},
	{"vlc --test-arg=kitty --no-foo --bar -baz", true, GenericPlayer{Name: "vlc", Args: []string{"--test-arg=kitty", "--no-foo", "--bar", "-baz"}, GetHelperScriptPath: nil}},
	{"mpv", true, GenericPlayer{Name: "mpv", Args: []string{""}, GetHelperScriptPath: nil}},
	{"iina", true, GenericPlayer{Name: "iina", Args: []string{""}, GetHelperScriptPath: nil}},
	{"iina --test-arg=kitty --no-foo --bar -baz", true, GenericPlayer{Name: "iina", Args: []string{"--test-arg=kitty", "--no-foo", "--bar", "-baz"}, GetHelperScriptPath: nil}},
}

func TestNewPlayer(t *testing.T) {
	for _, test := range playerTests {
		p, err := NewPlayer(test.Player)
		gotValid := (err == nil)

		if gotValid {
			log.Printf("NewPlayer(%v) = %v, want %v\n", test.Player, p.Name, p.Args)
		} else {
			t.Errorf("NewPlayer(%v) = %v, want %v", test.Player, p.Name, p.Args)
		}
	}
}
