package sitemap

import (
	"encoding/xml"
	"fmt"
	"sort"
	"time"
)

// URLSet represents the root element of a sitemap
type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

// URL represents a single URL entry in the sitemap
type URL struct {
	XMLName    xml.Name  `xml:"url"`
	Loc        string    `xml:"loc"`
	LastMod    string    `xml:"lastmod,omitempty"`
	ChangeFreq string    `xml:"changefreq,omitempty"`
	Priority   float64   `xml:"priority,omitempty"`
	LastModded time.Time `xml:"-"` // Internal field for sorting
}

// NewURLSet creates a new URLSet with the standard sitemap namespace
func NewURLSet() *URLSet {
	return &URLSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  make([]URL, 0),
	}
}

// AddURL adds a new URL to the sitemap
func (us *URLSet) AddURL(loc string, lastMod time.Time) {
	url := URL{
		Loc:        loc,
		LastMod:    lastMod.Format("2006-01-02"),
		LastModded: lastMod,
	}
	us.URLs = append(us.URLs, url)
}

// SortByLastMod sorts URLs by last modification date in descending order
func (us *URLSet) SortByLastMod() {
	sort.Slice(us.URLs, func(i, j int) bool {
		return us.URLs[i].LastModded.After(us.URLs[j].LastModded)
	})
}

// Validate checks if the sitemap is valid according to the sitemap protocol
func (us *URLSet) Validate() error {
	if len(us.URLs) == 0 {
		return fmt.Errorf("sitemap must contain at least one URL")
	}

	if len(us.URLs) > 50000 {
		return fmt.Errorf("sitemap cannot contain more than 50,000 URLs")
	}

	for _, url := range us.URLs {
		if url.Loc == "" {
			return fmt.Errorf("URL location cannot be empty")
		}

		if len(url.Loc) > 2048 {
			return fmt.Errorf("URL location cannot exceed 2048 characters: %s", url.Loc)
		}

		if url.Priority != 0 && (url.Priority < 0.0 || url.Priority > 1.0) {
			return fmt.Errorf("URL priority must be between 0.0 and 1.0: %s", url.Loc)
		}

		if url.ChangeFreq != "" {
			validFreqs := map[string]bool{
				"always":  true,
				"hourly":  true,
				"daily":   true,
				"weekly":  true,
				"monthly": true,
				"yearly":  true,
				"never":   true,
			}
			if !validFreqs[url.ChangeFreq] {
				return fmt.Errorf("invalid change frequency for URL %s: %s", url.Loc, url.ChangeFreq)
			}
		}
	}

	return nil
}

// Size returns the number of URLs in the sitemap
func (us *URLSet) Size() int {
	return len(us.URLs)
}

// Clone creates a deep copy of the URLSet
func (us *URLSet) Clone() *URLSet {
	clone := NewURLSet()
	clone.URLs = make([]URL, len(us.URLs))
	copy(clone.URLs, us.URLs)
	return clone
}
