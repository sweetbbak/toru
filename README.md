<p align="center">
  <img src="assets/toru.png" />
<br>
<a href="http://makeapullrequest.com"><img src="https://img.shields.io/badge/PRs-welcome-brightgreen.svg"></a>
<a href="#Linux"><img src="https://img.shields.io/badge/os-linux-brightgreen">
<a href="#MacOS"><img src="https://img.shields.io/badge/os-mac-brightgreen">
<a href="#Android"><img src="https://img.shields.io/badge/os-android-brightgreen">
<a href="#Windows"><img src="https://img.shields.io/badge/os-windows-yellowgreen">
<a href="#iOS"><img src="https://img.shields.io/badge/os-ios-yellow">
<a href="#Steam-deck"><img src="https://img.shields.io/badge/os-steamdeck-yellow">
<br>
<a href="https://www.buymeacoffee.com/sweetbabyalaska"><img src="https://img.shields.io/badge/Buy%20Me%20a%20Coffee-ffdd00?style=for-the-badge&logo=buy-me-a-coffee&logoColor=black"></a>
<a href="https://github.com/sweetbbak"><img src="https://img.shields.io/badge/creator-sweet-green"></a>
<br>
</p>

<p align="center">
<a href="#golang"><img src="https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white">
<a href="go"><img src="https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black">
<a href="linux"><img src="https://img.shields.io/badge/Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white">
<a href="bsd"><img src="https://img.shields.io/badge/-OpenBSD-%23FCC771?style=for-the-badge&logo=openbsd&logoColor=black">
<a href="mac"><img src="https://img.shields.io/badge/mac%20os-000000?style=for-the-badge&logo=macos&logoColor=F0F0F0">
</p>

<h3 align="center">
Toru: Stream Torrents instantly, straight from the Command Line
</h3>

Toru, the command line tool designed for anime enthusiasts. Browse the latest releases on nyaa.si or search for specific titles with Toru's intuitive interface. Stream episodes directly from the command line using MPV or your favorite video player or browser.
Toru serves the selected anime over localhost, making it as accessible as any HTTP link. It also functions as a convenient torrent client for magnet links. Simplify your streaming and torrenting with Toru.

![example of toru in progress](assets/example.gif)

![example of toru in progress](assets/search.png)

## Table of Contents

