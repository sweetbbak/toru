package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/anacrolix/torrent"
	"github.com/charmbracelet/lipgloss"
	fzf "github.com/ktr0731/go-fuzzyfinder"
	"github.com/sweetbbak/toru/pkg/libtorrent"
	"github.com/sweetbbak/toru/pkg/nyaa"
	"github.com/sweetbbak/toru/pkg/search"
)

func runSearch(cl *libtorrent.Client) error {
	s := search.NewSearch()
	// list category options and return
	if searchopts.List {
		search.List()
		return nil
	}

	// build the query
	if searchopts.Category != "" {
		s.Category = searchopts.Category
	}
	if searchopts.Filter != "" {
		s.Filter = searchopts.Filter
	}
	if searchopts.SortBy != "" {
		s.SortBy = searchopts.SortBy
	}
	if searchopts.SortOrder != "" {
		s.SortOrder = searchopts.SortOrder
	}
	if searchopts.User != "" {
		s.User = searchopts.User
	}
	if searchopts.Args.Query != "" {
		s.Args.Query = searchopts.Args.Query
	}
	if searchopts.Page != 0 {
		s.Page = searchopts.Page
	}

	if searchopts.Latest {
		// search.LatestAnime(searchopts.Args.Query, 1)
		s = &search.Search{
			SortBy:    "id",
			SortOrder: "desc",
			Page:      1,
			Category:  "subs",
		}
	}

	if options.Proxy != "" {
		s.ProxyURL = options.Proxy
	}

	// make the request for results to nyaa.si
	m, err := s.Query()
	if err != nil {
		return err
	}

	// print and/or handle results
	// TODO: change os.Stdout to a generic io.Writer and support a file directly
	if searchopts.Json {
		err := m.WriteToJson(os.Stdout)
		if err != nil {
			return err
		}
		return nil
	}

	if searchopts.Stdout {
		m.PrintResults()
	}

	if searchopts.Stream || searchopts.Interactive {
		return SelectAndPlay(cl, m)
	}

	if !searchopts.Stdout && !searchopts.Json && !searchopts.Stream {
		m.PrintResults()
	}

	return nil
}

// basic search and play
func InteractiveSearch(cl *libtorrent.Client, query string) error {
	header := cutePrint.Align(lipgloss.Center).Render("Toru, stream anime, no strings attached")
	fmt.Println(header)

	var s string
	var err error
	if query == "" {
		s, err = Prompt("Search for an anime: ")
		if err != nil {
			return err
		}
	} else {
		s = query
	}

	m, err := search.LatestAnime(s, options.Proxy, 1)
	if err != nil {
		return err
	}

	cj := path.Join(cl.DataDir, "cache.json")
	err = m.Cache(cj)
	if err != nil {
		log.Println(err)
	}

LOOP:
	choice, err := fzfMenu(m.Media)
	if err != nil {
		return err
	}

	err = PlayTorrent(cl, choice.Magnet)
	if err != nil {
		return err
	}

	action, err := fzfMain()
	if err != nil {
		return err
	}

	switch action {
	case "select":
		goto LOOP
	case "search":
		if err := InteractiveSearch(cl, ""); err != nil {
			log.Fatal(err)
		}
	case "exit":
		os.Exit(0)
	}

	return nil
}

func SelectAndPlay(cl *libtorrent.Client, res *search.Results) error {
LOOP:
	// select a result
	choice, err := fzfMenu(res.Media)
	if err != nil {
		return err
	}

	// play that result
	err = PlayTorrent(cl, choice.Magnet)
	if err != nil {
		return err
	}

	// select, search or exit
	action, err := fzfMain()
	if err != nil {
		return err
	}

	switch action {
	case "select":
		goto LOOP
	case "search":
		if err := InteractiveSearch(cl, ""); err != nil {
			log.Fatal(err)
		}
	case "exit":
		os.Exit(0)
	}

	return nil

}

func fzfMain() (string, error) {
	m := []string{"select", "search", "exit"}
	idx, err := fzf.Find(
		m,
		func(i int) string {
			return m[i]
		},
	)

	// User has selected nothing
	if err != nil {
		if errors.Is(err, fzf.ErrAbort) {
			return "exit", err
		} else {
			return "exit", err
		}
	}

	return m[idx], nil
}

func fzfEpisodes(m []*torrent.File) (string, error) {
	l := len(m) - 1
	idx, err := fzf.Find(
		m,
		func(i int) string {
			// fpath := m[i].DisplayPath()
			fpath := m[l-i].DisplayPath()
			return fpath
		},
	)

	// User has selected nothing
	if err != nil {
		if errors.Is(err, fzf.ErrAbort) {
			os.Exit(0)
		} else {
			return "", err
		}
	}

	return m[idx].DisplayPath(), nil
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

			return FormatMedia(m[i])

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

func fzfMenuMulti(m []nyaa.Media) ([]nyaa.Media, error) {
	idxs, err := fzf.FindMulti(
		m,
		func(i int) string {
			return m[i].Name
		},
		fzf.WithPreviewWindow(func(i, width, height int) string {
			if i == -1 {
				return "lol"
			}

			return FormatMedia(m[i])

		}),
	)

	var matches []nyaa.Media
	for _, item := range idxs {
		matches = append(matches, m[item])
	}

	// User has selected nothing
	if err != nil {
		if errors.Is(err, fzf.ErrAbort) {
			os.Exit(0)
		} else {
			return nil, err
		}
	}

	return matches, nil
}
