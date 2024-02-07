package nyaa

import (
	"io/fs"
	"strconv"
	"time"
)

type Filter uint8

const (
	FilterNoFilter Filter = iota
	FilterNoRemakes
	FilterTrustedOnly
)

type Category string

const (
	CategoryNone                           Category = ""
	CategoryAllCategories                  Category = "0_0"
	CategoryAnime                          Category = "1_0"
	CategoryAnimeMusicVideo                Category = "1_1"
	CategoryAnimeEnglishTranslated         Category = "1_2"
	CategoryAnimeNonEnglishTranslated      Category = "1_3"
	CategoryAnimeRaw                       Category = "1_4"
	CategoryAudio                          Category = "2_0"
	CategoryAudioLossless                  Category = "2_1"
	CategoryAudioLossy                     Category = "2_2"
	CategoryLiterature                     Category = "3_0"
	CategoryLiteratureEnglishTranslated    Category = "3_1"
	CategoryLiteratureNonEnglishTranslated Category = "3_2"
	CategoryLiteratureRaw                  Category = "3_3"
	CategoryLiveAction                     Category = "4_0"
	CategoryLiveActionEnglishTranslated    Category = "4_1"
	CategoryLiveActionIdolPromotionalVideo Category = "4_2"
	CategoryLiveActionNonEnglishTranslated Category = "4_3"
	CategoryLiveActionRaw                  Category = "4_4"
	CategoryPictures                       Category = "5_0"
	CategoryPicturesGraphics               Category = "5_1"
	CategoryPicturesPhotos                 Category = "5_2"
	CategorySoftware                       Category = "6_0"
	CategorySoftwareApplications           Category = "6_1"
	CategorySoftwareGames                  Category = "6_2"
)

type SortBy string

const (
	SortByNone      SortBy = ""
	SortByComments  SortBy = "comments"
	SortBySize      SortBy = "size"
	SortByDate      SortBy = "id"
	SortBySeeders   SortBy = "seeders"
	SortByLeechers  SortBy = "leechers"
	SortByDownloads SortBy = "downloads"
)

type SortOrder string

const (
	SortOrderNone       = ""
	SortOrderDescending = "desc"
	SortOrderAscending  = "asc"
)

type SearchParameters struct {
	Filter    Filter
	Category  Category
	User      string
	SortBy    SortBy
	SortOrder SortOrder
	Page      uint
}

type Media struct {
	MediaPartial
	IsFull bool

	Submitter   string
	Information string
	Hash        string
	Description string
	Files       []fs.FileInfo
	Comments    []Comment
}

type MediaPartial struct {
	Name         string
	ID           uint
	Category     Category
	Torrent      string
	Magnet       string
	Size         uint64
	Date         time.Time
	Seeders      uint
	Leechers     uint
	Downloads    uint
	CommentCount uint
}

func (m *MediaPartial) ViewURL() string {
	return NyaaURL + "/view/" + strconv.Itoa(int(m.ID))
}

type FileInfo struct {
	Name  string
	Size  uint64
	Path  []string
	IsDir bool
}

type Comment struct {
	Author          string
	AuthorAvatarURL string
	Date            time.Time
}
