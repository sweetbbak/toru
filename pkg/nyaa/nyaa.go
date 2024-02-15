package nyaa

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/docker/go-units"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

var NyaaURL = "https://nyaa.si"

var PageLimit uint = 100

// set a proxy url like nyaa.iss.ink
// set a proxy url like nyaa.iss.ink
func SetProxyURL(proxyURL string) {
	NyaaURL = proxyURL
}

func Search(search string, parameters ...SearchParameters) ([]Media, error) {
	params, err := getOneParameterSet(parameters)
	if err != nil {
		return nil, err
	}

	doc, err := requestHTML(search, params)
	if err != nil {
		return nil, errors.Wrap(err, "error getting the nyaa page")
	}

	medias, err := parseSearchPageHTML(doc)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing html")
	}

	return medias, nil
}

func getOneParameterSet(parameters []SearchParameters) (SearchParameters, error) {
	params := SearchParameters{}
	if len(parameters) == 1 {
		params = parameters[0]
	}
	if len(parameters) > 1 {
		return SearchParameters{}, errors.New("only one parameter set accepted")
	}

	if params.Page > PageLimit {
		return params, errors.New("exceeded page limit")
	}

	return params, nil
}

func requestHTML(search string, params SearchParameters) (*goquery.Document, error) {
	URL, err := urlForParams(search, params)
	if err != nil {
		return nil, errors.Wrap(err, "error creating url for search")
	}

	rep, err := http.Get(URL)
	if err != nil {
		return nil, errors.Wrap(err, "error requesting results")
	}
	defer rep.Body.Close()

	if rep.StatusCode < 200 || rep.StatusCode >= 300 {
		return nil, errors.Errorf("non-OK HTTP status code: %d %s", rep.StatusCode, rep.Status)
	}

	doc, err := goquery.NewDocumentFromReader(rep.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing response html")
	}

	return doc, nil
}

func urlForParams(search string, parameters SearchParameters) (string, error) {
	var baseURL string
	if parameters.Proxy != "" {
		baseURL = parameters.Proxy
	} else {
		baseURL = NyaaURL
	}

	if parameters.User != "" {
		baseURL += "/user/" + url.PathEscape(parameters.User)
	}

	URL, err := url.Parse(baseURL)
	if err != nil {
		return "", errors.Wrap(err, "error parsing nyaa url")
	}

	query := URL.Query()
	query.Set("f", strconv.FormatInt(int64(parameters.Filter), 10))
	query.Set("c", string(parameters.Category))
	query.Set("q", search)
	query.Set("s", string(parameters.SortBy))
	query.Set("o", string(parameters.SortOrder))
	query.Set("p", strconv.FormatInt(int64(parameters.Page), 10))
	URL.RawQuery = query.Encode()

	return URL.String(), nil
}

func parseSearchPageHTML(doc *goquery.Document) ([]Media, error) {
	selection := doc.Find(".torrent-list tbody tr")
	medias := make([]Media, selection.Length())
	errChan := make(chan error)
	group := sync.WaitGroup{}
	group.Add(selection.Length())
	doneChan := make(chan struct{})
	go func() {
		group.Wait()
		doneChan <- struct{}{}
	}()

	go func() {
		selection.Each(func(i int, selection *goquery.Selection) {
			media, err := parseMediaElement(selection)
			if err != nil {
				errChan <- errors.Wrap(err, "error parsing media element")
				return
			}
			medias[i] = media
			group.Done()
		})
	}()

	select {
	case err := <-errChan:
		return nil, err
	case <-doneChan:
	}

	return medias, nil
}

func parseMediaElement(elem *goquery.Selection) (Media, error) {
	media := Media{}

	links := elem.Find("td a:not(.comments)").Nodes
	err := parseMediaElementLinks(links, &media)
	if err != nil {
		return media, err
	}

	nodes := elem.Find("td").Nodes
	err = parseMediaElementTexts(nodes, &media)
	if err != nil {
		return media, err
	}

	err = parseMediaElementTimestamp(nodes, &media)
	if err != nil {
		return media, err
	}

	err = parseMediaElementCommentCount(elem, &media)
	if err != nil {
		return media, err
	}

	return media, nil
}

