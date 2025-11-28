# gbpm Design

## Overview

`gbpm` is a small Go-based package manager that targets **Git Bash on Windows**.

Key ideas:

- **User-space only** – everything lives under `~/.gbpm`
- **Git-backed registry** – manifests live in a git repo (`git-bash-package-manager-registry`)
- **YAML manifests** – describe how to download and install a package
- **Single binary** – compiled with Go, no runtime dependencies

## Directories

By default:

- `GBPM_HOME` (default: `~/.gbpm`)
- `GBPM_BIN`  (`$GBPM_HOME/bin`)
- `GBPM_CACHE` (`$GBPM_HOME/cache`)
- `GBPM_REGISTRY` (`$GBPM_HOME/registry`)

Env vars can override:

- `GBPM_HOME`
- `GBPM_BIN`
- `GBPM_CACHE`
- `GBPM_REGISTRY`

## Registry

The registry is a git repo with this layout:

```text
packages/
  <name>/
    <name>.yaml
    # potential future files (README, logo, etc.)
```

The CLI:

* Clones registry into `GBPM_REGISTRY` on first `gbpm update`
* Pulls latest changes on subsequent `gbpm update`
* Resolves `gbpm install <name>` to `packages/<name>/<name>.yaml`

## Package Life Cycle

### Install

1. Resolve manifest (local file or registry).
2. Validate:

   * name, version, supported platform
3. Download asset into cache:

   * `GBPM_CACHE/<name>/<version>/<filename>`
4. Verify checksum (later).
5. Extract if needed (zip/tar.gz).
6. Copy file(s) into `GBPM_BIN`.
7. Record in `state.json`.

### Uninstall

1. Look up package in `state.json`.
2. Remove installed files from filesystem (best-effort).
3. Remove from `state.json`.

## State

State is stored in a single JSON file:

```json
{
  "installed": {
    "fzf": {
      "name": "fzf",
      "version": "0.46.1",
      "files": [
        "C:/Users/User/.gbpm/bin/fzf.exe"
      ],
      "installed_at": "2025-11-27T12:00:00Z"
    }
  }
}
```

This can be evolved to support multiple versions or rollback later.

## Implementation Notes

### Downloading

- Use Go's `net/http` for downloads
- Store in cache before extracting
- Support resume/retry (future)

### Extraction

- Use `archive/zip` for `.zip` files
- Use `archive/tar` + `compress/gzip` for `.tar.gz`
- Extract to temporary directory, then copy needed files

### Path Handling

- Convert between Windows and Unix paths as needed
- Support both forward and backward slashes
- Normalize paths for comparison

### Error Handling

- Provide clear error messages
- Clean up on failure (remove partial installs)
- Log operations for debugging (future)
