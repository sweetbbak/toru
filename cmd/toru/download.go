package main

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/sweetbbak/toru/pkg/libtorrent"
)

func DownloadTorrent(cl *libtorrent.Client) error {
	var tfile string

	if download.TorrentFile != "" {
		tfile = download.TorrentFile
	} else if download.Args.Query != "" {
		tfile = download.Args.Query
	} else {
		return fmt.Errorf("download: missing argument (magnet|torrent|url) OR --torrent flag")
	}

	t, err := cl.AddTorrent(tfile)
	if err != nil {
		return err
	}

	go func() {
		name := t.Name()
		fmt.Printf("Downloading: '%s'\n", name)
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

	t.DownloadAll()
	if cl.TorrentClient.WaitAll() {
		return nil
	} else {
		return fmt.Errorf("Unable to completely download torrent: %s", t.Name())
	}
}
