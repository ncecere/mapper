package crawler

import (
	"fmt"
	"net/url"
	"time"
)

// Config holds the configuration for the crawler
type Config struct {
	// BaseURL is the starting point for crawling
	BaseURL *url.URL

	// MaxDepth defines how deep the crawler should traverse from the base URL
	// A depth of 0 means only crawl the base URL
	// A depth of 1 means crawl the base URL and all directly linked pages
	MaxDepth int

	// MaxConcurrent defines the maximum number of concurrent requests
	MaxConcurrent int

	// RequestTimeout defines the timeout for each HTTP request
	RequestTimeout time.Duration

	// RateLimit defines the minimum time between requests to the same domain
	RateLimit time.Duration

	// UserAgent is the User-Agent string to use in HTTP requests
	UserAgent string

	// FollowRedirects determines if the crawler should follow HTTP redirects
	FollowRedirects bool

	// ExcludePatterns contains regex patterns for URLs to exclude from crawling
	ExcludePatterns []string

	// IncludePatterns contains regex patterns for URLs to include in crawling
	// If empty, all URLs not matching exclude patterns are included
	IncludePatterns []string
}

// DefaultConfig returns a Config with sensible default values
func DefaultConfig(baseURL string) (*Config, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	return &Config{
		BaseURL:         parsedURL,
		MaxDepth:        3,
		MaxConcurrent:   5,
		RequestTimeout:  10 * time.Second,
		RateLimit:       time.Second,
		UserAgent:       "Mapper/1.0 (+https://github.com/ncecere/mapper)",
		FollowRedirects: true,
	}, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.BaseURL == nil {
		return fmt.Errorf("base URL is required")
	}

	if c.MaxDepth < 0 {
		return fmt.Errorf("max depth must be non-negative")
	}

	if c.MaxConcurrent < 1 {
		return fmt.Errorf("max concurrent must be at least 1")
	}

	if c.RequestTimeout < time.Second {
		return fmt.Errorf("request timeout must be at least 1 second")
	}

	if c.RateLimit < 0 {
		return fmt.Errorf("rate limit must be non-negative")
	}

	if c.UserAgent == "" {
		return fmt.Errorf("user agent is required")
	}

	return nil
}

// Option defines a function that modifies a Config
type Option func(*Config)

// WithMaxDepth sets the maximum crawl depth
func WithMaxDepth(depth int) Option {
	return func(c *Config) {
		c.MaxDepth = depth
	}
}

// WithMaxConcurrent sets the maximum number of concurrent requests
func WithMaxConcurrent(max int) Option {
	return func(c *Config) {
		c.MaxConcurrent = max
	}
}

// WithRequestTimeout sets the request timeout
func WithRequestTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.RequestTimeout = timeout
	}
}

// WithRateLimit sets the rate limit
func WithRateLimit(limit time.Duration) Option {
	return func(c *Config) {
		c.RateLimit = limit
	}
}

// WithUserAgent sets the User-Agent string
func WithUserAgent(userAgent string) Option {
	return func(c *Config) {
		c.UserAgent = userAgent
	}
}

// WithFollowRedirects sets whether to follow redirects
func WithFollowRedirects(follow bool) Option {
	return func(c *Config) {
		c.FollowRedirects = follow
	}
}

// WithExcludePatterns sets the URL patterns to exclude
func WithExcludePatterns(patterns []string) Option {
	return func(c *Config) {
		c.ExcludePatterns = patterns
	}
}

// WithIncludePatterns sets the URL patterns to include
func WithIncludePatterns(patterns []string) Option {
	return func(c *Config) {
		c.IncludePatterns = patterns
	}
}
