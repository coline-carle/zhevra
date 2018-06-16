package metadata

import (
	"encoding/json"
	"errors"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// CurseforgeReader reader struct for addon page
type CurseforgeReader struct {
	doc *goquery.Document
}

// TwitchButton on addon page for download install
type TwitchButton struct {
	ProjectID   int64  `json:"ProjectID"`
	ProjectName string `json:"ProjectName"`
}

// NewCurseforgeReader create a reader for curseforge addon page
func NewCurseforgeReader(r io.Reader) (*CurseforgeReader, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	return &CurseforgeReader{
		doc: doc,
	}, nil
}

// Downloads number of download  of the addon
func (r *CurseforgeReader) Downloads() (int64, error) {
	var strDownloads string
	re := regexp.MustCompile(`^\s*(?:\d{1,3},?)+\s*`)
	r.doc.Find("div.infobox__content > p > span").EachWithBreak(
		func(i int, span *goquery.Selection) bool {
			text := span.Text()
			if re.MatchString(text) {
				strDownloads = text
				return false
			}
			return true
		})
	strDownloads = strings.Map(func(r rune) rune {
		if r == ',' {
			return -1
		}
		return r
	}, strDownloads)
	return strconv.ParseInt(strDownloads, 10, 64)
}

// LastMod of the addon
func (r *CurseforgeReader) LastMod() (time.Time, error) {
	stdDate := r.doc.Find("div.infobox__content abbr.standard-date").First()
	epochStr, exists := stdDate.Attr("data-epoch")
	if !exists {
		return time.Time{}, errors.New("lastmod attribute not found")
	}
	epoch, err := strconv.ParseInt(epochStr, 10, 64)
	if err != nil {
		return time.Time{}, errors.New("lastmod could not convert epoch to int")
	}
	t := time.Unix(epoch, 0)
	return t, nil
}

// Description of he addon
func (r *CurseforgeReader) Description() (string, error) {
	descNode := r.doc.Find(`meta[property="og:description"]`).First()
	descContent, exists := descNode.Attr("content")
	if !exists {
		return "", errors.New("description attribute not found")
	}
	return descContent, nil
}

// ID cuse id of the addon
func (r *CurseforgeReader) ID() (int64, error) {
	twitchNode := r.doc.Find(`a.button--twitch`)
	jsonData, exists := twitchNode.Attr("data-action-value")
	if !exists {
		return 0, errors.New("twitchnode not found")
	}
	var twitchButton TwitchButton
	err := json.NewDecoder(strings.NewReader(jsonData)).Decode(&twitchButton)
	if err != nil {
		return 0, err
	}
	return twitchButton.ProjectID, nil
}

// Name of the addon
func (r *CurseforgeReader) Name() (string, error) {
	titleNode := r.doc.Find(`meta[property="og:title"]`).First()
	titleContent, exists := titleNode.Attr("content")
	if !exists {
		return "", errors.New("name attribute not found")
	}
	return titleContent, nil
}

// Upstream return the url of the project
func (r *CurseforgeReader) Upstream() (string, error) {
	upstreamNode := r.doc.Find("p.infobox__cta a").First()
	upstreamHref, exists := upstreamNode.Attr("href")
	if !exists {
		return "", errors.New("upstream attribute not found")
	}
	return upstreamHref, nil
}

// GameVersion Return the addon game version
func (r *CurseforgeReader) GameVersion() (string, error) {
	re := regexp.MustCompile(`\d{1,2}\.\d{1,2}.\d{1,2}`)
	spanText := r.doc.Find("span.stats--game-version").Text()
	return re.FindString(spanText), nil
}
