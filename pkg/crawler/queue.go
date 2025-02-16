package crawler

import (
	"net/url"
	"sync"
)

// URLQueue manages the queue of URLs to be crawled
type URLQueue struct {
	mu sync.Mutex

	// queue holds the URLs to be processed
	queue []*QueueItem

	// seen tracks URLs that have been seen to prevent duplicates
	seen map[string]bool

	// baseHost is the host of the base URL to ensure we stay within domain
	baseHost string
}

// QueueItem represents a URL in the queue with its depth
type QueueItem struct {
	URL   *url.URL
	Depth int
}

// NewURLQueue creates a new URLQueue instance
func NewURLQueue(baseURL *url.URL) *URLQueue {
	return &URLQueue{
		queue:    make([]*QueueItem, 0),
		seen:     make(map[string]bool),
		baseHost: baseURL.Host,
	}
}

// Push adds a URL to the queue if it hasn't been seen and matches criteria
func (q *URLQueue) Push(urls []*url.URL, depth int) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for _, u := range urls {
		// Skip if URL has been seen
		if q.seen[u.String()] {
			continue
		}

		// Skip if URL is not in the same domain
		if u.Host != q.baseHost {
			continue
		}

		// Mark URL as seen and add to queue
		q.seen[u.String()] = true
		q.queue = append(q.queue, &QueueItem{
			URL:   u,
			Depth: depth,
		})
	}
}

// Pop removes and returns the next URL from the queue
func (q *URLQueue) Pop() *QueueItem {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.queue) == 0 {
		return nil
	}

	item := q.queue[0]
	q.queue = q.queue[1:]
	return item
}

// Len returns the current length of the queue
func (q *URLQueue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.queue)
}

// HasSeen returns true if the URL has been seen
func (q *URLQueue) HasSeen(u *url.URL) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.seen[u.String()]
}

// SeenCount returns the number of unique URLs seen
func (q *URLQueue) SeenCount() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.seen)
}

// Clear empties the queue and seen URLs
func (q *URLQueue) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.queue = make([]*QueueItem, 0)
	q.seen = make(map[string]bool)
}

// GetProcessedURLs returns a slice of all processed URLs
func (q *URLQueue) GetProcessedURLs() []*url.URL {
	q.mu.Lock()
	defer q.mu.Unlock()

	urls := make([]*url.URL, 0, len(q.seen))
	for urlStr := range q.seen {
		if u, err := url.Parse(urlStr); err == nil {
			urls = append(urls, u)
		}
	}
	return urls
}

// IsInDomain checks if a URL is in the same domain as the base URL
func (q *URLQueue) IsInDomain(u *url.URL) bool {
	return u.Host == q.baseHost
}
