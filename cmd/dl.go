package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"
	"toru/pkg/dl"
	// "github.com/anacrolix/torrent/metainfo"
)

var storageDir string = getStorage()

const projectDir = "toru"

func getStorage() string {
	s, _ := os.UserCacheDir()
	p := path.Join(s, projectDir)

	_, err := os.Stat(p)
	if err != nil {
		// create storage directory
		err := os.MkdirAll(p, 0o755)
		if err != nil {
			log.Fatal(err)
		}
	}
	return p
}

func serveTorrent(t *torrent.Torrent, port int) {
	for _, f := range t.Files() {
		ext := path.Ext(f.Path())
		// janky way to skip non-video files and to try and guess the correct file. It should usually be a single file
		switch ext {
		case ".mp4", ".mkv", ".avi", ".avif", ".av1", ".mov", ".flv", ".f4v", ".webm", ".wmv", ".mpeg", ".mpg", ".mlv", ".hevc", ".flac", ".flic":
		default:
			continue
		}

		done := make(chan bool)
		fname := f.DisplayPath()
		fname = escapeUrl(fname)

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "video/mp4")
			http.ServeContent(w, r, fname, time.Unix(f.Torrent().Metainfo().CreationDate, 0), f.NewReader())
		})

		port := fmt.Sprintf(":%d", port)

		go func() {
			if err := http.ListenAndServe(port, nil); err != nil {
				log.Fatal(err)
			}
		}()

		// print the link to the video
		link := fmt.Sprintf("http://localhost%s/%s\n", port, fname)
		fmt.Println(link)

		// block until at least the video is done
		<-done
	}
}

// returns a *torrent from a torrent file
func torrentFromFile(filename string, cl *torrent.Client) (*torrent.Torrent, error) {
	tor, err := cl.AddTorrentFromFile(filename)
	if err != nil {
		return nil, err
	}

	<-tor.GotInfo()
	return tor, nil
}

// Downloads a torrent file form a URL and then returns a torrent
func torrentFromURL(url string, cl *torrent.Client) (*torrent.Torrent, error) {
	path, err := dl.Download(url)
	if err != nil {
		return nil, err
	}

	return torrentFromFile(path, cl)
}

// Get a new default torrent client
func NewClient() (*torrent.Client, error) {
	cfg := torrent.NewDefaultClientConfig()
	cfg.Debug = opts.Verbose

	// set torrent directory to be project directory
	cfg.DefaultStorage = storage.NewFileByInfoHash(storageDir)

	cl, err := torrent.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	defer cl.Close()
	return cl, nil
}
