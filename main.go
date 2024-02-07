package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	Magnet  string `short:"m" long:"magnet" description:"download the complete torrent from a magnet link"`
	Verbose bool   `short:"v" long:"verbose" description:"print debugging information and verbose output"`
}

var Debug = func(string, ...interface{}) {}

func downloadMagnet(magnet string) error {
	cfg := torrent.NewDefaultClientConfig()
	cfg.Debug = opts.Verbose

	cl, err := torrent.NewClient(cfg)
	if err != nil {
		return err
	}
	defer cl.Close()

	tor, err := cl.AddMagnet(magnet)
	if err != nil {
		return err
	}

	<-tor.GotInfo()
	tor.DownloadAll()
	cl.WaitAll()
	return nil
}

func streamMagnet(magnet string) error {
	cfg := torrent.NewDefaultClientConfig()
	cfg.Debug = opts.Verbose

	cl, err := torrent.NewClient(cfg)
	if err != nil {
		return err
	}
	defer cl.Close()

	tor, err := cl.AddMagnet(magnet)
	if err != nil {
		return err
	}

	// block until we get the torrent info
	<-tor.GotInfo()

	name := tor.Info().BestName()
	infoHash := tor.InfoHash().HexString()
	creation := tor.Metainfo().CreationDate

	Debug("%s - hash [%v] creation [%v]\n", name, infoHash, creation)

	for _, f := range tor.Files() {
		Debug("path %v\n", magnet, f.Path())
		ext := path.Ext(f.Path())

		// janky way to skip non-video files and to try and guess the correct file. It should usually be a single file
		switch ext {
		case ".mp4", ".mkv", ".avi", ".avif", ".av1", ".mov", ".flv", ".f4v", ".webm", ".wmv", ".mpeg", ".mpg", ".mlv", ".hevc", ".flac", ".flic":
		default:
			continue
		}

		done := make(chan bool)

		go func() {
			// copy video to filesystem
			tr := f.NewReader()
			defer tr.Close()
			ch := done
			var description bytes.Buffer
			if _, err := io.Copy(&description, tr); err != nil {
				panic(err)
			}
			ch <- true
		}()

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "video/mp4")
			http.ServeContent(w, r, f.DisplayPath(), time.Unix(f.Torrent().Metainfo().CreationDate, 0), f.NewReader())
		})

		go func() {
			if err := http.ListenAndServe(":8080", nil); err != nil {
				log.Fatal(err)
			}
		}()

		// print the link to the video
		fmt.Printf("'http://localhost:8080/%s'\n", f.DisplayPath())

		// block until at least the video is done
		<-done
	}

	return nil
}

func Torrenter(args []string) error {
	err := streamMagnet(opts.Magnet)
	return err
}

func main() {
	args, err := flags.Parse(&opts)
	if flags.WroteHelp(err) {
		os.Exit(0)
	}
	if err != nil {
		log.Fatal(err)
	}

	if opts.Verbose {
		Debug = log.Printf
	}

	if err := Torrenter(args); err != nil {
		log.Fatal(err)
	}
}
