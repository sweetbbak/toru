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
	version    = "v0.3.2"
)

var (
	options     Options
	runner      Run
	streaming   Stream
	searchopts  Search
	download    Download
	completions Completions
	wget        WGET
	latest      Latest
)

var parser = flags.NewParser(&options, flags.Default)
var osargs []string

func init() {
	s, err := parser.AddCommand("stream", "stream torrents", "stream torrents", &streaming)
	if err != nil {
		log.Fatal(err)
	}
	s1, err := parser.AddCommand("search", "search torrents and output them in a given format", "search Nyaa.si for content", &searchopts)
	if err != nil {
		log.Fatal(err)
	}
	r, err := parser.AddCommand("run", "run an interactive terminal session", "run an interactive terminal session with toru", &runner)
	if err != nil {
		log.Fatal(err)
	}
	d, err := parser.AddCommand("download", "select one or many torrents to download", "download torrent from .torrent file, magnet or URL to a .torrent file", &download)
	if err != nil {
		log.Fatal(err)

	}
	_, err = parser.AddCommand("wget", "wget a torrent file", "wget a torrent file", &wget)
	if err != nil {
		log.Fatal(err)
	}
	_, err = parser.AddCommand("latest", "get the latest anime", "get the latest anime from nyaa.si", &latest)
	if err != nil {
		log.Fatal(err)
	}
	_, err = parser.AddCommand("version", "print version and debugging info", "print version and debugging info", &options)
	if err != nil {
		log.Fatal(err)
	}
	_, err = parser.AddCommand("init", "source zsh or bash completions", "", &completions)
	if err != nil {
		log.Fatal(err)
	}

	s.Aliases = []string{"s", "play"}
	s1.Aliases = []string{"se", "q"}
	r.Aliases = []string{"r"}
	d.Aliases = []string{"dl", "d"}

	// port for server *NOT TORRENT PORT*
	options.Port = libtorrent.INTERNAL_STREAM_PORT // port 8888 or the next open port available
	options.TorrentPort = -1                       // set to -1 -> or user input here -> or set by libtorrent on backend to any open port

	if len(os.Args) == 1 {
		osargs = append(osargs, "run")
	} else {
		osargs = os.Args[1:]
	}
}

func main() {
	args, err := parser.ParseArgs(osargs)
	if err != nil {
		if flags.WroteHelp(err) {
			os.Exit(0)
		} else {
			log.Fatal(err)
		}
	}

	var optArg string // optional positional search argument
	if len(args) > 0 {
		optArg = args[0]
	} else {
		optArg = ""
	}

	// dump shell completions to be sourced by shell init files
	// keep this high up to avoid uneccesary hangups
	if parser.Active.Name == "init" {
		if err := Completers(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	cl := libtorrent.NewClient(binaryName, options.Port)
	cl.DisableIPV6 = options.DisableIPV6

	if options.TorrentPort != -1 {
		cl.TorrentPort = options.TorrentPort
	}

	if parser.Active.Name == "download" {
		if err := DownloadMain(cl); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

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
	case "wget":
		if err := DownloadTorrent(cl); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	case "run", "interactive":
		if err := InteractiveSearch(cl, optArg); err != nil {
			log.Fatal(err)
		}
	case "latest":
		searchopts.Latest = true
		searchopts.Interactive = true
		if err := runSearch(cl); err != nil {
			log.Fatal(err)
		}
	case "version":
		fmt.Printf("%s %s %s/%s\n", binaryName, version, runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}
}
