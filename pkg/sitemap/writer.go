package sitemap

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

// Writer handles sitemap file generation
type Writer struct {
	// Indentation for XML output
	indent bool
}

// NewWriter creates a new sitemap writer
func NewWriter(indent bool) *Writer {
	return &Writer{
		indent: indent,
	}
}

// WriteToFile writes the sitemap to a file
func (w *Writer) WriteToFile(urlset *URLSet, filename string) error {
	// Validate sitemap before writing
	if err := urlset.Validate(); err != nil {
		return fmt.Errorf("invalid sitemap: %w", err)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create or truncate the file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Create encoder
	encoder := xml.NewEncoder(file)
	if w.indent {
		encoder.Indent("", "  ")
	}

	// Write XML header
	if _, err := file.WriteString(xml.Header); err != nil {
		return fmt.Errorf("failed to write XML header: %w", err)
	}

	// Encode sitemap
	if err := encoder.Encode(urlset); err != nil {
		return fmt.Errorf("failed to encode sitemap: %w", err)
	}

	return nil
}

// WriteToString returns the sitemap as a string
func (w *Writer) WriteToString(urlset *URLSet) (string, error) {
	// Validate sitemap
	if err := urlset.Validate(); err != nil {
		return "", fmt.Errorf("invalid sitemap: %w", err)
	}

	// Marshal to XML
	var output []byte
	var err error
	if w.indent {
		output, err = xml.MarshalIndent(urlset, "", "  ")
	} else {
		output, err = xml.Marshal(urlset)
	}
	if err != nil {
		return "", fmt.Errorf("failed to marshal sitemap: %w", err)
	}

	// Combine XML header and content
	return xml.Header + string(output), nil
}

// Compare compares two sitemaps and returns the differences
func (w *Writer) Compare(original, new *URLSet) (added, removed []URL) {
	// Create maps for quick lookup
	originalURLs := make(map[string]URL)
	for _, url := range original.URLs {
		originalURLs[url.Loc] = url
	}

	newURLs := make(map[string]URL)
	for _, url := range new.URLs {
		newURLs[url.Loc] = url
	}

	// Find added URLs
	for loc, url := range newURLs {
		if _, exists := originalURLs[loc]; !exists {
			added = append(added, url)
		}
	}

	// Find removed URLs
	for loc, url := range originalURLs {
		if _, exists := newURLs[loc]; !exists {
			removed = append(removed, url)
		}
	}

	return added, removed
}

// ValidateFile checks if an existing sitemap file is valid
func (w *Writer) ValidateFile(filename string) error {
	// Read file
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse XML
	var urlset URLSet
	if err := xml.Unmarshal(data, &urlset); err != nil {
		return fmt.Errorf("failed to parse sitemap: %w", err)
	}

	// Validate content
	if err := urlset.Validate(); err != nil {
		return fmt.Errorf("sitemap validation failed: %w", err)
	}

	return nil
}
