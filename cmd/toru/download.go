package main

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/sweetbbak/toru/pkg/libtorrent"
)

func DownloadTorrent(cl *libtorrent.Client) error {
	if download.Args.Query == "" {
		return fmt.Errorf("download: missing argument (magnet|torrent|url)")
	}

	t, err := cl.AddTorrent(download.Args.Query)
	if err != nil {
		return err
	}

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

	t.DownloadAll()
	cl.TorrentClient.WaitAll()
	return nil
}
