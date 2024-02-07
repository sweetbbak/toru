package main

import (
	// "bytes"
	"errors"
	"fmt"
	"syscall"

	// "io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
	"toru/pkg/nyaa"

	"github.com/anacrolix/torrent"
	"github.com/jessevdk/go-flags"
	fzf "github.com/ktr0731/go-fuzzyfinder"
	// "github.com/phayes/freeport"
)

var opts struct {
	Magnet   string `short:"m" long:"magnet" description:"download or stream the complete torrent from a magnet link"`
	Search   string `short:"s" long:"search" description:"search for an anime title"`
	Stream   bool   `short:"S" long:"stream" description:"automaticaly stream the torrent using a video player"`
	Download bool   `short:"d" long:"download" command:"search" description:"search for an anime title"`
	Latest   bool   `short:"l" long:"latest" command:"view latest anime interactively"`
	Player   string `short:"P" long:"player" description:"video player to use for playing torrents. Pass a string containing the player name and options"`
	Page     uint   `short:"p" long:"page" description:"get search results on page N of nyaa.si [default 1]"`
	Port     int    `long:"port" description:"default 8080"`
	Verbose  bool   `short:"v" long:"verbose" description:"print debugging information and verbose output"`
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

func streamMagnet(magnet string, usePlayer bool) error {
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

		// go func() {
		// 	// copy video to filesystem
		// 	tr := f.NewReader()
		// 	defer tr.Close()

		// 	ch := done
		// 	var description bytes.Buffer

		// 	if _, err := io.Copy(&description, tr); err != nil {
		// 		panic(err)
		// 	}
		// 	ch <- true
		// }()

		fname := f.DisplayPath()
		fname = escapeUrl(fname)

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "video/mp4")
			http.ServeContent(w, r, fname, time.Unix(f.Torrent().Metainfo().CreationDate, 0), f.NewReader())
		})

		// p, err := freeport.GetFreePort()
		// if err != nil {
		// 	return err
		// }

		port := fmt.Sprintf(":%d", 8080)

		go func() {
			if err := http.ListenAndServe(port, nil); err != nil {
				log.Fatal(err)
			}
		}()

		// print the link to the video
		link := fmt.Sprintf("http://localhost%s/%s\n", port, fname)
		fmt.Println(link)

		if usePlayer {
			timeout := time.Duration(time.Second * 10)
			client := http.Client{
				Timeout: timeout,
			}
			_, err := client.Get(link)

			player := opts.Player

			cmd := exec.Command(player, link)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			cmd.SysProcAttr = &syscall.SysProcAttr{
				Setsid: true,
			}

			err = cmd.Run()
			if err != nil {
				log.Println(err)
			}
		}

		// block until at least the video is done
		<-done
	}

	return nil
}

func escapeUrl(u string) string {
	u = strings.ReplaceAll(u, "'", "")
	u = strings.ReplaceAll(u, "\n", "")
	u = strings.ReplaceAll(u, " ", "_")
	u = strings.ReplaceAll(u, "_-_", "_")
	u = strings.ReplaceAll(u, "__", "_")
	u = strings.ReplaceAll(u, "--", "-")
	return u
}

// TODO: make use of all search parameters and expose them to user
func search(q string, page uint) ([]nyaa.Media, error) {
	m, err := nyaa.Search(q, nyaa.SearchParameters{
		Category:  nyaa.CategoryAnimeEnglishTranslated,
		SortBy:    nyaa.SortByDate,
		SortOrder: nyaa.SortOrderDescending,
		Page:      page,
	})
	if err != nil {
		return nil, err
	}

	return m, nil
}

func fzfMenu(m []nyaa.Media) (nyaa.Media, error) {
	idx, err := fzf.Find(
		m,
		func(i int) string {
			return m[i].Name
		},
		fzf.WithPreviewWindow(func(i, width, height int) string {
			if i == -1 {
				return "lol"
			}
			return fmt.Sprintf("%s\n%s\nDate - %s\n%s\nDownloads %d\n[\x1b[32m%v\x1b[0m|\x1b[31m%v\x1b[0m]\nSubmitted by - %v\nSize - %v",
				m[i].Name,
				WrapString(m[i].Description, 55),
				m[i].Date,
				m[i].Category,
				m[i].Downloads,
				m[i].Seeders,
				m[i].Leechers,
				m[i].Submitter,
				m[i].Size,
			)
		}),
	)
	// User has selected nothing
	if err != nil {
		if errors.Is(err, fzf.ErrAbort) {
			os.Exit(0)
		} else {
			return nyaa.Media{}, err
		}
	}

	return m[idx], nil
}

func Torrenter(args []string) error {
	if opts.Magnet != "" {
		err := streamMagnet(opts.Magnet, true)
		if err != nil {
			return err
		}
	}

	var media []nyaa.Media
	var err error
	if opts.Search != "" {
		media, err = search(opts.Search, opts.Page)
		if err != nil {
			return err
		}
	}

	if opts.Latest {
		media, err = search("", 1)
		if err != nil {
			return err
		}
	}

	if len(media) < 1 {
		return fmt.Errorf("No results")
	}

	sel, err := fzfMenu(media)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n%s\nDate - %s\n%s\nDownloads %d\n[\x1b[32m%v\x1b[0m|\x1b[31m%v\x1b[0m]\nSubmitted by - %v\nSize - %v\n%v\n",
		sel.Name,
		WrapString(sel.Description, 55),
		sel.Date,
		sel.Category,
		sel.Downloads,
		sel.Seeders,
		sel.Leechers,
		sel.Submitter,
		sel.Size,
		sel.Magnet,
	)

	if err := streamMagnet(sel.Magnet, true); err != nil {
		return err
	}

	return nil
}

func init() {
	opts.Page = 1
	opts.Player = "mpv"
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
