# gbpm – Git Bash Package Manager

`gbpm` is a lightweight, user-space package manager designed specifically for **Git Bash on Windows**.

It installs CLI tools and scripts into your home directory without requiring admin rights, using a simple, git-backed registry of YAML manifests.

---

## Goals

- **User-space only**: No admin rights, no touching `C:\Program Files`.
- **Git Bash native**: Designed for the Git for Windows / MinTTY environment.
- **Simple manifests**: YAML package descriptions stored in a git registry.
- **Self-contained**: Single Go binary + git + curl (or built-in HTTP).

---

## High-Level Design

- Installs into: `~/.gbpm`
  - Binaries: `~/.gbpm/bin`
  - Cache: `~/.gbpm/cache`
  - Registry clone: `~/.gbpm/registry`
  - State: `~/.gbpm/state.json`
- Registry = separate GitHub repo (`Foggy-Forge/git-bash-package-manager-registry`) with YAML manifests:
  - `packages/<name>/<name>.yaml`

---

## Core Commands (planned)

- `gbpm version` – show version
- `gbpm doctor` – check environment, PATH, and directories
- `gbpm paths` – print gbpm paths
- `gbpm install <name>` – install from registry
- `gbpm install --file <manifest.yaml>` – install from local manifest
- `gbpm list` – list installed packages
- `gbpm uninstall <name>` – uninstall a package
- `gbpm update` – update registry (git pull)
- `gbpm upgrade` – upgrade installed packages (later)
- `gbpm info <name>` – show manifest info

---

## Installation

### Quick Install (Recommended)

For Git Bash on Windows or any Unix-like system:

```bash
curl -fsSL https://raw.githubusercontent.com/Foggy-Forge/git-bash-package-manager/main/install.sh | bash
```

The installer will:
- Download the latest release for your platform
- Install to `~/.gbpm/bin/gbpm`
- Optionally add `~/.gbpm/bin` to your PATH

### Manual Installation

1. Download the latest release from [Releases](https://github.com/Foggy-Forge/git-bash-package-manager/releases)
2. Extract and move the binary to a directory in your PATH
3. Make it executable: `chmod +x gbpm`

### Building from Source

```bash
git clone https://github.com/Foggy-Forge/git-bash-package-manager.git
cd git-bash-package-manager
go mod tidy
go build -o gbpm ./cmd/gbpm

# optional: add local build to PATH
export PATH="$(pwd):$PATH"

gbpm version
gbpm doctor
```

---

## Roadmap

### Milestone 0 – Skeleton CLI

* [x] Basic command scaffold (`version`, `doctor`, `paths`)
* [x] Paths abstraction (`~/.gbpm`, `bin`, `cache`, `registry`)
* [x] `go build` + minimal docs

### Milestone 1 – Install from Local Manifest

* [ ] Define manifest spec (see `docs/manifest-spec.md`)
* [ ] Implement `gbpm install --file manifest.yaml`

  * [ ] YAML parsing
  * [ ] Download to cache
  * [ ] Extract (zip/tar.gz)
  * [ ] Copy binary to `~/.gbpm/bin`
  * [ ] Update `state.json`

### Milestone 2 – Git Registry

* [x] Create `git-bash-package-manager-registry` repo
* [ ] Add example manifests (`fzf`, `bat`, `rg`)
* [ ] Implement `gbpm update` (clone/pull)
* [ ] Implement `gbpm install <name>` via registry lookup
* [ ] Implement `gbpm list` and `gbpm uninstall`

See more details in [`docs/design.md`](./docs/design.md) and [`docs/manifest-spec.md`](./docs/manifest-spec.md).

---

## License

MIT License - see [LICENSE](LICENSE) for details.
