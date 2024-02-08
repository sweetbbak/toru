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
<a href="https://github.com/sweetbbak"><img src="https://img.shields.io/badge/creator-sweet-green"></a>
<br>
</p>

<h3 align="center">
A CLI tool to browse and stream anime with the power of web-torrents. Ani-cli meets Miru.
</h3>


`toru` allows you to stream torrents from the command line. You can view the latest anime on nyaa.si
or provide a search query and then pick an anime from a fzf-like interface, and then stream that episode
directly from the command line in MPV or your favorite video player (including the browser).

`toru` will serve the selected anime over port `8080` by default on `localhost` and it can be treated
like any other http link. `toru` can also be used as a makeshift torrent client for downloading magnet links

<p align="center">
  <img src="assets/example.webm" />
</p>

## Table of Contents

- [Install](#install)
- [Example Usage](#examples)
- [Contribution Guidelines](./CONTRIBUTING.md)
- [Disclaimer](./disclaimer.md)

## Install

```sh
go install github.com/sweetbbak/toru@latest
# OR
git clone https://github.com/sweetbbak/toru.git && cd toru
go build
```

## Examples

```sh
# Directly stream from magnet links
toru --magnet 'magnet:?xt=urn:btih:1a4fe542f61743b794272e1acdd3878b1fa73c5a&dn=%5BSubsPlease%5D%20Akuyaku%20Reijou%20Level%2099%20-%2005%20%28480p%29%20%5B0D52BF4C%5D.mkv&tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce'
# outputs a link that you can use to stream the torrent
# 'http://localhost:8080/Video content.mkv'

# View the latest on nyaa.si in an interactive terminal session
toru --latest

# Search for a series
toru --search "Akuyaku"
```

## Features

## Roadmap
