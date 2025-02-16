package sitemap

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Builder handles the construction of sitemaps
type Builder struct {
	// baseURL is the starting point URL
	baseURL *url.URL

	// urlset is the sitemap being built
	urlset *URLSet

	// options for sitemap generation
	options BuilderOptions
}

// BuilderOptions configures the sitemap building process
type BuilderOptions struct {
	// DefaultChangeFreq is the default change frequency for URLs
	DefaultChangeFreq string

	// DefaultPriority is the default priority for URLs
	DefaultPriority float64

	// IncludeLastMod determines if lastmod dates should be included
	IncludeLastMod bool

	// SortByLastMod determines if URLs should be sorted by last modification date
	SortByLastMod bool

	// ExcludePaths are paths that should be excluded from the sitemap
	ExcludePaths []string

	// StripQueryParams determines if query parameters should be stripped from URLs
	StripQueryParams bool
}

// DefaultBuilderOptions returns the default options for sitemap building
func DefaultBuilderOptions() BuilderOptions {
	return BuilderOptions{
		DefaultChangeFreq: "weekly",
		DefaultPriority:   0.5,
		IncludeLastMod:    true,
		SortByLastMod:     true,
		StripQueryParams:  true,
	}
}

// NewBuilder creates a new sitemap builder
func NewBuilder(baseURL *url.URL, options BuilderOptions) *Builder {
	return &Builder{
		baseURL: baseURL,
		urlset:  NewURLSet(),
		options: options,
	}
}

// AddURL adds a URL to the sitemap
func (b *Builder) AddURL(loc string, lastMod time.Time) error {
	// Parse the URL to validate it
	parsedURL, err := url.Parse(loc)
	if err != nil {
		return fmt.Errorf("invalid URL %s: %w", loc, err)
	}

	// Ensure URL is in the same domain as base URL
	if parsedURL.Host != b.baseURL.Host {
		return fmt.Errorf("URL %s is not in the same domain as base URL", loc)
	}

	// Check against excluded paths
	for _, excludePath := range b.options.ExcludePaths {
		if parsedURL.Path == excludePath {
			return nil
		}
	}

	// Strip query parameters if configured
	if b.options.StripQueryParams {
		parsedURL.RawQuery = ""
	}

	// Create URL entry
	url := URL{
		Loc:        parsedURL.String(),
		LastModded: lastMod,
	}

	// Add optional fields based on configuration
	if b.options.IncludeLastMod {
		url.LastMod = lastMod.Format("2006-01-02")
	}

	if b.options.DefaultChangeFreq != "" {
		url.ChangeFreq = b.options.DefaultChangeFreq
	}

	if b.options.DefaultPriority != 0 {
		url.Priority = b.options.DefaultPriority
	}

	// Add to urlset
	b.urlset.URLs = append(b.urlset.URLs, url)
	return nil
}

// Build finalizes and returns the sitemap
func (b *Builder) Build() (*URLSet, error) {
	// Sort URLs if configured
	if b.options.SortByLastMod {
		sort.Slice(b.urlset.URLs, func(i, j int) bool {
			return b.urlset.URLs[i].LastModded.After(b.urlset.URLs[j].LastModded)
		})
	}

	// Validate the sitemap
	if err := b.urlset.Validate(); err != nil {
		return nil, fmt.Errorf("sitemap validation failed: %w", err)
	}

	return b.urlset, nil
}

// SetChangeFreq sets the change frequency for all URLs
func (b *Builder) SetChangeFreq(freq string) {
	for i := range b.urlset.URLs {
		b.urlset.URLs[i].ChangeFreq = freq
	}
}

// SetPriority sets the priority for all URLs
func (b *Builder) SetPriority(priority float64) {
	for i := range b.urlset.URLs {
		b.urlset.URLs[i].Priority = priority
	}
}

// Clear removes all URLs from the sitemap
func (b *Builder) Clear() {
	b.urlset = NewURLSet()
}

// Count returns the number of URLs in the sitemap
func (b *Builder) Count() int {
	return len(b.urlset.URLs)
}

// GetURLs returns all URLs currently in the sitemap
func (b *Builder) GetURLs() []URL {
	urls := make([]URL, len(b.urlset.URLs))
	copy(urls, b.urlset.URLs)
	return urls
}
