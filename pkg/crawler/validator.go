package crawler

import (
	"net/url"
	"regexp"
	"strings"
)

// URLValidator handles URL validation and filtering
type URLValidator struct {
	// baseURL is the starting point URL
	baseURL *url.URL

	// excludePatterns contains compiled regex patterns for URLs to exclude
	excludePatterns []*regexp.Regexp

	// includePatterns contains compiled regex patterns for URLs to include
	includePatterns []*regexp.Regexp
}

// NewURLValidator creates a new URLValidator instance
func NewURLValidator(baseURL *url.URL, excludePatterns, includePatterns []string) (*URLValidator, error) {
	v := &URLValidator{
		baseURL:         baseURL,
		excludePatterns: make([]*regexp.Regexp, 0, len(excludePatterns)),
		includePatterns: make([]*regexp.Regexp, 0, len(includePatterns)),
	}

	// Compile exclude patterns
	for _, pattern := range excludePatterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		v.excludePatterns = append(v.excludePatterns, re)
	}

	// Compile include patterns
	for _, pattern := range includePatterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		v.includePatterns = append(v.includePatterns, re)
	}

	return v, nil
}

// IsValid checks if a URL should be crawled based on various criteria
func (v *URLValidator) IsValid(u *url.URL) bool {
	// Skip empty URLs
	if u == nil || u.String() == "" {
		return false
	}

	// Skip URLs with unsupported schemes
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	// Skip URLs not in the same domain
	if u.Host != v.baseURL.Host {
		return false
	}

	// Skip common non-content file types
	if v.isNonContentFile(u.Path) {
		return false
	}

	// Check against exclude patterns
	urlStr := u.String()
	for _, pattern := range v.excludePatterns {
		if pattern.MatchString(urlStr) {
			return false
		}
	}

	// If include patterns are specified, URL must match at least one
	if len(v.includePatterns) > 0 {
		matched := false
		for _, pattern := range v.includePatterns {
			if pattern.MatchString(urlStr) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	return true
}

// isNonContentFile checks if the URL points to a non-HTML resource
func (v *URLValidator) isNonContentFile(path string) bool {
	// List of file extensions to skip
	nonContentExts := []string{
		".jpg", ".jpeg", ".png", ".gif", ".ico", ".css", ".js",
		".pdf", ".doc", ".docx", ".ppt", ".pptx", ".xls", ".xlsx",
		".zip", ".tar", ".gz", ".rar", ".exe", ".mp3", ".mp4",
		".avi", ".mov", ".wmv", ".flv", ".svg", ".woff", ".woff2",
		".ttf", ".eot",
	}

	path = strings.ToLower(path)
	for _, ext := range nonContentExts {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}

	return false
}

// NormalizePath ensures consistent path formatting
func (v *URLValidator) NormalizePath(path string) string {
	// Remove trailing slash unless it's the root path
	if path != "/" && strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}

	// Ensure path starts with /
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return path
}

// ShouldFollowRedirect determines if a redirect should be followed
func (v *URLValidator) ShouldFollowRedirect(redirectURL *url.URL) bool {
	// Only follow redirects to the same domain
	return redirectURL.Host == v.baseURL.Host
}

// GetDomain returns the domain being crawled
func (v *URLValidator) GetDomain() string {
	return v.baseURL.Host
}

// IsSubpath checks if a URL path is a subpath of another
func (v *URLValidator) IsSubpath(parent, child string) bool {
	parent = v.NormalizePath(parent)
	child = v.NormalizePath(child)
	return strings.HasPrefix(child, parent)
}
