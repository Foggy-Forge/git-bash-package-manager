package manifest

import (
	"fmt"
	"os"
	"runtime"

	"gopkg.in/yaml.v3"
)

// Manifest represents a package manifest
type Manifest struct {
	Name        string     `yaml:"name"`
	Version     string     `yaml:"version"`
	Description string     `yaml:"description,omitempty"`
	Homepage    string     `yaml:"homepage,omitempty"`
	License     string     `yaml:"license,omitempty"`
	Platforms   []Platform `yaml:"platforms"`
	Install     Install    `yaml:"install"`
}

// Platform represents a platform-specific artifact
type Platform struct {
	OS       string `yaml:"os"`
	Arch     string `yaml:"arch"`
	Archive  bool   `yaml:"archive,omitempty"`
	URL      string `yaml:"url"`
	Checksum string `yaml:"checksum,omitempty"`
}

// Install represents installation steps
type Install struct {
	Steps []InstallStep `yaml:"steps"`
}

// InstallStep represents a single installation step
type InstallStep struct {
	Type string `yaml:"type"`
	From string `yaml:"from,omitempty"`
	To   string `yaml:"to,omitempty"`
}

// LoadManifest loads and parses a manifest file
func LoadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	if err := m.Validate(); err != nil {
		return nil, fmt.Errorf("invalid manifest: %w", err)
	}

	return &m, nil
}

// Validate checks if the manifest is valid
func (m *Manifest) Validate() error {
	if m.Name == "" {
		return fmt.Errorf("name is required")
	}
	if m.Version == "" {
		return fmt.Errorf("version is required")
	}
	if len(m.Platforms) == 0 {
		return fmt.Errorf("at least one platform is required")
	}
	if len(m.Install.Steps) == 0 {
		return fmt.Errorf("at least one install step is required")
	}
	return nil
}

// GetPlatform returns the platform matching the current OS and architecture
func (m *Manifest) GetPlatform() (*Platform, error) {
	currentOS := runtime.GOOS
	currentArch := runtime.GOARCH

	for i := range m.Platforms {
		p := &m.Platforms[i]
		if p.OS == currentOS && p.Arch == currentArch {
			return p, nil
		}
	}

	return nil, fmt.Errorf("no platform found for %s/%s", currentOS, currentArch)
}
