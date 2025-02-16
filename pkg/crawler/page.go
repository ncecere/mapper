package crawler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// Page represents a crawled web page
type Page struct {
	// URL of the page
	URL *url.URL

	// Depth represents how many links deep this page is from the start URL
	Depth int

	// LastModified is the last modification time of the page
	// This is extracted from the Last-Modified header or current time if not available
	LastModified time.Time

	// Links contains all unique URLs found on the page
	Links []*url.URL

	// Error holds any error encountered while processing the page
	Error error
}

// NewPage creates a new Page instance
func NewPage(pageURL *url.URL, depth int) *Page {
	return &Page{
		URL:   pageURL,
		Depth: depth,
		Links: make([]*url.URL, 0),
	}
}

// Process fetches and processes the page content
func (p *Page) Process(client *http.Client) error {
	req, err := http.NewRequest(http.MethodGet, p.URL.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Extract last modified time
	if lastMod := resp.Header.Get("Last-Modified"); lastMod != "" {
		if t, err := time.Parse(time.RFC1123, lastMod); err == nil {
			p.LastModified = t
		}
	}
	if p.LastModified.IsZero() {
		p.LastModified = time.Now()
	}

	return p.parseHTML(resp.Body)
}

// parseHTML parses the HTML content and extracts links
func (p *Page) parseHTML(body io.Reader) error {
	doc, err := html.Parse(body)
	if err != nil {
		return fmt.Errorf("failed to parse HTML: %w", err)
	}

	var links []*url.URL
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			// Check for <a> tags with href
			if n.Data == "a" {
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						if link := p.normalizeURL(attr.Val); link != nil {
							links = append(links, link)
						}
						break
					}
				}
			}

			// Check for <link> tags with href (e.g., for canonical URLs)
			if n.Data == "link" {
				var rel, href string
				for _, attr := range n.Attr {
					switch attr.Key {
					case "rel":
						rel = attr.Val
					case "href":
						href = attr.Val
					}
				}
				if (rel == "canonical" || rel == "alternate") && href != "" {
					if link := p.normalizeURL(href); link != nil {
						links = append(links, link)
					}
				}
			}
		}

		// Traverse child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)
	p.Links = uniqueURLs(links)
	return nil
}

// normalizeURL converts a relative or absolute URL to a normalized absolute URL
func (p *Page) normalizeURL(rawURL string) *url.URL {
	// Remove fragment
	if idx := strings.Index(rawURL, "#"); idx != -1 {
		rawURL = rawURL[:idx]
	}

	// Skip empty URLs and javascript: links
	if rawURL == "" || strings.HasPrefix(rawURL, "javascript:") {
		return nil
	}

	// Parse the URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil
	}

	// Convert relative URLs to absolute
	if !parsedURL.IsAbs() {
		parsedURL = p.URL.ResolveReference(parsedURL)
	}

	// Ensure URL has a scheme
	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "https"
	}

	return parsedURL
}

// uniqueURLs removes duplicate URLs from a slice while preserving order
func uniqueURLs(urls []*url.URL) []*url.URL {
	seen := make(map[string]bool)
	unique := make([]*url.URL, 0, len(urls))

	for _, u := range urls {
		key := u.String()
		if !seen[key] {
			seen[key] = true
			unique = append(unique, u)
		}
	}

	return unique
}
