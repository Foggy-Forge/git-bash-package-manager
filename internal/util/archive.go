package util

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ExtractZip extracts a zip archive to the destination directory
func ExtractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("failed to open zip: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		if err := extractZipFile(f, dest); err != nil {
			return err
		}
	}

	return nil
}

func extractZipFile(f *zip.File, dest string) error {
	// Construct file path
	fpath := filepath.Join(dest, f.Name)

	// Check for ZipSlip vulnerability
	if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
		return fmt.Errorf("illegal file path: %s", fpath)
	}

	if f.FileInfo().IsDir() {
		// Create directory
		return os.MkdirAll(fpath, os.ModePerm)
	}

	// Create parent directory
	if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
		return err
	}

	// Create file
	outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Open file in archive
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// Copy contents
	_, err = io.Copy(outFile, rc)
	return err
}

// ExtractTarGz extracts a tar.gz archive to the destination directory
func ExtractTarGz(src, dest string) error {
	file, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open tar.gz: %w", err)
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar: %w", err)
		}

		// Construct file path
		fpath := filepath.Join(dest, header.Name)

		// Check for ZipSlip vulnerability
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return err
			}
		case tar.TypeReg:
			// Create parent directory
			if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			// Create file
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// Copy contents
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}

	return nil
}

// Extract automatically detects and extracts an archive
func Extract(src, dest string) error {
	if strings.HasSuffix(src, ".zip") {
		return ExtractZip(src, dest)
	} else if strings.HasSuffix(src, ".tar.gz") || strings.HasSuffix(src, ".tgz") {
		return ExtractTarGz(src, dest)
	}
	return fmt.Errorf("unsupported archive format: %s", src)
}
