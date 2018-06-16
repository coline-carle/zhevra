package curseforge

import (
	"bytes"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"testing"
	"time"
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

func readerFromFixure(t *testing.T, name string) *AddonPageReader {
	b := bytes.NewReader(fixtures[name])
	r, err := NewAddonPageReader(b)
	if err != nil {
		t.Fatalf("NewAddonPageReader: unexpected error: %s\n", err)
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
		actual, err := r.Name()
		if err != nil {
			t.Fatalf("unexpected error %s\n", err)
		}
		if actual != test.expected {
			t.Errorf("expected '%s', got: '%s'", test.expected, actual)
		}
	}
}

var lastModTest = []struct {
	fixture  string
	expected time.Time
}{
	{"bigwigs", time.Unix(1528927095, 0)},
}

func TestCurseforgeLastMod(t *testing.T) {
	for _, test := range lastModTest {
		r := readerFromFixure(t, test.fixture)
		actual, err := r.LastMod()
		if err != nil {
			t.Fatalf("unexpected error %s\n", err)
		}
		if actual != test.expected {
			t.Errorf("expected '%s', got: '%s'", test.expected, actual)
		}
	}
}

var idTests = []struct {
	fixture  string
	expected int64
}{
	{"bigwigs", 2382},
}

func TestCurseforgeID(t *testing.T) {
	for _, test := range idTests {
		r := readerFromFixure(t, test.fixture)
		actual, err := r.ID()
		if err != nil {
			t.Fatalf("unexpected error %s\n", err)
		}
		if actual != test.expected {
			t.Errorf("expected '%d', got: '%d'", test.expected, actual)
		}
	}
}

var downloadsTests = []struct {
	fixture  string
	expected int64
}{
	{"bigwigs", 33340527},
}

func TestCurseforgeDownloads(t *testing.T) {
	for _, test := range downloadsTests {
		r := readerFromFixure(t, test.fixture)
		actual, err := r.Downloads()
		if err != nil {
			t.Fatalf("unexpected error %s\n", err)
		}
		if actual != test.expected {
			t.Errorf("expected '%d', got: '%d'", test.expected, actual)
		}
	}
}

var gameVersionTests = []struct {
	fixture  string
	expected string
}{
	{"bigwigs", "8.0.1"},
}

func TestCurseforgeGameVersion(t *testing.T) {
	for _, test := range gameVersionTests {
		r := readerFromFixure(t, test.fixture)
		actual, err := r.GameVersion()
		if err != nil {
			t.Fatalf("unexpected error %s\n", err)
		}
		if actual != test.expected {
			t.Errorf("expected '%s', got: '%s'", test.expected, actual)
		}
	}
}

var descriptionTests = []struct {
	fixture  string
	expected string
}{
	{
		"bigwigs",
		"Modular, lightweight, non-intrusive approach to boss encounter warnings." +
			" The efficiently coded alternative to Deadly Boss Mods (DBM) for spell &" +
			" ability...",
	},
}

func TestCurseforgeDescription(t *testing.T) {
	for _, test := range descriptionTests {
		r := readerFromFixure(t, test.fixture)
		actual, err := r.Description()
		if err != nil {
			t.Fatalf("unexpected error %s\n", err)
		}
		if actual != test.expected {
			t.Errorf("expected '%s', got: '%s'", test.expected, actual)
		}
	}
}
