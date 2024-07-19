package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/pterm/pterm"
	"github.com/sweetbbak/toru/pkg/libtorrent"
	"github.com/sweetbbak/toru/pkg/player"
	"github.com/sweetbbak/toru/pkg/search"
)

// truncate a string to a specific length
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

// print the progress of a specific torrent until context cancel is triggered
func Progress(t *torrent.Torrent, ctx context.Context) {
	title := TruncateString(t.Name(), 33)
	fmt.Println(title)

	if t == nil {
		log.Fatal("torrent is nil")
	}

	p, err := pterm.DefaultProgressbar.WithTotal(100).Start()
	if err != nil {
		log.Println(err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			p.Stop()
			return
		default:
			pc := float64(t.BytesCompleted()) / float64(t.Length()) * 100
			numpeers := len(t.PeerConns())
			p.Increment().Current = int(pc)
			p.UpdateTitle(fmt.Sprintf("peers [%v]", numpeers))
			time.Sleep(time.Millisecond * 50)
		}
	}
}

// takes any type of torrent file/url/magnet, adds it to the client and streams it
// torfile is any of magnet, link or path to .torrent file.
func StreamTorrent(torfile string, cl *libtorrent.Client) (player.MediaEntry, *torrent.Torrent, error) {
	success, _ := pterm.DefaultSpinner.Start("getting torrent info")
	t, err := cl.AddTorrent(torfile)
	if err != nil {
		return player.MediaEntry{}, nil, err
	}
	success.Success("Success!")

	files := t.Files()
	filesCount := len(files)

	var link string

	if filesCount != 1 {
		fpath, err := fzfEpisodes(files)
		if err != nil {
			return player.MediaEntry{}, nil, err
		}

		link = cl.ServeTorrentEpisode(t, fpath)
	} else {
		if len(files) < 1 {
			return player.MediaEntry{}, nil, fmt.Errorf("oops something went wrong :P - no detected files in torrent in func stream torrent")
		}

		file := files[0].DisplayPath()
		link = cl.ServeTorrentEpisode(t, file)
	}

	return player.MediaEntry{URL: link}, t, nil
}

// play a single torrent from a provided magnet, torrent or torrent URL
func PlayTorrent(cl *libtorrent.Client, magnet string) error {
	link, t, err := StreamTorrent(magnet, cl)
	if err != nil {
		return err
	}

	p, err := player.NewPlayer(options.Player)
	if err != nil {
		return err
	}

	// get a new player and start the media
	proc, err := p.Open(link)
	if err != nil {
		return err
	}

	// run progress in the background and cancel the progress bar when the player exits
	ctx, cancel := context.WithCancel(context.Background())
	go Progress(t, ctx)

	var px *os.ProcessState

	// wait for player to close
	px, err = proc.Wait()

	for !px.Exited() {
		time.Sleep(time.Millisecond * 500)
	}

	cancel()

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
