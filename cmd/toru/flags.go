package main

import "github.com/jessevdk/go-flags"

// Global application options
type Options struct {
	// verbosity with level
	Verbose []bool `short:"v" long:"verbose" description:"Verbose output"`
	Version bool   `short:"V" long:"version" description:"display version info and exit"`
	Player  string `short:"p" long:"player" description:"set a custom video player. Use {url} as a placeholder if the url is not the last argument in the string"`
	Port    string `short:"P" long:"port" description:"set the port that torrents are streamed over"`
}

// Streaming options
type Stream struct {
	Magnet      string         `short:"m" long:"magnet" description:"stream directly from the provided torrent magnet link"`
	TorrentFile string         `short:"t" long:"torrent" description:"stream directly from the provided torrent file or torrent URL"`
	Remove      bool           `long:"rm" description:"remove cached files after exiting"`
	Latest      bool           `short:"l" long:"latest" description:"view the latest anime and select an episode"`
	FromJson    flags.Filename `short:"j" long:"from-json" description:"resume selection from prior search saved as json [see: toru search --help]"`

	// optional magnet link or torrent file as a trailing argument instead of explicitly defined
	Args struct {
		Query string
	} `positional-args:"yes"`
}

// Downloading options
type Download struct {
	Directory string `short:"d" long:"dir" description:"parent directory to download torrents to"`

	// magnet link, torrent or torrent file url
	Args struct {
		Query string
	} `positional-args:"yes" required:"yes" description:"magnet link, torrent file, or torrent file URL"`
}

// Non-interactive CLI search options
type Search struct {
	SortBy      string `short:"b" long:"sort-by" description:"sort results by a category [size|date|seeders|leechers|downloads]"`
	SortOrder   string `short:"o" long:"sort-order" description:"sort by ascending or descending: options [asc|desc]" choice:"asc" choice:"desc"`
	User        string `short:"u" long:"user" description:"search for content by a user"`
	Filter      string `short:"f" long:"filter" description:"filter content. Options: [no-remakes|trusted]"`
	Proxy       string `short:"x" long:"proxy" description:"use a proxy URL like nyaa.iss.ink"`
	Page        uint   `short:"p" long:"page" description:"which results page to display [default 1]"`
	Stream      bool   `short:"s" long:"stream" description:"stream selected torrents after search"`
	Download    bool   `short:"d" long:"download" description:"download selected torrents after search"`
	Multi       bool   `short:"m" long:"multi" description:"choose multiple torrents to queue for downloading or streaming"`
	Latest      bool   `short:"n" long:"latest" description:"view the latest anime"`
	Category    string `short:"c" long:"category" description:"f"`
	List        bool   `short:"l" long:"list" description:"list all accepted arguments for searching by categories"`
	Json        bool   `short:"j" long:"json" description:"output search results as Json"`
	Stdout      bool   `short:"P" long:"print" description:"output search results in a pretty and readable format to stdout"`
	Interactive bool   `short:"i" long:"fzf" description:"interact with the search results using fzf"`

	// actual search query, doesn't matter where it is placed and is OPTIONAL
	Args struct {
		Query string
	} `positional-args:"yes"`
}
