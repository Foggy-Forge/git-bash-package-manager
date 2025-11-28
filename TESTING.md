# Testing Guide: Installing tree with gbpm

This guide walks through testing gbpm by installing the `tree` command.

## Prerequisites

1. Windows with Git Bash installed
2. gbpm v0.1.0 installed (via install.sh)
3. Internet connection

## Test 1: Install from Registry (Once manifest is in registry)

```bash
# Update registry to get latest manifests
gbpm update

# Install tree
gbpm install tree

# Verify installation
tree --version
tree .

# Check installed packages
gbpm list

# Test uninstall
gbpm uninstall tree
tree --version  # Should fail - not found
```

## Test 2: Install from Local Manifest (For testing)

If you want to test before the registry manifest is published:

1. Create `tree.yaml` locally:

```yaml
name: tree
version: "1.5.2.2"
description: "Recursive directory listing command"
homepage: "https://gnuwin32.sourceforge.net/packages/tree.htm"
license: "GPL-2.0-or-later"

platforms:
  - os: windows
    arch: amd64
    archive: true
    url: "https://sourceforge.net/projects/gnuwin32/files/tree/1.5.2.2/tree-1.5.2.2-bin.zip/download"

install:
  steps:
    - type: extract
      to: "{{ .TmpDir }}/tree"
    
    - type: copy
      from: "{{ .TmpDir }}/tree/bin/tree.exe"
      to: "{{ .BinDir }}/tree.exe"
```

2. Install from file:

```bash
gbpm install --file tree.yaml
```

## Expected Behavior

### During Installation:
```
Installing tree v1.5.2.2...
Downloading from https://sourceforge.net/projects/gnuwin32/files/tree/1.5.2.2/tree-1.5.2.2-bin.zip/download...
Downloading... 10%
Downloading... 20%
...
Downloading... 100%
Step 1/2: extract
Step 2/2: copy
✓ Successfully installed tree v1.5.2.2
```

### File Locations:
- **Downloaded archive**: `~/.gbpm/cache/tree/1.5.2.2/tree-1.5.2.2-bin.zip`
- **Installed binary**: `~/.gbpm/bin/tree.exe`
- **State tracking**: `~/.gbpm/state.json`

### Verify State:
```bash
cat ~/.gbpm/state.json
```

Should show:
```json
{
  "installed": {
    "tree": {
      "name": "tree",
      "version": "1.5.2.2",
      "files": [
        "/c/Users/YourUser/.gbpm/bin/tree.exe"
      ],
      "installed_at": "2025-11-28T..."
    }
  }
}
```

## Troubleshooting

### If download fails:
- Check internet connection
- Try downloading URL manually in browser
- SourceForge URLs may redirect - this is normal

### If extraction fails:
- Check cache directory: `ls ~/.gbpm/cache/tree/1.5.2.2/`
- Verify zip file downloaded correctly
- Check file size matches expected

### If binary doesn't work:
- Verify it's in PATH: `which tree`
- Check permissions: `ls -la ~/.gbpm/bin/tree.exe`
- Try running directly: `~/.gbpm/bin/tree.exe --version`

### If PATH not set:
```bash
export PATH="$HOME/.gbpm/bin:$PATH"
# Or add to ~/.bashrc permanently
```

## Clean Up After Testing

```bash
# Uninstall tree
gbpm uninstall tree

# Clear cache (optional)
rm -rf ~/.gbpm/cache/tree

# Check nothing remains
gbpm list  # Should not show tree
which tree # Should not find it
```

## What This Tests

✅ Download from external URL (SourceForge)  
✅ ZIP extraction  
✅ Multi-level directory navigation in archive  
✅ Binary installation to bin directory  
✅ State tracking  
✅ PATH integration  
✅ Uninstall cleanup  

This is a perfect real-world test case!
