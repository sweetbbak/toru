package main

import (
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

	go func() {
		Progress(t)
	}()

	t.DownloadAll()
	if cl.TorrentClient.WaitAll() {
		return nil
	} else {
		return fmt.Errorf("Unable to completely download torrent: %s", t.Name())
	}
}
