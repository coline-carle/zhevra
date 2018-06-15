package metadata

import (
	"errors"
	"io"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

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

// Name of the addon
func (r *CurseforgeReader) Name() string {
	return r.doc.Find("h2.name").Text()
}
