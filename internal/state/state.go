package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// State represents the gbpm state
type State struct {
	Installed map[string]*Package `json:"installed"`
}

// Package represents an installed package
type Package struct {
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Files       []string  `json:"files"`
	InstalledAt time.Time `json:"installed_at"`
}

// Load loads the state from the state file
func Load(statePath string) (*State, error) {
	// If file doesn't exist, return empty state
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		return &State{
			Installed: make(map[string]*Package),
		}, nil
	}

	data, err := os.ReadFile(statePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %w", err)
	}

	if s.Installed == nil {
		s.Installed = make(map[string]*Package)
	}

	return &s, nil
}

// Save saves the state to the state file
func (s *State) Save(statePath string) error {
	// Create parent directory
	if err := os.MkdirAll(filepath.Dir(statePath), 0755); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(statePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

// AddPackage adds a package to the state
func (s *State) AddPackage(pkg *Package) {
	if s.Installed == nil {
		s.Installed = make(map[string]*Package)
	}
	s.Installed[pkg.Name] = pkg
}

// RemovePackage removes a package from the state
func (s *State) RemovePackage(name string) {
	delete(s.Installed, name)
}

// GetPackage returns a package from the state
func (s *State) GetPackage(name string) (*Package, bool) {
	pkg, ok := s.Installed[name]
	return pkg, ok
}

// IsInstalled checks if a package is installed
func (s *State) IsInstalled(name string) bool {
	_, ok := s.Installed[name]
	return ok
}
