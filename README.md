# Mapper - Sitemap Generator

A command-line tool for generating XML sitemaps by crawling websites. The crawler stays within the specified domain and supports various configuration options.

## Features

- Concurrent web crawling with configurable limits
- Domain-scoped crawling (stays within the same domain)
- Rate limiting to prevent server overload
- Progress visualization with real-time statistics
- Configurable crawl depth and concurrency
- XML sitemap generation following the sitemap protocol
- Support for lastmod dates, change frequency, and priority

## Installation

```bash
go install github.com/bitop-dev/mapper@latest
```

Or clone and build from source:

```bash
git clone https://github.com/bitop-dev/mapper.git
cd mapper
go build
```

## Usage

Basic usage:
```bash
mapper generate https://example.com
```

With options:
```bash
mapper generate \
  --depth 3 \
  --concurrent 5 \
  --rate-limit 500ms \
  --output sitemap.xml \
  https://example.com
```

### Available Options

- `--depth, -d`: Maximum crawl depth (default: 3)
- `--output, -o`: Output file path (default: sitemap.xml)
- `--concurrent, -c`: Maximum concurrent requests (default: 5)
- `--timeout, -t`: Request timeout (default: 10s)
- `--rate-limit, -r`: Rate limit between requests (default: 1s)
- `--exclude, -e`: Paths to exclude (e.g., /admin/*)
- `--no-follow-redirects`: Don't follow redirects
- `--strip-query`: Strip query parameters from URLs (default: true)
- `--user-agent`: Custom User-Agent string
- `--debug`: Enable debug mode
- `--config`: Config file path (default: $HOME/.mapper.yaml)

### Configuration File

You can create a configuration file at `~/.mapper.yaml` with default settings:

```yaml
concurrent_requests: 5
request_timeout: 10s
rate_limit: 1s
user_agent: "Mapper/1.0"
debug: false
```

## Example Output

```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://example.com/</loc>
    <lastmod>2025-02-16</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.5</priority>
  </url>
  ...
</urlset>
```

## Development

The project follows standard Go project layout:

- `cmd/`: Command line interface
- `pkg/crawler/`: Core crawling logic
- `pkg/sitemap/`: Sitemap generation
- `pkg/ui/`: Terminal UI components

### Building from Source

1. Clone the repository
2. Install dependencies: `go mod download`
3. Build: `go build`

## License

MIT License
