package ui

import (
	"fmt"
	"sync"
	"time"
)

// Stats holds crawling statistics
type Stats struct {
	TotalURLs     int
	ProcessedURLs int
	ErrorCount    int
	StartTime     time.Time
}

// Progress manages crawling progress display
type Progress struct {
	mu        sync.Mutex
	stats     Stats
	startTime time.Time
}

// NewProgress creates a new progress tracker
func NewProgress() *Progress {
	return &Progress{
		startTime: time.Now(),
	}
}

// Update updates the crawling statistics and displays progress
func (p *Progress) Update(stats Stats) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.stats = stats
	p.display()
}

// display shows the current progress
func (p *Progress) display() {
	duration := time.Since(p.startTime).Round(time.Second)
	fmt.Printf("\rProcessed: %d URLs • Errors: %d • Time: %v     ", // Extra spaces to clear any previous output
		p.stats.ProcessedURLs,
		p.stats.ErrorCount,
		duration)
}

// Done marks the progress as complete
func (p *Progress) Done() {
	fmt.Println() // Add newline after progress output
}
