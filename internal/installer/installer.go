package installer

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/Foggy-Forge/git-bash-package-manager/internal/manifest"
	"github.com/Foggy-Forge/git-bash-package-manager/internal/paths"
	"github.com/Foggy-Forge/git-bash-package-manager/internal/state"
	"github.com/Foggy-Forge/git-bash-package-manager/internal/util"
)

// Installer handles package installation
type Installer struct {
	Paths     *paths.Paths
	State     *state.State
	StatePath string
}

// New creates a new installer
func New(p *paths.Paths, statePath string) (*Installer, error) {
	s, err := state.Load(statePath)
	if err != nil {
		return nil, err
	}

	return &Installer{
		Paths:     p,
		State:     s,
		StatePath: statePath,
	}, nil
}

// Install installs a package from a manifest
func (i *Installer) Install(m *manifest.Manifest) error {
	fmt.Printf("Installing %s v%s...\n", m.Name, m.Version)

	// Check if already installed
	if i.State.IsInstalled(m.Name) {
		existing, _ := i.State.GetPackage(m.Name)
		if existing.Version == m.Version {
			return fmt.Errorf("package %s v%s is already installed", m.Name, m.Version)
		}
		fmt.Printf("Upgrading from v%s to v%s\n", existing.Version, m.Version)
	}

	// Get platform
	platform, err := m.GetPlatform()
	if err != nil {
		return err
	}

	// Download asset
	cacheDir := filepath.Join(i.Paths.Cache, m.Name, m.Version)
	
	// Determine filename - use archive extension if platform says it's an archive
	filename := filepath.Base(platform.URL)
	if platform.Archive && !hasArchiveExtension(filename) {
		// SourceForge and similar may not have extension in URL
		// Default to .zip for Windows archives
		filename = m.Name + "-" + m.Version + ".zip"
	}
	
	cachePath := filepath.Join(cacheDir, filename)

	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		fmt.Printf("Downloading from %s...\n", platform.URL)
		if err := util.DownloadWithProgress(platform.URL, cachePath); err != nil {
			return fmt.Errorf("failed to download: %w", err)
		}
	} else {
		fmt.Println("Using cached download...")
	}

	// Create temp directory for extraction
	tmpDir, err := os.MkdirTemp("", "gbpm-install-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Template context
	ctx := map[string]string{
		"TmpDir":   tmpDir,
		"BinDir":   i.Paths.Bin,
		"Home":     i.Paths.Home,
		"CacheDir": i.Paths.Cache,
	}

	// Track installed files
	var installedFiles []string

	// Execute install steps
	for idx, step := range m.Install.Steps {
		fmt.Printf("Step %d/%d: %s\n", idx+1, len(m.Install.Steps), step.Type)
		
		switch step.Type {
		case "extract":
			to, err := renderTemplate(step.To, ctx)
			if err != nil {
				return fmt.Errorf("failed to render template: %w", err)
			}
			
			if platform.Archive {
				if err := util.Extract(cachePath, to); err != nil {
					return fmt.Errorf("failed to extract: %w", err)
				}
			}

		case "copy":
			from, err := renderTemplate(step.From, ctx)
			if err != nil {
				return fmt.Errorf("failed to render from template: %w", err)
			}
			
			to, err := renderTemplate(step.To, ctx)
			if err != nil {
				return fmt.Errorf("failed to render to template: %w", err)
			}

			if err := copyFile(from, to); err != nil {
				return fmt.Errorf("failed to copy: %w", err)
			}
			
			installedFiles = append(installedFiles, to)

		default:
			return fmt.Errorf("unknown step type: %s", step.Type)
		}
	}

	// Update state
	pkg := &state.Package{
		Name:        m.Name,
		Version:     m.Version,
		Files:       installedFiles,
		InstalledAt: time.Now(),
	}
	i.State.AddPackage(pkg)

	if err := i.State.Save(i.StatePath); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	fmt.Printf("✓ Successfully installed %s v%s\n", m.Name, m.Version)
	return nil
}

// Uninstall uninstalls a package
func (i *Installer) Uninstall(name string) error {
	pkg, ok := i.State.GetPackage(name)
	if !ok {
		return fmt.Errorf("package %s is not installed", name)
	}

	fmt.Printf("Uninstalling %s v%s...\n", pkg.Name, pkg.Version)

	// Remove files
	for _, file := range pkg.Files {
		fmt.Printf("Removing %s\n", file)
		if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
			fmt.Printf("Warning: failed to remove %s: %v\n", file, err)
		}
	}

	// Remove from state
	i.State.RemovePackage(name)
	if err := i.State.Save(i.StatePath); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	fmt.Printf("✓ Successfully uninstalled %s\n", name)
	return nil
}

// renderTemplate renders a template string with the given context
func renderTemplate(tmpl string, ctx map[string]string) (string, error) {
	t, err := template.New("step").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, ctx); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	// Create destination directory
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source: %w", err)
	}
	defer srcFile.Close()

	// Create destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination: %w", err)
	}
	defer dstFile.Close()

	// Copy contents
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy: %w", err)
	}

	// Copy permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source: %w", err)
	}
	
	if err := os.Chmod(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	return nil
}

// hasArchiveExtension checks if a filename has a known archive extension
func hasArchiveExtension(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".zip" || ext == ".tar" || ext == ".gz" || ext == ".tgz"
}
