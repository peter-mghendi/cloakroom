package utility

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

// Download downloads a file from 'url' to 'destinationPath'.
// It implements several best practices:
//  1. Retries with exponential backoff.
//  2. Downloads to a temporary .partial file, then renames on success.
//  3. (Optional) Verifies the file's SHA-256 checksum if non-empty.
//  4. Tracks progress via a progress bar.
//  5. Respects context cancellation.
//
// Arguments:
//   - ctx: to allow cancellation (e.g., from signals or parent context).
//   - p: mpb.Progress pointer for multi-file progress bar handling.
//   - url: the direct download URL.
//   - destinationPath: full path of the final file on disk.
//   - expectedSHA256: if not empty, verifies the downloaded file matches this checksum (hex-encoded).
//   - maxRetries: how many times to attempt with exponential backoff.
//
// Returns an error if something goes wrong or if checksum verification fails.
func Download(
	ctx context.Context,
	progress *mpb.Progress,
	url string,
	destination string,
	hash *string,
	retries int,
) error {

	// Create the final directory if needed
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return fmt.Errorf("mkdir failed for %s: %v", filepath.Dir(destination), err)
	}

	// Partial file approach
	partialPath := destination + ".partial"

	// For the progress bar labeling
	filename := filepath.Base(destination)

	var attempt int
	var backoff time.Duration
	var lastErr error

	for attempt = 0; attempt <= retries; attempt++ {
		// If context is canceled, bail out immediately
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Begin the single download attempt
		lastErr = fetch(ctx, progress, url, partialPath, filename)
		if lastErr == nil {
			// If the download succeeded, rename the partial file => final destination
			if err := os.Rename(partialPath, destination); err != nil {
				return fmt.Errorf("rename failed: %w", err)
			}

			// If we have a checksum, verify it
			if hash != nil {
				if err := verify(destination, *hash); err != nil {
					lastErr = err
				} else {
					return nil
				}
			} else {
				return nil
			}
		}

		// If we reach here, either the download or checksum failed
		//  => we’ll retry if attempt < retries
		if attempt < retries {
			backoff = exponentialBackoff(attempt)
			fmt.Printf("[Retry %d/%d] Retrying in %s due to error: %v\n",
				attempt+1, retries, backoff, lastErr)
			time.Sleep(backoff)
		}
	}

	return fmt.Errorf("download failed after %d attempts: last error: %w", attempt, lastErr)
}

// fetch performs a single attempt at downloading the file into partialPath.
// It also creates/updates a progress bar for the read operation.
func fetch(ctx context.Context, p *mpb.Progress, url, partialPath, fileLabel string) error {
	// Remove any leftover partial file before starting fresh
	_ = os.Remove(partialPath)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP GET failed: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	// Create the partial file
	out, err := os.Create(partialPath)
	if err != nil {
		return fmt.Errorf("create partial file: %w", err)
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {

		}
	}(out)

	// Set up the bar
	totalSize := resp.ContentLength
	if totalSize < 0 {
		// If server doesn't send Content-Length, set it to 0
		totalSize = 0
	}

	// Create a bar for this file download
	bar := p.AddBar(
		totalSize,
		mpb.PrependDecorators(
			// Use a fixed width so multi-file downloads line up well
			decor.Name(fmt.Sprintf("%-25s", fileLabel)),
			decor.Percentage(),
		),
		mpb.AppendDecorators(
			decor.CountersKibiByte("% .2f / % .2f"),
			decor.AverageETA(decor.ET_STYLE_MMSS),
		),
	)

	// Wrap resp.Body with the bar’s ProxyReader
	reader := bar.ProxyReader(resp.Body)
	defer func(reader io.ReadCloser) {
		err := reader.Close()
		if err != nil {

		}
	}(reader)

	// Copy to disk
	if _, err := io.Copy(out, reader); err != nil {
		return fmt.Errorf("io copy failed: %w", err)
	}

	// If we got here, it means the download completed successfully
	return nil
}

// verify checks the SHA-256 of the downloaded file
// against the expected hex-encoded string. Returns an error if mismatched.
func verify(filePath, expectedHex string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open for checksum: %w", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return fmt.Errorf("copy for checksum: %w", err)
	}
	actualHex := hex.EncodeToString(h.Sum(nil))

	// Normalize uppercase vs. lowercase
	expectedHex = strings.ToLower(strings.TrimSpace(expectedHex))
	if actualHex != expectedHex {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedHex, actualHex)
	}
	return nil
}

// exponentialBackoff returns a simple exponential backoff duration
// e.g. attempt=0 => ~1s, attempt=1 => ~2s, attempt=2 => ~4s, etc.
func exponentialBackoff(attempt int) time.Duration {
	base := 1 * time.Second
	factor := math.Pow(2, float64(attempt))
	return time.Duration(factor) * base
}
