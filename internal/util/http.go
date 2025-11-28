package util

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Download downloads a file from URL to the specified destination
func Download(url, dest string) error {
	// Create destination directory
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Create the file
	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// DownloadWithProgress downloads a file with progress indication
func DownloadWithProgress(url, dest string) error {
	// Create destination directory
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Create the file
	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Get file size
	size := resp.ContentLength

	// Create progress reader
	reader := &ProgressReader{
		Reader: resp.Body,
		Total:  size,
	}

	// Write the body to file with progress
	_, err = io.Copy(out, reader)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Println() // New line after progress
	return nil
}

// ProgressReader tracks download progress
type ProgressReader struct {
	Reader   io.Reader
	Total    int64
	Current  int64
	lastPct  int
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	pr.Current += int64(n)
	
	if pr.Total > 0 {
		pct := int(float64(pr.Current) / float64(pr.Total) * 100)
		if pct != pr.lastPct && pct%10 == 0 {
			fmt.Printf("\rDownloading... %d%%", pct)
			pr.lastPct = pct
		}
	}
	
	return n, err
}
