package main

import (
	"log"
	"os"
	"toru/pkg/libtorrent"

	"github.com/jessevdk/go-flags"
)

var (
	options    Options
	streaming  Stream
	searchopts Search
	download   Download
)

var parser = flags.NewParser(&options, flags.Default)

func init() {
	parser.AddCommand("stream", "stream torrents", "stream torrents", &streaming)
	parser.AddCommand("search", "search torrents and output them in a given format", "search Nyaa.si for content", &searchopts)
	parser.AddCommand("run", "run an interactive terminal session", "run an interactive terminal session with toru", &options)
	d, _ := parser.AddCommand("download", "download torrents", "download torrent from .torrent file, magnet or URL to a .torrent file", &download)
	d.Aliases = []string{"dl", "d"}
}

func main() {
	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}

	// TODO: func add config parsing here

	cl := libtorrent.NewClient("tori", 8080)
	if err := cl.Init(); err != nil {
		log.Fatal(err)
	}

	switch parser.Active.Name {
	case "search":
		if err := runSearch(cl); err != nil {
			log.Fatal(err)
		}
	case "stream":
		if err := runStream(cl); err != nil {
			log.Fatal(err)
		}
	case "dl", "download":
		if err := runSearch(cl); err != nil {
			log.Fatal(err)
		}
	case "run", "interactive":
		if err := InteractiveSearch(cl); err != nil {
			log.Fatal(err)
		}
	}
}
