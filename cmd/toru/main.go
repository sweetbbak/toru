package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/jessevdk/go-flags"
	"github.com/sweetbbak/toru/pkg/libtorrent"
)

const (
	binaryName = "toru"
	version    = "v0.1"
)

var (
	options    Options
	streaming  Stream
	searchopts Search
	download   Download
)

var parser = flags.NewParser(&options, flags.Default)

func init() {
	s, err := parser.AddCommand("stream", "stream torrents", "stream torrents", &streaming)
	if err != nil {
		log.Fatal(err)
	}
	s1, err := parser.AddCommand("search", "search torrents and output them in a given format", "search Nyaa.si for content", &searchopts)
	if err != nil {
		log.Fatal(err)
	}
	r, err := parser.AddCommand("run", "run an interactive terminal session", "run an interactive terminal session with toru", &options)
	if err != nil {
		log.Fatal(err)
	}
	d, err := parser.AddCommand("download", "download torrents", "download torrent from .torrent file, magnet or URL to a .torrent file", &download)
	if err != nil {
		log.Fatal(err)
	}
	_, err = parser.AddCommand("version", "print version and debugging info", "print version and debugging info", &options)
	if err != nil {
		log.Fatal(err)
	}

	s.Aliases = []string{"s", "play"}
	s1.Aliases = []string{"se", "q"}
	r.Aliases = []string{"r"}
	d.Aliases = []string{"dl", "d"}
	options.Port = "8080"
}

func main() {
	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
		default:
			os.Exit(1)
		}
	}

	// TODO: func add config parsing here

	cl := libtorrent.NewClient(binaryName, options.Port)
	cl.DisableIPV6 = options.DisableIPV6

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
		if err := DownloadTorrent(cl); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	case "run", "interactive":
		if err := InteractiveSearch(cl); err != nil {
			log.Fatal(err)
		}
	case "version":
		fmt.Printf("%s %s %s/%s\n", binaryName, version, runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}
}
