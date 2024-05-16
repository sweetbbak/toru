package mpv

import (
	_ "embed"
	"os"
	"path"
)

//go:embed script.lua
var script []byte

func WriteScript() error {
	if _, err := os.Stat(getScriptPath()); !os.IsNotExist(err) {
		return nil
	}

	return os.WriteFile(getScriptPath(), []byte(script), 0644)
}

func getScriptPath() string {
	return path.Join(os.TempDir(), "toru-mpv-script.lua")
}

func GetScriptPath() string {
	WriteScript()
	return getScriptPath()
}
