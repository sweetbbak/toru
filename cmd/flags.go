package main

// Global application options
type Options struct {
	// verbosity with level
	Verbose []bool `short:"v" long:"verbose" description:"Verbose output"`
}

// Streaming options
type Stream struct {
	Player      string `short:"p" long:"player" description:"which player to use along with its options and flags ie: 'mpv --vo=kitty' "`
	Magnet      string `short:"m" long:"magnet" description:"stream torrent directly from the provided magnet link"`
	TorrentFile string `short:"t" long:"torrent" description:"stream torrent directly from the provided torrent file"`
	Remove      bool   `long:"rm" description:"remove cached files after exiting"`
	Latest      bool   `short:"l" long:"latest" description:"view the latest anime on nyaa.si"`

	// optional magnet link or torrent file
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

// type Category struct {
// 	Anime           bool `short:"a" long:"anime" description:"sort by all anime [default]"`
// 	AnimeMusicVideo bool `long:"music-video" description:"sort by music videos"`
// 	AnimeEng        bool `long:"subs" description:"sort by english translated anime"`
// 	AnimeNonEng     bool `long:"non-english" description:"sort by non english translated anime"`
// 	AnimeRaw        bool `long:"raw" description:"sort by raw anime"`
// 	All             bool `short:"A" long:"all" description:"sort by all categories"`
// 	Audio           bool `long:"audio" description:"sort by all audio"`
// 	AudioLossless   bool `long:"audio-lossless" description:"sort by lossless audio"`
// 	Literature      bool `long:"literature" description:"sort by all literature"`
// 	LiveAction      bool `long:"live-action" description:"sort by live action anime"`
// 	Pictures        bool `long:"pictures" description:"sort by pictures"`
// 	Software        bool `long:"software" description:"sort by software"`
// }
