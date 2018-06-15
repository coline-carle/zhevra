package metadata_test

import (
	"bytes"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/wow-sweetlie/zhevra/addon/metadata"
)

var fixtures map[string][]byte
var bigwigs []byte

func init() {
	fixtures = make(map[string][]byte)
	files, err := filepath.Glob("./fixtures/curseforge/*.html")
	if err != nil {
		log.Fatalf("error loading curseforge fixtures: %s", err)
	}
	for _, file := range files {
		name := filepath.Base(file)
		name = strings.Trim(name, ".html")
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatalf("error loading fixture: %s - %s", file, err)
		}
		fixtures[name] = data
		log.Printf("loaded: %s - %s\n", name, file)
	}
}

func readerFromFixure(t *testing.T, name string) *metadata.CurseforgeReader {
	b := bytes.NewReader(fixtures[name])
	r, err := metadata.NewCurseforgeReader(b)
	if err != nil {
		t.Fatalf("NewCurseforgeReader: unexpected error: %s\n", err)
	}
	return r
}

var nameTests = []struct {
	fixture  string
	expected string
}{
	{"bigwigs", "BigWigs Bossmods"},
}

func TestCurseforgeName(t *testing.T) {
	for _, test := range nameTests {
		r := readerFromFixure(t, test.fixture)
		actual := r.Name()
		if actual != test.expected {
			t.Errorf("r.Name(): expected '%s', got: '%s'", test.expected, actual)
		}
	}
}

var lastModTest = []struct {
	fixture  string
	expected time.Time
}{
	{"bigwigs", time.Unix(1527139371, 0)},
}

func TestCurseforgeLastMod(t *testing.T) {
	for _, test := range lastModTest {
		r := readerFromFixure(t, test.fixture)
		actual, err := r.LastMod()
		if err != nil {
			t.Fatalf("r.Lastmod(): unexpected error %s\n", err)
		}
		if actual != test.expected {
			t.Errorf("r.LastMod(): expected '%s', got: '%s'", test.expected, actual)
		}
	}
}
