package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/pterm/pterm"
	"github.com/sweetbbak/toru/pkg/libtorrent"
	"github.com/sweetbbak/toru/pkg/search"
)

//
//
//
//

func DownloadMain(cl *libtorrent.Client) error {
	var outputName string
	if download.Directory != "" {
		outputName = string(download.Directory)
	} else {
		outputName = "toru-media"
	}

	// create download dir
	if err := CreateOutput(outputName); err != nil {
		return err
	}

	// set download dir lol
	cl.SetDownloadDir(outputName)

	// no need to serve torrents
	// cl.SetServerOFF(true)
	tmp := os.TempDir()
	opt := libtorrent.SetDataDir(tmp)

	if err := cl.Init(opt); err != nil {
		return err
	}

	torrents, err := DownloadMultiple(cl)
	if err != nil {
		return err
	}

	// Create a multi printer for managing multiple printers
	multi := pterm.DefaultMultiPrinter

	var pbars []*pterm.ProgressbarPrinter

	for _, t := range torrents {
		pb, _ := pterm.DefaultProgressbar.WithTotal(100).WithWriter(multi.NewWriter()).Start(TruncateString(t.Name(), 30))
		pbars = append(pbars, pb)
	}

	_, err = multi.Start()
	if err != nil {
		log.Println(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				multi.Stop()
				return
			default:
				for i, t := range torrents {
					pb := pbars[i]
					pc := float64(t.BytesCompleted()) / float64(t.Length()) * 100
					numpeers := len(t.PeerConns())
					pb.Increment().Current = int(pc)
					pb.UpdateTitle(fmt.Sprintf("peers [%02d]", numpeers))
					time.Sleep(time.Millisecond * 5)
				}
			}
		}
	}()

	for {
		if cl.TorrentClient.WaitAll() {
			cancel()
			println("done!")
			return nil
		}
	}

}

// TODO: create a function that just handles query building
func SearchAnime(q *Search, term string) (*search.Results, error) {
	s := search.NewSearch()

	// build the query
	if q.Category != "" {
		s.Category = q.Category
	}
	if q.Filter != "" {
		s.Filter = q.Filter
	}
	if q.SortBy != "" {
		s.SortBy = q.SortBy
	}
	if q.SortOrder != "" {
		s.SortOrder = q.SortOrder
	}
	if q.User != "" {
		s.User = q.User
	}
	if q.Args.Query != "" {
		s.Args.Query = q.Args.Query
	}
	if term != "" {
		s.Args.Query = q.Args.Query
	}
	if q.Page != 0 {
		s.Page = q.Page
	}

	if q.Latest {
		s = &search.Search{
			SortBy:    "id",
			SortOrder: "desc",
			Page:      1,
			Category:  "subs",
		}
	}

	if options.Proxy != "" {
		s.ProxyURL = options.Proxy
	}

	// make the request for results to nyaa.si
	m, err := s.Query()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func DownloadMultiple(cl *libtorrent.Client) ([]*torrent.Torrent, error) {
	q := &Search{
		SortBy:    download.SortBy,
		SortOrder: download.SortOrder,
		User:      download.User,
		Filter:    download.Filter,
		Page:      download.Page,
		Latest:    download.Latest,
		Category:  download.Category,
	}

	q.Args.Query = download.Query

	m, err := SearchAnime(q, download.Query)
	if err != nil {
		return nil, err
	}

	choices, err := fzfMenuMulti(m.Media)
	if err != nil {
		return nil, err
	}

	var torrents []*torrent.Torrent
	for _, item := range choices {
		t, err := cl.AddTorrent(item.Magnet)
		if err != nil {
			log.Println(err)
		}

		t.DownloadAll()
		torrents = append(torrents, t)
	}

	return torrents, nil
}

func CreateOutput(dir string) error {
	_, err := os.Stat(dir)
	if err == nil {
		return err
	} else {
		return os.MkdirAll(dir, 0o755)
	}
}
