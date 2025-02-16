# Mapper - Sitemap Generator

A command-line tool for generating XML sitemaps by crawling websites. The crawler stays within the specified domain and supports various configuration options.

## Features

- Concurrent web crawling with configurable limits
- Domain-scoped crawling (stays within the same domain)
- Rate limiting to prevent server overload
- Simple progress display with real-time statistics
- Configurable crawl depth and concurrency
- XML sitemap generation following the sitemap protocol
- Support for lastmod dates, change frequency, and priority

## Installation

### macOS
1. Download the latest release for macOS:
   ```bash
   # For Apple Silicon (M1/M2)
   curl -LO "https://github.com/ncecere/mapper/releases/latest/download/mapper_Darwin_arm64.tar.gz"
   # For Intel Macs
   curl -LO "https://github.com/ncecere/mapper/releases/latest/download/mapper_Darwin_x86_64.tar.gz"
   ```

2. Extract the archive:
   ```bash
   tar xzf mapper_Darwin_*.tar.gz
   ```

3. Make the binary executable and move it to your PATH:
   ```bash
   chmod +x mapper
   sudo mv mapper /usr/local/bin/
   ```

### Linux

1. Download the latest release for your architecture:
   ```bash
   # For 64-bit AMD/Intel (most common)
   curl -LO "https://github.com/ncecere/mapper/releases/latest/download/mapper_Linux_x86_64.tar.gz"
   # For ARM64
   curl -LO "https://github.com/ncecere/mapper/releases/latest/download/mapper_Linux_arm64.tar.gz"
   ```

2. Extract the archive:
   ```bash
   tar xzf mapper_Linux_*.tar.gz
   ```

3. Make the binary executable and move it to your PATH:
   ```bash
   chmod +x mapper
   sudo mv mapper /usr/local/bin/
   ```

### Windows

1. Download the latest release from the [releases page](https://github.com/ncecere/mapper/releases/latest)
   - Choose `mapper_Windows_x86_64.zip` for 64-bit Windows

2. Extract the ZIP file using File Explorer or PowerShell:
   ```powershell
   Expand-Archive -Path mapper_Windows_x86_64.zip -DestinationPath C:\mapper
   ```

3. Add to PATH:
   - Open Start Menu and search for "Environment Variables"
   - Click "Edit the system environment variables"
   - Click "Environment Variables" button
   - Under "System Variables", find and select "Path"
   - Click "Edit"
   - Click "New"
   - Add `C:\mapper`
   - Click "OK" on all windows

4. Verify installation (open a new PowerShell window):
   ```powershell
   mapper --version
   ```

## Project Structure

```
mapper/
├── cmd/                    # Command line interface
│   ├── root.go            # Root command setup
│   └── generate.go        # Generate command implementation
├── pkg/
│   ├── crawler/           # Web crawler package
│   │   ├── config.go      # Crawler configuration
│   │   ├── crawler.go     # Core crawler implementation
│   │   ├── page.go        # Page processing
│   │   ├── queue.go       # URL queue management
│   │   └── validator.go   # URL validation
│   ├── sitemap/           # Sitemap generation
│   │   ├── builder.go     # Sitemap construction
│   │   ├── types.go       # Data structures
│   │   └── writer.go      # XML output
│   └── ui/                # User interface
│       └── progress.go    # Progress display
├── .gitignore             # Git ignore patterns
├── go.mod                 # Go module definition
├── go.sum                 # Go module checksums
├── main.go                # Application entry point
└── README.md             # Project documentation
```

## Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/ncecere/mapper.git
   cd mapper
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the project:
   ```bash
   go build
   ```

## Git Workflow

1. Create a new branch for your changes:
   ```bash
   git checkout -b feature/your-feature
   ```

2. Make your changes and commit them:
   ```bash
   git add .
   git commit -m "Description of your changes"
   ```

3. Push your changes:
   ```bash
   git push origin feature/your-feature
   ```

## Usage Examples

### Basic Usage

Generate a sitemap for a website:
```bash
mapper generate https://example.com
```

### Advanced Options

1. Control crawling depth and concurrency:
   ```bash
   mapper generate \
     --depth 3 \
     --concurrent 5 \
     https://example.com
   ```

2. Configure rate limiting:
   ```bash
   mapper generate \
     --rate-limit 500ms \
     --timeout 10s \
     https://example.com
   ```

3. Exclude specific paths:
   ```bash
   mapper generate \
     --exclude "/admin/*" \
     --exclude "/private/*" \
     https://example.com
   ```

4. Custom output location:
   ```bash
   mapper generate \
     --output /path/to/sitemap.xml \
     https://example.com
   ```

### Configuration File

Create a `~/.mapper.yaml` file for default settings:

```yaml
concurrent_requests: 5
request_timeout: 10s
rate_limit: 1s
user_agent: "Mapper/1.0"
debug: false
```

## Output Format

The generated sitemap follows the [Sitemap Protocol](https://www.sitemaps.org/protocol.html):

```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://example.com/</loc>
    <lastmod>2025-02-16</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.5</priority>
  </url>
</urlset>
```

## Design Principles

1. **Modularity**: Each package has a specific responsibility:
   - `crawler`: Handles web crawling and URL management
   - `sitemap`: Manages sitemap generation and XML output
   - `ui`: Handles user interface and progress display

2. **Configuration**: Flexible configuration through:
   - Command-line flags
   - Configuration file
   - Environment variables

3. **Error Handling**: Comprehensive error handling with:
   - Detailed error messages
   - Debug mode for troubleshooting
   - Proper cleanup on interruption

4. **Performance**: Efficient operation through:
   - Concurrent crawling
   - Rate limiting
   - Memory-efficient URL queue
   - Duplicate URL detection

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

MIT License
