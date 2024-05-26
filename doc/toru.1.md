% toru(1) | General Commands Manual

# NAME

Toru - stream anime using torrents without waiting for downloads

# SYNOPSIS

`toru [options] [command] [command-options] [optional-query]`

# DESCRIPTION

toru allows you to stream torrents from the terminal. You must specify a base command and then can optionally add flags to modify the behavior.
By default toru stores torrents and data in XDG_CACHE_HOME. toru uses nyaa.si to provide search capabilities for anime.
You can also use toru to download torrents using their magnet/url/path or you can stream any torrent that you have the link to (not just anime)

# OPTIONS

`-h, --help`

: Show the help message and exit

`-V, --version`

: Show version info and exit

`-v, --verbose`

: Verbose output

`-p, --player PLAYER_CMD...`

: Set a special player for toru to use for playing torrent video files

`-P, --port PORT...`

: Specify the port that toru serves videos on.

# COMMANDS

## download

download torrents

download torrent from .torrent file, magnet or URL to a .torrent file

**Usage: toru \[OPTIONS\] download \[download-OPTIONS\]**

**Aliases: dl, d**

: ## download torrents

download torrent from .torrent file, magnet or URL to a .torrent file

**-d**, **\--dir**

: parent directory to download torrents to

## Help Options

**-h**, **\--help**

: Show this help message

## run

run an interactive terminal session

run an interactive terminal session with toru

**Usage: toru \[OPTIONS\] run \[run-OPTIONS\]**

**Aliases: , r**

: ## run an interactive terminal session

run an interactive terminal session with toru

**-v**, **\--verbose**

: Verbose output

**-V**, **\--version**

: display version info and exit

**-p**, **\--player**

: set a custom video player. Use {url} as a placeholder if the url is
not the last argument in the string

**-P**, **\--port**

: set the port that torrents are streamed over

## Help Options

**-h**, **\--help**

: Show this help message

## search

search torrents and output them in a given format

search Nyaa.si for content

**Usage: toru \[OPTIONS\] search \[search-OPTIONS\]**

**Aliases: se, q**

: ## search torrents and output them in a given format

search Nyaa.si for content

**-b**, **\--sort-by**

: sort results by a category
\[size\|date\|seeders\|leechers\|downloads\]

**-o**, **\--sort-order**

: sort by ascending or descending: options \[asc\|desc\]

**-u**, **\--user**

: search for content by a user

**-f**, **\--filter**

: filter content. Options: \[no-remakes\|trusted\]

**-p**, **\--page**

: which results page to display \[default 1\]

**-s**, **\--stream**

: stream selected torrents after search

**-d**, **\--download**

: download selected torrents after search

**-m**, **\--multi**

: choose multiple torrents to queue for downloading or streaming

**-n**, **\--latest**

: view the latest anime

**-c**, **\--category**

: f

**-l**, **\--list**

: list all accepted arguments for searching by categories

**-j**, **\--json**

: output search results as Json

**-P**, **\--print**

: output search results in a pretty and readable format to stdout

**-i**, **\--fzf**

: interact with the search results using fzf

## Help Options

**-h**, **\--help**

: Show this help message

## stream

stream torrents

stream torrents

**Usage: toru \[OPTIONS\] stream \[stream-OPTIONS\]**

**Aliases: s, play**

: ## stream torrents

stream torrents

**-m**, **\--magnet**

: stream directly from the provided torrent magnet link

**-t**, **\--torrent**

: stream directly from the provided torrent file or torrent URL

**\--rm**

: remove cached files after exiting

**-l**, **\--latest**

: view the latest anime and select an episode

**-j**, **\--from-json**

: resume selection from prior search saved as json \[see: toru search
\--help\]

## Help Options

**-h**, **\--help**

: Show this help message

## version

print version and debugging info

print version and debugging info

**Usage: toru \[OPTIONS\] version \[version-OPTIONS\]**

## print version and debugging info

print version and debugging info

**-v**, **\--verbose**

: Verbose output

**-V**, **\--version**

: display version info and exit

**-p**, **\--player**

: set a custom video player. Use {url} as a placeholder if the url is
not the last argument in the string

**-P**, **\--port**

: set the port that torrents are streamed over

## Help Options

**-h**, **\--help**

: Show this help message
