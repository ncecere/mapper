package crawler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Result represents the outcome of crawling a URL
type Result struct {
	URL         string    // The URL that was crawled
	LastMod     time.Time // Last modification time
	StatusCode  int       // HTTP status code
	Error       error     // Any error that occurred
	Depth       int       // Depth from the start URL
	TimeToFetch time.Duration
}

// Crawler manages the web crawling process
type Crawler struct {
	config    *Config
	queue     *URLQueue
	validator *URLValidator
	client    *http.Client

	// Statistics
	stats struct {
		sync.Mutex
		processed int
		errors    int
		start     time.Time
	}

	// Channels for coordination
	results chan *Result
	done    chan struct{}
}

// NewCrawler creates a new Crawler instance
func NewCrawler(config *Config) (*Crawler, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	validator, err := NewURLValidator(
		config.BaseURL,
		config.ExcludePatterns,
		config.IncludePatterns,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create validator: %w", err)
	}

	client := &http.Client{
		Timeout: config.RequestTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if !config.FollowRedirects {
				return http.ErrUseLastResponse
			}
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			return nil
		},
	}

	c := &Crawler{
		config:    config,
		queue:     NewURLQueue(config.BaseURL),
		validator: validator,
		client:    client,
		results:   make(chan *Result),
		done:      make(chan struct{}),
	}

	return c, nil
}

// Start begins the crawling process
func (c *Crawler) Start(ctx context.Context) (<-chan *Result, error) {
	// Initialize statistics
	c.stats.start = time.Now()

	// Add the start URL to the queue
	c.queue.Push([]*url.URL{c.config.BaseURL}, 0)

	// Start worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < c.config.MaxConcurrent; i++ {
		wg.Add(1)
		go c.worker(ctx, &wg)
	}

	// Start a goroutine to close results channel when done
	go func() {
		wg.Wait()
		close(c.results)
		close(c.done)
	}()

	return c.results, nil
}

// worker processes URLs from the queue
func (c *Crawler) worker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Get next URL from queue
			item := c.queue.Pop()
			if item == nil {
				return
			}

			// Skip if URL is invalid
			if !c.validator.IsValid(item.URL) {
				continue
			}

			// Skip if beyond max depth
			if item.Depth > c.config.MaxDepth {
				continue
			}

			// Process the page
			start := time.Now()
			page := NewPage(item.URL, item.Depth)
			err := page.Process(c.client)
			duration := time.Since(start)

			// Update statistics
			c.stats.Lock()
			c.stats.processed++
			if err != nil {
				c.stats.errors++
			}
			c.stats.Unlock()

			// Send result
			c.results <- &Result{
				URL:         item.URL.String(),
				LastMod:     page.LastModified,
				StatusCode:  http.StatusOK,
				Error:       err,
				Depth:       item.Depth,
				TimeToFetch: duration,
			}

			// If page was processed successfully, add its links to the queue
			if err == nil {
				c.queue.Push(page.Links, item.Depth+1)
			}

			// Rate limiting
			if c.config.RateLimit > 0 {
				time.Sleep(c.config.RateLimit)
			}
		}
	}
}

// Wait blocks until crawling is complete
func (c *Crawler) Wait() {
	<-c.done
}

// Stats returns current crawling statistics
func (c *Crawler) Stats() (processed, errors int, duration time.Duration) {
	c.stats.Lock()
	defer c.stats.Unlock()
	return c.stats.processed, c.stats.errors, time.Since(c.stats.start)
}

// GetProcessedURLs returns all successfully processed URLs
func (c *Crawler) GetProcessedURLs() []*url.URL {
	return c.queue.GetProcessedURLs()
}
