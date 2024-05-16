package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/pterm/pterm"
	"github.com/sweetbbak/toru/pkg/libtorrent"
	"github.com/sweetbbak/toru/pkg/player"
	"github.com/sweetbbak/toru/pkg/search"
)

func TruncateString(s string, length int) string {
	var l int
	var sb strings.Builder

	// early return if string is shorter then requested length
	if length >= len(s) {
		return s
	}

	for _, r := range s {
		if l <= length {
			sb.WriteRune(r)
		} else {

		}
		l++
	}
	return sb.String()
}

func Progress(t *torrent.Torrent) {
	title := TruncateString(t.Name(), 33)
	fmt.Println(title)

	p, _ := pterm.DefaultProgressbar.WithTotal(100).Start()

	for !t.Complete.Bool() {
		pc := float64(t.BytesCompleted()) / float64(t.Length()) * 100
		numpeers := len(t.PeerConns())
		p.Increment().Current = int(pc)
		p.UpdateTitle(fmt.Sprintf("peers [%v]", numpeers))
		time.Sleep(time.Millisecond * 50)
	}
}

// takes any type of torrent file/url/magnet, adds it to the client and streams it
func StreamTorrent(torfile string, cl *libtorrent.Client) (string, error) {
	success, _ := pterm.DefaultSpinner.Start("getting torrent info")
	t, err := cl.AddTorrent(torfile)
	if err != nil {
		return "", err
	}
	success.Success("Success!")

	torrentFiles := len(t.Files())
	var episode int

	if torrentFiles != 1 {
		episode, err = PromptEpisodeInRangeWithDefaultToMax(1, torrentFiles)
		if err != nil {
			return "", err
		}
	}

	link := cl.ServeTorrent(t, episode)

	// consider deleting this as it sometimes conflicts with the fzf user interface
	go func() {
		for !t.Complete.Bool() {
			Progress(t)
		}
	}()

	return link, nil
}

// play a single torrent from a provided magnet, torrent or torrent URL
func PlayTorrent(cl *libtorrent.Client, magnet string) error {
	l, err := StreamTorrent(magnet, cl)
	if err != nil {
		return err
	}
	p, err := player.NewPlayer(options.Player)
	if err != nil {
		return err
	}

	// get a new player and start the media
	proc, err := p.Open(l)
	if err != nil {
		return err
	}

	var px *os.ProcessState
	// wait for player to close
	px, err = proc.Wait()
	for !px.Exited() {
	}

	return nil
}

// stream subcommand entry point
func runStream(cl *libtorrent.Client) error {
	if streaming.Magnet != "" {
		err := PlayTorrent(cl, streaming.Magnet)
		if err != nil {
			return err
		}
	}

	if streaming.TorrentFile != "" {
		err := PlayTorrent(cl, streaming.TorrentFile)
		if err != nil {
			return err
		}
	}

	if streaming.Args.Query != "" {
		err := PlayTorrent(cl, streaming.Args.Query)
		if err != nil {
			return err
		}
	}

	if streaming.Latest {
		res, err := search.LatestAnime(streaming.Args.Query, options.Proxy, 1)
		if err != nil {
			return err
		}

		if err := SelectAndPlay(cl, res); err != nil {
			return err
		}
	}

	return nil
}
