package paths

import (
	"os"
	"path/filepath"
)

type Paths struct {
	Home     string
	Bin      string
	Cache    string
	Registry string
}

func NewDefault() *Paths {
	home, _ := os.UserHomeDir()

	gbpmHome := getEnvOrDefault("GBPM_HOME", filepath.Join(home, ".gbpm"))

	return &Paths{
		Home:     gbpmHome,
		Bin:      filepath.Join(gbpmHome, "bin"),
		Cache:    filepath.Join(gbpmHome, "cache"),
		Registry: filepath.Join(gbpmHome, "registry"),
	}
}

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
