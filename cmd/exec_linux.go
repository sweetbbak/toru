package main

import (
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func playLink(link string) error {
	// wait until link is being served by our client (takes a couples seconds)
	timeout := time.Duration(time.Second * 10)
	client := http.Client{
		Timeout: timeout,
	}
	_, err := client.Get(link)

	var player string
	if opts.Player != "" {
		player = opts.Player
	} else {
		// [TODO): add utils func to find players
		player = "mpv"
	}

	cmd := exec.Command(player, link)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// setsid on Windows? what about android?
	// on windows its:
	//  cmd.SysProcAttr = &syscall.SysProcAttr{
	// CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	// }
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
