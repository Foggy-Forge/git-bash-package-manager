# gbpm Manifest Specification (Draft)

Manifests are YAML files that describe how to install a package.

## Top-Level Fields

```yaml
name: fzf
version: 0.46.1
description: "A general-purpose command-line fuzzy finder"
homepage: "https://github.com/junegunn/fzf"

license: "MIT"

platforms:
  - os: windows
    arch: amd64
    archive: true
    url: "https://github.com/junegunn/fzf/releases/download/v0.46.1/fzf-0.46.1-windows_amd64.zip"
    checksum: "sha256:deadbeef..."  # optional at first

install:
  steps:
    - type: extract
      to: "{{ .TmpDir }}/fzf"
    - type: copy
      from: "{{ .TmpDir }}/fzf/fzf.exe"
      to: "{{ .BinDir }}/fzf.exe"
```

## Fields

* `name` (string, required)
  Unique package name (used on CLI).

* `version` (string, required)
  Package version (semantic version recommended, but not enforced).

* `description` (string, optional)

* `homepage` (string, optional)

* `license` (string, optional)

## `platforms`

A list of platform-specific artifacts.

```yaml
platforms:
  - os: windows
    arch: amd64
    archive: true
    url: "https://..."
    checksum: "sha256:..."
```

Fields:

* `os` (string, required) — e.g. `windows`
* `arch` (string, required) — e.g. `amd64`
* `archive` (bool, optional, default: false)
  If true, asset is an archive (zip/tar.gz).
* `url` (string, required)
  Download URL for the asset.
* `checksum` (string, optional)
  Format: `algo:value`, e.g. `sha256:deadbeef...`.

The CLI picks the first entry matching the current `GOOS`/`GOARCH`.

## `install.steps`

An ordered list of actions executed to install the package.

Supported `type`s (initial version):

### `extract`

Extracts an archive to a directory.

```yaml
- type: extract
  to: "{{ .TmpDir }}/fzf"
```

* Uses the downloaded asset.
* Supports zip and tar.gz (initially zip is enough for Windows).

### `copy`

Copies a single file.

```yaml
- type: copy
  from: "{{ .TmpDir }}/fzf/fzf.exe"
  to: "{{ .BinDir }}/fzf.exe"
```

Template variables:

* `{{ .TmpDir }}` — a temp directory for this install.
* `{{ .BinDir }}` — resolved bin directory (e.g. `~/.gbpm/bin`).
* (Later) `{{ .Home }}`, `{{ .CacheDir }}`.

More step types can be added later (e.g. `chmod`, `shell`, `rename`).

## Example Manifests

### fzf

```yaml
name: fzf
version: 0.46.1
description: "A general-purpose command-line fuzzy finder"
homepage: "https://github.com/junegunn/fzf"
license: "MIT"

platforms:
  - os: windows
    arch: amd64
    archive: true
    url: "https://github.com/junegunn/fzf/releases/download/v0.46.1/fzf-0.46.1-windows_amd64.zip"

install:
  steps:
    - type: extract
      to: "{{ .TmpDir }}/fzf"
    - type: copy
      from: "{{ .TmpDir }}/fzf/fzf.exe"
      to: "{{ .BinDir }}/fzf.exe"
```

### bat

```yaml
name: bat
version: 0.24.0
description: "A cat clone with syntax highlighting"
homepage: "https://github.com/sharkdp/bat"
license: "MIT/Apache-2.0"

platforms:
  - os: windows
    arch: amd64
    archive: true
    url: "https://github.com/sharkdp/bat/releases/download/v0.24.0/bat-v0.24.0-x86_64-pc-windows-msvc.zip"

install:
  steps:
    - type: extract
      to: "{{ .TmpDir }}/bat"
    - type: copy
      from: "{{ .TmpDir }}/bat/bat-v0.24.0-x86_64-pc-windows-msvc/bat.exe"
      to: "{{ .BinDir }}/bat.exe"
```

### ripgrep

```yaml
name: ripgrep
version: 14.1.0
description: "A line-oriented search tool that recursively searches the current directory for a regex pattern"
homepage: "https://github.com/BurntSushi/ripgrep"
license: "MIT/Unlicense"

platforms:
  - os: windows
    arch: amd64
    archive: true
    url: "https://github.com/BurntSushi/ripgrep/releases/download/14.1.0/ripgrep-14.1.0-x86_64-pc-windows-msvc.zip"

install:
  steps:
    - type: extract
      to: "{{ .TmpDir }}/ripgrep"
    - type: copy
      from: "{{ .TmpDir }}/ripgrep/ripgrep-14.1.0-x86_64-pc-windows-msvc/rg.exe"
      to: "{{ .BinDir }}/rg.exe"
```
