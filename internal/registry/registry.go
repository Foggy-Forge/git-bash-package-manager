package registry

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const defaultRegistry = "https://github.com/Foggy-Forge/git-bash-package-manager-registry.git"

// Registry manages the package registry
type Registry struct {
	Path string
	URL  string
}

// New creates a new registry manager
func New(path string) (*Registry, error) {
	return &Registry{
		Path: path,
		URL:  getRegistryURL(),
	}, nil
}

// Clone clones the registry repository
func (r *Registry) Clone() error {
	// Check if already cloned
	if _, err := os.Stat(filepath.Join(r.Path, ".git")); err == nil {
		return fmt.Errorf("registry already exists, use 'gbpm update' to update")
	}

	fmt.Printf("Cloning registry from %s...\n", r.URL)

	// Create parent directory
	if err := os.MkdirAll(filepath.Dir(r.Path), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Clone repository
	cmd := exec.Command("git", "clone", r.URL, r.Path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone registry: %w", err)
	}

	fmt.Println("✓ Registry cloned successfully")
	return nil
}

// Pull updates the registry repository
func (r *Registry) Pull() error {
	// Check if registry exists
	if _, err := os.Stat(filepath.Join(r.Path, ".git")); os.IsNotExist(err) {
		return r.Clone()
	}

	fmt.Println("Updating registry...")

	cmd := exec.Command("git", "-C", r.Path, "pull")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update registry: %w", err)
	}

	fmt.Println("✓ Registry updated successfully")
	return nil
}

// FindManifest finds a manifest file for a package
func (r *Registry) FindManifest(name string) (string, error) {
	// Check if registry exists
	if _, err := os.Stat(r.Path); os.IsNotExist(err) {
		return "", fmt.Errorf("registry not found, run 'gbpm update' first")
	}

	// Look for manifest
	manifestPath := filepath.Join(r.Path, "packages", name, name+".yaml")
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		return "", fmt.Errorf("package '%s' not found in registry", name)
	}

	return manifestPath, nil
}

// getRegistryURL gets the registry URL from environment or uses default
func getRegistryURL() string {
	if url := os.Getenv("GBPM_REGISTRY_URL"); url != "" {
		return url
	}
	return defaultRegistry
}
