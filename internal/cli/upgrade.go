package cli

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

const (
	repoOwner = "Foggy-Forge"
	repoName  = "git-bash-package-manager"
)

func newUpgradeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade gbpm to the latest version",
		Long:  "Download and install the latest version of gbpm.",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Checking for updates...")

			// Get current executable path
			exePath, err := os.Executable()
			if err != nil {
				return fmt.Errorf("failed to get executable path: %w", err)
			}
			exePath, err = filepath.EvalSymlinks(exePath)
			if err != nil {
				return fmt.Errorf("failed to resolve symlinks: %w", err)
			}

			// Get latest release info
			latestVersion, err := getLatestVersion()
			if err != nil {
				return fmt.Errorf("failed to check for updates: %w", err)
			}

			if latestVersion == version {
				fmt.Printf("Already on the latest version: v%s\n", version)
				return nil
			}

			fmt.Printf("Current version: v%s\n", version)
			fmt.Printf("Latest version: %s\n", latestVersion)
			fmt.Println()

			// Construct download URL
			binaryName := getBinaryName()
			downloadURL := fmt.Sprintf(
				"https://github.com/%s/%s/releases/download/%s/%s",
				repoOwner, repoName, latestVersion, binaryName,
			)

			// Download to temp file
			fmt.Println("Downloading update...")
			tmpFile, err := os.CreateTemp("", "gbpm-upgrade-*")
			if err != nil {
				return fmt.Errorf("failed to create temp file: %w", err)
			}
			tmpPath := tmpFile.Name()
			defer os.Remove(tmpPath)

			resp, err := http.Get(downloadURL)
			if err != nil {
				return fmt.Errorf("failed to download: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("download failed: %s", resp.Status)
			}

			if _, err := io.Copy(tmpFile, resp.Body); err != nil {
				tmpFile.Close()
				return fmt.Errorf("failed to write update: %w", err)
			}
			tmpFile.Close()

			// Make executable
			if err := os.Chmod(tmpPath, 0755); err != nil {
				return fmt.Errorf("failed to set permissions: %w", err)
			}

			// Backup current binary
			backupPath := exePath + ".bak"
			if err := os.Rename(exePath, backupPath); err != nil {
				return fmt.Errorf("failed to backup current binary: %w", err)
			}

		// Move new binary into place
		if err := os.Rename(tmpPath, exePath); err != nil {
			// Restore backup on failure
			_ = os.Rename(backupPath, exePath)
			return fmt.Errorf("failed to install update: %w", err)
		}

		// Remove backup
		_ = os.Remove(backupPath)
		fmt.Printf("✓ Successfully upgraded to %s\n", latestVersion)
		fmt.Println("\nRun 'gbpm version' to verify.")
		return nil
		},
	}
}

func getLatestVersion() (string, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)
	
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Simple JSON parsing for tag_name
	bodyStr := string(body)
	tagStart := strings.Index(bodyStr, `"tag_name":`)
	if tagStart == -1 {
		return "", fmt.Errorf("could not find tag_name in response")
	}
	
	tagStart += len(`"tag_name":"`)
	tagEnd := strings.Index(bodyStr[tagStart:], `"`)
	if tagEnd == -1 {
		return "", fmt.Errorf("could not parse tag_name")
	}

	return bodyStr[tagStart : tagStart+tagEnd], nil
}

func getBinaryName() string {
	os := runtime.GOOS
	arch := runtime.GOARCH
	
	if os == "windows" {
		return fmt.Sprintf("gbpm-%s-%s.exe", os, arch)
	}
	return fmt.Sprintf("gbpm-%s-%s", os, arch)
}
