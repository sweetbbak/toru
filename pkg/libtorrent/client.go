package libtorrent

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"
)

const (
	INTERNAL_STREAM_PORT = "8888"
)

type Client struct {
	// client / project name, will be the default directory name
	Name string
	// directory to download torrents to
	DataDir string
	// Seed or no
	Seed bool
	// Port to stream torrents on
	Port string
	// Port to stream torrents on
	TorrentPort int
	// Default torrent client options
	TorrentClient *torrent.Client
	// server
	srv *http.Server
	// torrents
	Torrents []*torrent.Torrent
	// Disable IPV6
	DisableIPV6 bool
}

// create a default client, must call Init afterwords
func NewClient(name string, port string) *Client {
	return &Client{
		Name: name,
		Port: port,
		Seed: false,
	}
}

// Initialize torrent configuration
func (c *Client) Init() error {
	cfg := torrent.NewDefaultClientConfig()
	s, err := c.getStorage()
	if err != nil {
		return err
	}

	cfg.DisableIPv6 = c.DisableIPV6

	// sanity check - get open port to allow for multiple instances of toru
	if c.TorrentPort < 5 {
		port, err := GetFreePort()
		if err != nil {
			c.TorrentPort = 42069
		} else {
			c.TorrentPort = port
		}
	}

	if c.Port == INTERNAL_STREAM_PORT {
		p, err := GetFreePortString()
		if err != nil {
			c.Port = INTERNAL_STREAM_PORT
		} else {
			c.Port = p
		}
	}

	cfg.ListenPort = c.TorrentPort
	c.DataDir = s
	cfg.DefaultStorage = storage.NewFileByInfoHash(c.DataDir)

	client, err := torrent.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("error creating a new torrent client: %v", err)
	}

	c.StartServer()
	c.TorrentClient = client
	return nil
}

// add a torrent and mark it entirely for download
func (c *Client) DownloadTorrent(torrent string) error {
	t, err := c.AddTorrent(torrent)
	if err != nil {
		return err
	}
	t.DownloadAll()
	return nil
}

// Is torrent file a video
func IsVideoFile(f *torrent.File) bool {
	ext := path.Ext(f.Path())
	switch ext {
	case ".mp4", ".mkv", ".avi", ".avif", ".av1", ".mov", ".flv", ".f4v", ".webm", ".wmv", ".mpeg", ".mpg", ".mlv", ".hevc", ".flac", ".flic":
		return true
	default:
		return false
	}
}

// get a free port in string format
func GetFreePortString() (string, error) {
	port, err := GetFreePort()
	if err != nil {
		return INTERNAL_STREAM_PORT, err
	}
	return fmt.Sprintf("%d", port), nil
}

// get a free port
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// HTTP handler for ServeTorrent allowing for query by hash for the torrent and filepath for the subsequent video
func (c *Client) handler(w http.ResponseWriter, r *http.Request) {
	ts := c.TorrentClient.Torrents()
	queries := r.URL.Query()
	// get hash of torrent
	hash := queries.Get("hash")

	// check if the user requested a specific episode
	fpath := queries.Get("filepath")

	// idk why but this is always mangled af
	hash = strings.TrimSpace(hash)
	hash = strings.ReplaceAll(hash, "\n", "")

	fpath = strings.TrimSpace(fpath)
	fpath = strings.ReplaceAll(fpath, "\n", "")

	if hash == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		log.Println("server handler: Hash query is empty")
		return
	}

	var targetTorrent *torrent.Torrent

	for _, ff := range ts {
		<-ff.GotInfo()
		ih := ff.InfoHash().String()
		if ih == hash {
			targetTorrent = ff
		}
	}

	if targetTorrent == nil {
		http.Error(w, http.StatusText(400), http.StatusInternalServerError)
		log.Println("server handler: couldnt find torrent by infohash")
		return
	}

	// file count
	fileCount := len(targetTorrent.Files())

	// the file we will serve after checking user input
	var targetFile *torrent.File

	if fileCount == 1 {
		// TODO: do bounds checking and make sure the file isn't a txt file or something
		targetFile = targetTorrent.Files()[0]
	}

	if fpath != "" && fileCount > 1 {
		for _, f := range targetTorrent.Files() {
			// TODO: I have a feeling this will be problematic with weird characters and stuff, handle that.
			if f.DisplayPath() == fpath {
				targetFile = f
				break
			}
		}
	}

	if targetFile == nil {
		http.Error(w, http.StatusText(400), http.StatusInternalServerError)
		log.Println("server handler: couldnt find torrent file requested")
		return
	}

	w.Header().Set("Content-Type", "video/mp4")
	http.ServeContent(w, r, targetFile.DisplayPath(), time.Unix(targetFile.Torrent().Metainfo().CreationDate, 0), targetFile.NewReader())
}