func parseMediaElementLinks(links []*html.Node, media *Media) error {
	if len(links) != 4 {
		return errors.Errorf("unexpected layout: expected 4 links, got %d", len(links))
	}

	href, err := getLinkHref(links, 0)
	if err != nil {
		return err
	}
	media.Category = hrefToCategory(href)

	href, err = getLinkHref(links, 1)
	if err != nil {
		return err
	}
	id, err := hrefToID(href)
	if err != nil {
		return errors.Wrap(err, "error parsing ID")
	}
	media.ID = id
	title, ok := getAttributeValueByKey(links[1], "title")
	if !ok {
		return errors.New("unexpected layout: link 1 does not have a title")
	}
	media.Name = title

	href, err = getLinkHref(links, 2)
	if err != nil {
		return err
	}
	media.Torrent = href

	href, err = getLinkHref(links, 3)
	if err != nil {
		return err
	}
	media.Magnet = href

	return nil
}

func getLinkHref(links []*html.Node, index int) (string, error) {
	href, ok := getAttributeValueByKey(links[index], "href")
	if !ok {
		return "", errors.Errorf("unexpected layout: link %d does not have an href", index)
	}
	return href, nil
}

func parseMediaElementTexts(nodes []*html.Node, media *Media) error {
	if len(nodes) != 8 {
		return errors.Errorf("unexpected layout: expected 8 nodes, got %d", len(nodes))
	}

	data, err := getFirstChildText(nodes, 3)
	if err != nil {
		return err
	}
	size, err := units.FromHumanSize(data)
	if err != nil {
		return errors.Wrap(err, "error parsing size")
	}
	media.Size = uint64(size)

	data, err = getFirstChildText(nodes, 5)
	if err != nil {
		return err
	}
	seeders, err := strconv.Atoi(data)
	if err != nil {
		return errors.Wrap(err, "error parsing se")
	}
	media.Seeders = uint(seeders)

	data, err = getFirstChildText(nodes, 6)
	if err != nil {
		return err
	}
	leechers, err := strconv.Atoi(data)
	if err != nil {
		return errors.Wrap(err, "error parsing leechers")
	}
	media.Leechers = uint(leechers)

	data, err = getFirstChildText(nodes, 7)
	if err != nil {
		return err
	}
	downloads, err := strconv.Atoi(data)
	if err != nil {
		return errors.Wrap(err, "error parsing downloads")
	}
	media.Downloads = uint(downloads)

	return nil
}

func getFirstChildText(nodes []*html.Node, index int) (string, error) {
	child := nodes[index].FirstChild
	if child == nil || child.Type != html.TextNode {
		return "", errors.Errorf("unexpected layout: expected node %d to have a text first child", index)
	}
	return child.Data, nil
}

func parseMediaElementTimestamp(nodes []*html.Node, media *Media) error {
	timestamp, ok := getAttributeValueByKey(nodes[4], "data-timestamp")
	if !ok {
		return errors.New("unexpected layout: expected node 5 to have a data-timestamp")
	}
	timestampInt, err := strconv.Atoi(timestamp)
	if err != nil {
		return errors.Wrap(err, "error parsing timestamp")
	}
	media.Date = time.Unix(int64(timestampInt), 0)
	return nil
}

func parseMediaElementCommentCount(elem *goquery.Selection, media *Media) error {
	nodes := elem.Find(".comments").Nodes
	if len(nodes) == 0 {
		return nil
	}
	if len(nodes) > 1 {
		return errors.New("found more than one comments element")
	}

	textChild := nodes[0].LastChild
	if textChild == nil || textChild.Type != html.TextNode {
		return errors.New("unexpected layout: expected comments elem to have a text last child")
	}
	commentCount, err := strconv.Atoi(textChild.Data)
	if err != nil {
		return errors.Wrap(err, "error parsing comment count")
	}
	media.CommentCount = uint(commentCount)

	return nil
}

func getAttributeValueByKey(node *html.Node, key string) (string, bool) {
	for _, attribute := range node.Attr {
		if attribute.Key == key {
			return attribute.Val, true
		}
	}
	return "", false
}

func hrefToCategory(href string) Category {
	return Category(strings.TrimPrefix(href, "/?c="))
}

func hrefToID(href string) (uint, error) {
	id, err := strconv.Atoi(strings.TrimPrefix(href, "/view/"))
	return uint(id), err
}
