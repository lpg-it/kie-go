package kie

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Downloader downloads generated images and videos from KIE.
//
// When configured with WithKIEClient, the Downloader automatically converts
// KIE tempfile URLs to temporary download URLs before downloading.
type Downloader struct {
	httpClient *http.Client
	concurrent int
	kieClient  *Client // Optional: for automatic download URL conversion
}

// DownloaderOption configures the Downloader.
type DownloaderOption func(*Downloader)

// NewDownloader creates a new Downloader.
func NewDownloader(opts ...DownloaderOption) *Downloader {
	d := &Downloader{
		httpClient: &http.Client{},
		concurrent: 5,
	}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

// WithDownloadConcurrency sets the number of concurrent downloads.
func WithDownloadConcurrency(n int) DownloaderOption {
	return func(d *Downloader) {
		if n > 0 {
			d.concurrent = n
		}
	}
}

// WithDownloadHTTPClient sets a custom HTTP client for downloads.
func WithDownloadHTTPClient(client *http.Client) DownloaderOption {
	return func(d *Downloader) {
		if client != nil {
			d.httpClient = client
		}
	}
}

// WithKIEClient sets a KIE client for automatic download URL conversion.
//
// When set, the Downloader will automatically call GetDownloadURL for
// KIE tempfile URLs before downloading, ensuring valid temporary links.
//
// Example:
//
//	client := kie.NewClient("your-api-key")
//	downloader := kie.NewDownloader(kie.WithKIEClient(client))
func WithKIEClient(client *Client) DownloaderOption {
	return func(d *Downloader) {
		d.kieClient = client
	}
}

// DownloadResult represents the result of a download operation.
type DownloadResult struct {
	// URL is the source URL
	URL string

	// LocalPath is the path where the file was saved
	LocalPath string

	// Size is the file size in bytes
	Size int64

	// Error is any error that occurred during download
	Error error
}

// DownloadFromTaskInfo downloads all result URLs from a TaskInfo.
//
// Example:
//
//	downloader := kie.NewDownloader()
//	results, err := downloader.DownloadFromTaskInfo(ctx, taskInfo, "./output/")
//	for _, r := range results {
//	    if r.Error != nil {
//	        log.Printf("Failed to download %s: %v", r.URL, r.Error)
//	    } else {
//	        log.Printf("Downloaded %s to %s (%d bytes)", r.URL, r.LocalPath, r.Size)
//	    }
//	}
func (d *Downloader) DownloadFromTaskInfo(ctx context.Context, info *TaskInfo, outputDir string) ([]DownloadResult, error) {
	urls, err := info.GetResultURLs()
	if err != nil {
		return nil, fmt.Errorf("failed to get result URLs: %w", err)
	}
	return d.DownloadURLs(ctx, urls, outputDir, info.TaskID)
}

// DownloadURLs downloads multiple URLs to a directory.
//
// The prefix is prepended to filenames to help identify the source.
func (d *Downloader) DownloadURLs(ctx context.Context, urls []string, outputDir string, prefix string) ([]DownloadResult, error) {
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	results := make([]DownloadResult, len(urls))
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, d.concurrent)

	for i, url := range urls {
		wg.Add(1)
		go func(idx int, u string) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				results[idx] = DownloadResult{URL: u, Error: ctx.Err()}
				return
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			}

			results[idx] = d.downloadFile(ctx, u, outputDir, prefix, idx)
		}(i, url)
	}

	wg.Wait()
	return results, nil
}

// downloadFile downloads a single file.
func (d *Downloader) downloadFile(ctx context.Context, url, outputDir, prefix string, index int) DownloadResult {
	result := DownloadResult{URL: url}

	// If kieClient is configured and URL is a KIE tempfile, convert it
	downloadURL := url
	if d.kieClient != nil && IsKIETempFileURL(url) {
		convertedURL, err := d.kieClient.GetDownloadURL(ctx, url)
		if err != nil {
			result.Error = fmt.Errorf("failed to get download URL: %w", err)
			return result
		}
		downloadURL = convertedURL
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		result.Error = fmt.Errorf("failed to create request: %w", err)
		return result
	}

	// Execute request
	resp, err := d.httpClient.Do(req)
	if err != nil {
		result.Error = fmt.Errorf("failed to download: %w", err)
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.Error = fmt.Errorf("download failed with status %d", resp.StatusCode)
		return result
	}

	// Generate filename
	filename := d.generateFilename(url, prefix, index)
	localPath := filepath.Join(outputDir, filename)

	// Create file
	f, err := os.Create(localPath)
	if err != nil {
		result.Error = fmt.Errorf("failed to create file: %w", err)
		return result
	}
	defer f.Close()

	// Copy data
	n, err := io.Copy(f, resp.Body)
	if err != nil {
		result.Error = fmt.Errorf("failed to write file: %w", err)
		return result
	}

	result.LocalPath = localPath
	result.Size = n
	return result
}

// generateFilename creates a filename from URL.
func (d *Downloader) generateFilename(url, prefix string, index int) string {
	// Try to extract extension from URL
	ext := ".bin"
	if idx := strings.LastIndex(url, "."); idx != -1 {
		potentialExt := url[idx:]
		// Only accept common extensions
		if len(potentialExt) <= 5 {
			switch strings.ToLower(potentialExt) {
			case ".png", ".jpg", ".jpeg", ".webp", ".gif", ".mp4", ".webm":
				ext = potentialExt
			}
		}
	}

	// Generate filename
	if prefix != "" {
		return fmt.Sprintf("%s_%d%s", prefix, index, ext)
	}
	return fmt.Sprintf("result_%d%s", index, ext)
}

// DownloadSingle downloads a single URL to a file.
func (d *Downloader) DownloadSingle(ctx context.Context, url, outputPath string) error {
	// Create parent directory
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	result := d.downloadFile(ctx, url, dir, "", 0)
	if result.Error != nil {
		return result.Error
	}

	// Rename to desired path
	if result.LocalPath != outputPath {
		if err := os.Rename(result.LocalPath, outputPath); err != nil {
			return fmt.Errorf("failed to rename file: %w", err)
		}
	}

	return nil
}