- [Install](#install)
- [Example Usage](#examples)
- [Contribution Guidelines](./CONTRIBUTING.md)
- [Disclaimer](./disclaimer.md)

## Install

Quick install a pre-built binary

```sh
export PREFIX="$HOME/bin"
wget "https://github.com/sweetbbak/toru/releases/download/v0.1/toru_$(uname -s)_$(uname -m).tar.gz" -O - | tar -xz
mv toru "${PREFIX}"
```

on Windows

```sh
iwr -Uri "https://github.com/sweetbbak/toru/releases/download/v0.1/toru_Windows_x86_64.zip" -OutFile toru_Windows_x86_64.zip
```

<details closed>
  <summary>Install Go</summary>
  <a href="https://go.dev/doc/install">Install go</a>
  This project requires go 1.21.7 or higher.
</details>

```sh
go install github.com/sweetbbak/toru/cmd/toru@latest
```

<details closed>
  <summary>Build from source</summary>

using Go is recommended alongside `gup` so that it can be easily updated

```sh
git clone https://github.com/sweetbbak/toru.git && cd toru
go build -o toru ./cmd/toru
```

you can also use the justfile

```sh
git clone https://github.com/sweetbbak/toru.git && cd toru
just
```

or the makefile

```sh
git clone https://github.com/sweetbbak/toru.git && cd toru
make build
```

Using `docker`, `podman` with the `shell.nix` file on non-nixOS distros
this will automatically pull the correct version of Go and install `just`
so that building is easy.

```sh
git clone https://github.com/sweetbbak/toru.git && cd toru
# mount the project directory inside the container
podman run --volume $(pwd):/toru -ti docker.io/nixos/nix:latest
# inside the container run:
cd /toru
nix-shell
just
```

### Building for different platforms and architectures

Run to find your target architecture and platform:

```sh
go tool dist list
```

then use the environment variables `GOOS` and `GOARCH` before using
the build command.

Example:

```sh
GOOS=linux GOARCH=arm64 go build -o toru ./cmd/toru
```

_If you do this_ feel free to create an issue for a platform reporting how it went
So far there is an issue with android and termux as well as arm 6.

</details>

if you are on nix or have nix installed you can just use the shell.nix directly and run `just` or `make` or use `go build -o toru ./cmd/toru`.

## Changelog

- added `--proxy` which allows use of nyaa proxy sites and sukebi
- added the ability to disable ipv6
- sub-command "run" now accepts a trailing search term argument

## Examples

Search for an anime:

```sh
toru search -i
toru search ""
```

![example of toru in progress](assets/search.png)

the selected torrent will begin playing and once the video player is closed
you will have the option to select another episode, make another search query,
or to exit.

#### _View_ the latest anime on nyaa.si in an interactive fzf session

```sh
toru search --latest
```

#### Search for a specific keyword or series

```sh
toru search "Akuyaku"
```

If you know the magnet link for the content you can directly download or stream it

```sh
toru stream --magnet 'magnet:?xt=urn:btih:1...announce'
toru download --magnet 'magnet:?xt=urn:btih:1...announce'
```

All of the above outputs a link that you can use to stream the torrent `'http://localhost:8080/stream?ep=torrent_info_hash'`
you can treat this link like any other http link and stream it with `mpv` or `vlc`, download it, use `yt-dlp`, or open it in the browser etc...

### Handling input and output

You can use toru to search for anime and other media types and then output the results in multiple formats.
Output in Json and parsing that output with `jq`:

```sh
toru search --latest --json | jq -r '.[]|.Name,.Magnet'
```

#### Open a cached search session from a json file

```sh
toru search --json "one piece" > cache.json
toru search --from-json cacne.json --interactive
```

#### Output in a human readable format:

```sh
toru search "akuyaku 99 1080p"
# Outputs:
[Erai-raws] Akuyaku Reijou Level 99 - 01 [1080p][Multiple Subtitle] [ENG][POR-BR][SPA-LA][SPA][ARA][FRE][GER][ITA][RUS]
2024-01-09 07:36:29
Downloads: 1203
[33|0]
Size: 727 MB
magnet:... [ magnet link here ]

```

### Creating your own CLI tool using toru

```sh
# Create a JSON file using toru
toru search --latest --json > out.json
# Here is a simple example of using fzf and toru to create a simple interace to select and play torrents
# you can also replace toru with any CLI bittorrent client
cat out.json | jq '.[].Name' | \
fzf --preview='cat out.json | jq -r ".[{n}]"' \
  --bind "enter:execute(cat out.json | jq -r '.[{n}].Magnet')+abort" | \
  xargs toru stream --magnet
```

> [!IMPORTANT]\
> toru is in a very early development phase! In order to provide a consistent and smooth experience
> the CLI interface is subject to change. PR's and advice on project sturcture, pkg organization and
> feedback on the UI of toru is much appreciated.
>
> Currently tested on Linux and Windows
> I would much appreciate someone reporting on the functionality on any BSD or Mac
> Android with termux currently has UDP issues. Idk much about how android works

## Features

- [x] Stream anime from torrents
- [x] add Nyaa.si as a source
- [ ] add a generic torrent tracker library for Korean and American movies
- [ ] package as various formats (AUR, DEB, Flatpak, AppImage, Release binaries)
- [ ] ensure compatibility across platforms and aim for consistent compatibility (should work but currently untested)

## Good issues to work on

- [ ] add a package for Arch, Fedora, Nix, Scoop or otherwise
- [ ] get toru working on Android. (currently an issue with connecting to peers via UDP)

## Roadmap

- Daemonize into the background and listen for commands on a socket (optional for user, sometimes this is annoying)
- Simple torrent client features (download|seed|add magnet|stream|search)
- Look into file and search caching
- Add other trackers besides `nyaa.si`
- Expand user interface with bubbletea
- Ensure we are not straining or leeching off of the network more than we are giving

## Contributing

PR's welcome! This project currently uses Golang 1.21.7 along with standard go formatting using `gopls`
TODO: add a development containerfile and automate building binaries for all platforms

## Why though?

Because scraping is annoying af and it constantly breaks. On top of that, _someone_ is paying for those servers.
Torrents are more resistant to takedowns and hopefully will have more longevity.

## Credits

torrenting library:
![anacrolix/torrent](https://github.com/anacrolix/torrent)

Nyaa package is modified from here:
![Feyko/nyaa-go](https://github.com/Feyko/nyaa-go)

## Support

Consider creating a PR, taking up a minor issue on the TODO list, leaving an issue to help improve functionality or buy
me a coffee!

![moe-visitor-counter](https://count.getloli.com/get/@sweetbbak?theme=asoul)