// start the server in the background, only done once internally.
func (c *Client) StartServer() {
	// :8080 for localhost:8080/
	port := fmt.Sprintf(":%s", c.Port)
	c.srv = &http.Server{Addr: port}
	http.HandleFunc("/stream", c.handler)

	go func() {
		if err := c.srv.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				return
			} else {
				log.Fatal(err)
			}
		}
	}()
}

// Generate a link that can be used with the default clients server to play a torrent
// that is already loaded into the client and allow for the specification of a file to play by filepath
func (c *Client) ServeTorrentEpisode(t *torrent.Torrent, filePath string) string {
	mh := t.InfoHash().String()
	return fmt.Sprintf("http://localhost:%s/stream?hash=%s&filepath=%s", c.Port, mh, filePath)
}

// Generate a link that can be used with the default clients server to play a torrent
// that is already loaded into the client
func (c *Client) ServeTorrent(t *torrent.Torrent) string {
	mh := t.InfoHash().String()
	return fmt.Sprintf("http://localhost:%s/stream?hash=%s", c.Port, mh)
}

// returns a slice of loaded torrents or nil
func (c *Client) ShowTorrents() []*torrent.Torrent {
	return c.TorrentClient.Torrents()
}

// generic add torrent function, takes magnets, URLs to torrent files and torrent files.
func (c *Client) AddTorrent(tor string) (*torrent.Torrent, error) {
	if strings.HasPrefix(tor, "magnet") {
		return c.AddMagnet(tor)
	} else if strings.Contains(tor, "http") {
		return c.AddTorrentURL(tor)
	} else {
		return c.AddTorrentFile(tor)
	}
}

func (c *Client) AddMagnet(magnet string) (*torrent.Torrent, error) {
	t, err := c.TorrentClient.AddMagnet(magnet)
	if err != nil {
		return nil, err
	}
	<-t.GotInfo()
	return t, nil
}

func (c *Client) AddTorrentFile(file string) (*torrent.Torrent, error) {
	t, err := c.TorrentClient.AddTorrentFromFile(file)
	if err != nil {
		return nil, err
	}
	<-t.GotInfo()
	return t, nil
}

func (c *Client) AddTorrentURL(url string) (*torrent.Torrent, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fname := path.Base(url)
	tmp := os.TempDir()
	path.Join(tmp, fname)

	file, err := os.Create(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return nil, err
	}

	t, err := c.TorrentClient.AddTorrentFromFile(file.Name())
	if err != nil {
		return nil, err
	}
	<-t.GotInfo()
	return t, nil
}

// stops the client and closes all connections to peers
func (c *Client) Close() (errs []error) {
	return c.TorrentClient.Close()
}

// look through the torrent files the client is handling and return a torrent with a
// matching info hash
func (c *Client) FindByInfoHhash(infoHash string) (*torrent.Torrent, error) {
	torrents := c.TorrentClient.Torrents()
	for _, t := range torrents {
		if t.InfoHash().AsString() == infoHash {
			return t, nil
		}
	}
	return nil, fmt.Errorf("No torrents match info hash: %v", infoHash)
}

func (c *Client) DropTorrent(t *torrent.Torrent) {
	t.Drop()
}

// Create storage path if it doesn't exist and return Path
func (c *Client) getStorage() (string, error) {
	s, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("Client error, couldnt get user cache directory: %v", err)
	}

	p := path.Join(s, c.Name)
	if p == "" || c.Name == "" {
		return "", fmt.Errorf("Client error, couldnt construct client path: Empty path or project name")
	}

	err = os.MkdirAll(p, 0o755)
	if err != nil {
		return "", fmt.Errorf("Client error, couldnt create project directory: %v", err)
	}

	_, err = os.Stat(p)
	if err == nil {
		return p, nil
	} else {
		return "", err
	}
}
