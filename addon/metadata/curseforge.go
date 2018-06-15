package metadata

import (
	"errors"
	"io"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Metadata struct {
	Name        string
	Description string
	Lastmod     time.Time
	GameVersion string
}

// CurseforgeReader reader struct for addon page
type CurseforgeReader struct {
	doc *goquery.Document
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

// LastMod of the addon
func (r *CurseforgeReader) LastMod() (time.Time, error) {
	stdDate := r.doc.Find("span.stats--last-updated abbr.standard-date").First()
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

// Name of the addon
func (r *CurseforgeReader) Name() (string, error) {
	return r.doc.Find("h2.name").Text(), nil
}

// GameVersion Return the addon game version
func (r *CurseforgeReader) GameVersion() (string, error) {
	re := regexp.MustCompile(`\d{1,2}\.\d{1,2}.\d{1,2}`)
	spanText := r.doc.Find("span.stats--game-version").Text()
	return re.FindString(spanText), nil
}
