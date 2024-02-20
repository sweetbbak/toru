package search

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/sweetbbak/toru/pkg/nyaa"
)

// map user input into categories for interfacing with the CLI
var catmap = map[string]nyaa.Category{
	"anime":          nyaa.CategoryAnime,
	"english":        nyaa.CategoryAnimeEnglishTranslated,
	"subs":           nyaa.CategoryAnimeEnglishTranslated,
	"non-english":    nyaa.CategoryAnimeNonEnglishTranslated,
	"subs-other":     nyaa.CategoryAnimeNonEnglishTranslated,
	"music-video":    nyaa.CategoryAnimeMusicVideo,
	"mv":             nyaa.CategoryAnimeMusicVideo,
	"all":            nyaa.CategoryAllCategories,
	"raw":            nyaa.CategoryAnimeRaw,
	"audio":          nyaa.CategoryAudio,
	"audio-lossless": nyaa.CategoryAudioLossless,
	"literature":     nyaa.CategoryLiterature,
	"novels":         nyaa.CategoryLiterature,
	"english-novels": nyaa.CategoryLiteratureEnglishTranslated,
	"pictures":       nyaa.CategoryPictures,
	"images":         nyaa.CategoryPictures,
	"software":       nyaa.CategorySoftware,
	"":               nyaa.CategoryAnimeEnglishTranslated,
}

type Search struct {
	SortBy    string
	SortOrder string
	User      string
	Filter    string
	List      bool
	Page      uint
	Stream    bool
	Download  bool
	Multi     bool
	Latest    bool
	Category  string
	ProxyURL  string

	Args struct {
		Query string
	}
}

type Results struct {
	Media []nyaa.Media
}

func GetSortBy(s string) nyaa.SortBy {
	switch s {
	case "size":
		return nyaa.SortBySize
	case "seeders":
		return nyaa.SortBySeeders
	case "leechers":
		return nyaa.SortByLeechers
	case "downloads":
		return nyaa.SortByDownloads
	case "date":
		return nyaa.SortByDate
	default:
		return nyaa.SortByNone
	}
}

func GetCategory(s string) (nyaa.Category, error) {
	val, ok := catmap[s]
	if !ok {
		return nyaa.CategoryAnime, fmt.Errorf("Unknown category: %v", s)
	}
	return val, nil
}

func NewSearch() *Search {
	return &Search{
		SortBy:    "date",
		SortOrder: "desc",
		Page:      1,
		Category:  "subs",
	}
}

// build a query from user input for Nyaa
func (search *Search) Query() (*Results, error) {
	s := nyaa.SearchParameters{}
	s.SortBy = GetSortBy(search.SortBy)

	if search.ProxyURL != "" {
		s.Proxy = search.ProxyURL
	}

	if search.Filter != "" {
		switch search.Filter {
		case "no-remakes":
			s.Filter = 1
		case "trusted":
			s.Filter = 2
		default:
			s.Filter = 0
		}
	}

	if search.Page != 0 {
		s.Page = search.Page
	} else {
		s.Page = 1
	}

	if search.User != "" {
		s.User = search.User
	}

	s.SortOrder = nyaa.SortOrderDescending

	if search.SortOrder != "" {
		switch search.SortOrder {
		case "asc":
			s.SortOrder = nyaa.SortOrderAscending
		case "desc":
			s.SortOrder = nyaa.SortOrderDescending
		default:
			s.SortOrder = nyaa.SortOrderNone
		}
	}

	cat, err := GetCategory(search.Category)
	if err != nil {
		return nil, err
	}

	s.Category = cat
	res := &Results{}

	// error parsing html or error getting nyaa page
	res.Media, err = nyaa.Search(search.Args.Query, s)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// get the latest english subbed anime
func LatestAnime(query, proxy string, page uint) (*Results, error) {
	p := nyaa.SearchParameters{
		Category:  nyaa.CategoryAnimeEnglishTranslated,
		SortBy:    nyaa.SortByDate,
		SortOrder: nyaa.SortOrderDescending,
		Page:      page,
	}

	if proxy != "" {
		p.Proxy = proxy
	}

	m, err := nyaa.Search(query, p)
	if err != nil {
		return nil, err
	}

	r := &Results{}
	r.Media = m

	return r, nil
}

func List() {
	fmt.Print("Parameter          Value\n")
	for k, v := range catmap {
		if k == "" || k == " " {
			fmt.Printf("%-18s %s\n", "[empty string]", v)
		} else {
			fmt.Printf("%-18s %s\n", k, v)
		}
	}
}

func FormatMedia(m nyaa.Media) string {
	return fmt.Sprintf("%s\n%s\nDownloads: %d\n[\x1b[32m%v\x1b[0m|\x1b[31m%v\x1b[0m]\nSize: %v\n%v\n%v\n%v\n%v\n",
		m.Name,
		m.Date.Format(time.DateTime),
		m.Downloads,
		m.Seeders,
		m.Leechers,
		humanize.Bytes(m.Size),
		m.Magnet,
		m.Hash,
		m.ID,
		m.Torrent,
	)
}

func (r *Results) PrintResults() {
	for _, m := range r.Media {
		fmt.Println(FormatMedia(m))
	}
}

func (r *Results) ToJson() (string, error) {
	b, err := json.Marshal(r.Media)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (r *Results) WriteToJson(w io.Writer) error {
	b, err := json.Marshal(r.Media)
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	return err
}

func (r *Results) Cache(fpath string) error {
	j, err := r.ToJson()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fpath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(j)
	return err
}

func (r *Results) ReadCache(fpath string) (*Results, error) {
	b, err := os.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	f := r
	err = json.Unmarshal(b, f)
	if err != nil {
		return nil, err
	}

	return f, nil
}

type CacheMedia struct {
	Category     string
	CommentCount float64
	Comments     interface{}
	Date         string
	Description  string
	Downloads    float64
	Files        interface{}
	Hash         string
	ID           float64
	Information  string
	IsFull       bool
	Leechers     float64
	Magnet       string
	Name         string
	Seeders      float64
	Size         float64
	Submitter    string
	Torrent      string
}

// Generic read cached results from a json file
func ReadCache(fpath string) ([]nyaa.Media, error) {
	m := []nyaa.Media{}

	b, err := os.ReadFile(fpath)
	if err != nil {
		return m, err
	}

	err = json.Unmarshal(b, &m)
	if err != nil {
		log.Fatal(err)
		return m, err
	}

	return m, nil
}
