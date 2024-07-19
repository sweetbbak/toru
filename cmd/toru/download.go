package main

import (
	"context"
	"fmt"

	"github.com/pterm/pterm"
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

	success, _ := pterm.DefaultSpinner.Start("getting torrent info")

	t, err := cl.AddTorrent(tfile)
	if err != nil {
		return err
	}

	success.Success("Success!")
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		Progress(t, ctx)
	}()

	t.DownloadAll()
	if cl.TorrentClient.WaitAll() {
		cancel()
		return nil
	} else {
		cancel()
		return fmt.Errorf("Unable to completely download torrent: %s", t.Name())
	}
}
