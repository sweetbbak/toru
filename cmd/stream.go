package main

import (
	"context"
	"fmt"
	"time"
	"toru/pkg/libtorrent"
	"toru/pkg/player"

	"github.com/dustin/go-humanize"
)

// takes any type of torrent file/url/magnet
func StreamTorrent(torfile string, cl *libtorrent.Client) (context.CancelFunc, string, error) {
	t, err := cl.AddMagnet(torfile)
	if err != nil {
		return nil, "", err
	}

	ctx, cancel := context.WithCancel(context.Background())
	link := cl.ServeTorrent(ctx, t)

	go func() {
		for !t.Complete.Bool() {
			c := t.BytesCompleted()
			total := t.Length()
			s := humanize.Bytes(uint64(c))
			x := humanize.Bytes(uint64(total))
			numpeers := len(t.PeerConns())
			fmt.Print("\r\x1b[1K") // return to beginning of line and clear
			fmt.Printf("Downloaded (%v/%v) from [%v] Peers...\n", s, x, numpeers)
			time.Sleep(time.Millisecond * 500)
		}
		println("Complete")
	}()

	fmt.Println(link)
	return cancel, link, nil
}

func PlayTorrent(cl *libtorrent.Client, magnet string) error {
	cancel, l, err := StreamTorrent(magnet, cl)
	if err != nil {
		return err
	}
	p, err := player.NewPlayer()
	if err != nil {
		return err
	}

	// get a new player and start the media
	proc, err := p.Open(l)
	if err != nil {
		return err
	}

	// wait for player to close
	proc.Wait()
	// cancel torrent server
	cancel()
	return nil
}

func runStream(cl *libtorrent.Client) error {
	if streaming.Magnet != "" {
		err := PlayTorrent(cl, streaming.Magnet)
		if err != nil {
			return err
		}
	}
	return nil
}
