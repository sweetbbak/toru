package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/sweetbbak/toru/pkg/libtorrent"
	"github.com/sweetbbak/toru/pkg/player"
	"github.com/sweetbbak/toru/pkg/search"
)

// takes any type of torrent file/url/magnet, adds it to the client and streams it
func StreamTorrent(torfile string, cl *libtorrent.Client) (string, error) {
	t, err := cl.AddTorrent(torfile)
	if err != nil {
		return "", err
	}

	link := cl.ServeTorrent(t)

	// consider deleting this as it sometimes conflicts with the fzf user interface
	go func() {
		for !t.Complete.Bool() {
			c := t.BytesCompleted()
			total := t.Length()
			s := humanize.Bytes(uint64(c))
			x := humanize.Bytes(uint64(total))
			numpeers := len(t.PeerConns())
			fmt.Printf("\x1b[2K\rDownloaded (%v/%v) from [%v] Peers...", s, x, numpeers)
			time.Sleep(time.Millisecond * 500)
		}
		println("Complete")
	}()

	fmt.Println(link)
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
