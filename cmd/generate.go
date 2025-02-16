package cmd

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/ncecere/mapper/pkg/crawler"
	"github.com/ncecere/mapper/pkg/sitemap"
	"github.com/ncecere/mapper/pkg/ui"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate [url]",
	Short: "Generate a sitemap for the specified URL",
	Long: `Generate an XML sitemap by crawling the specified website.
The crawler stays within the same domain as the provided URL.

Example:
  mapper generate https://example.com
  mapper generate --depth 3 --output sitemap.xml https://example.com`,
	Args: cobra.ExactArgs(1),
	RunE: runGenerate,
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Local flags
	generateCmd.Flags().IntP("depth", "d", 3, "maximum crawl depth")
	generateCmd.Flags().StringP("output", "o", "sitemap.xml", "output file path")
	generateCmd.Flags().IntP("concurrent", "c", 5, "maximum concurrent requests")
	generateCmd.Flags().DurationP("timeout", "t", 10*time.Second, "request timeout")
	generateCmd.Flags().DurationP("rate-limit", "r", time.Second, "rate limit between requests")
	generateCmd.Flags().StringSliceP("exclude", "e", []string{}, "paths to exclude (e.g., /admin/*)")
	generateCmd.Flags().Bool("no-follow-redirects", false, "don't follow redirects")
	generateCmd.Flags().Bool("strip-query", true, "strip query parameters from URLs")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	// Parse URL
	baseURL, err := url.Parse(args[0])
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Get flags
	depth, _ := cmd.Flags().GetInt("depth")
	outputPath, _ := cmd.Flags().GetString("output")
	concurrent, _ := cmd.Flags().GetInt("concurrent")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	rateLimit, _ := cmd.Flags().GetDuration("rate-limit")
	excludePaths, _ := cmd.Flags().GetStringSlice("exclude")
	noFollowRedirects, _ := cmd.Flags().GetBool("no-follow-redirects")
	stripQuery, _ := cmd.Flags().GetBool("strip-query")

	// Create crawler config
	config, err := crawler.DefaultConfig(baseURL.String())
	if err != nil {
		return fmt.Errorf("failed to create crawler config: %w", err)
	}

	config.MaxDepth = depth
	config.MaxConcurrent = concurrent
	config.RequestTimeout = timeout
	config.RateLimit = rateLimit
	config.FollowRedirects = !noFollowRedirects
	config.UserAgent = GetUserAgent()
	config.ExcludePatterns = excludePaths

	// Create crawler
	c, err := crawler.NewCrawler(config)
	if err != nil {
		return fmt.Errorf("failed to create crawler: %w", err)
	}

	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		<-sigChan
		fmt.Println("\nReceived interrupt signal. Shutting down...")
		cancel()
	}()

	fmt.Printf("Starting crawler for %s\n", baseURL)

	// Start crawler
	results, err := c.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed to start crawler: %w", err)
	}

	// Create sitemap builder
	builderOpts := sitemap.DefaultBuilderOptions()
	builderOpts.StripQueryParams = stripQuery
	builder := sitemap.NewBuilder(baseURL, builderOpts)

	// Create progress tracker
	progress := ui.NewProgress()

	// Process results
	var processedCount, errorCount int
	for result := range results {
		if result.Error != nil {
			errorCount++
			if GetDebugMode() {
				fmt.Printf("\nError crawling %s: %v", result.URL, result.Error)
			}
			continue
		}

		if err := builder.AddURL(result.URL, result.LastMod); err != nil && GetDebugMode() {
			fmt.Printf("\nError adding URL %s: %v", result.URL, err)
		}

		processedCount++
		progress.Update(ui.Stats{
			ProcessedURLs: processedCount,
			ErrorCount:    errorCount,
		})
	}

	// Mark progress as complete
	progress.Done()

	// Wait for crawler to finish
	c.Wait()

	// Build sitemap
	urlset, err := builder.Build()
	if err != nil {
		return fmt.Errorf("failed to build sitemap: %w", err)
	}

	// Create sitemap writer
	writer := sitemap.NewWriter(true)

	// Ensure output directory exists
	if dir := filepath.Dir(outputPath); dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// Write sitemap to file
	if err := writer.WriteToFile(urlset, outputPath); err != nil {
		return fmt.Errorf("failed to write sitemap: %w", err)
	}

	// Print summary
	fmt.Printf("\nSitemap generated successfully:\n")
	fmt.Printf("- URLs processed: %d\n", processedCount)
	fmt.Printf("- Errors: %d\n", errorCount)
	fmt.Printf("- Output file: %s\n", outputPath)

	return nil
}
