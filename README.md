# doppel

Fast CLI tool to find and remove duplicate files. Works even when files are renamed.

**Features:**
- 🔍 Smart detection: Perceptual hashing for images, SHA-256 for other files
- 🖼️ Finds visually similar images even if they're different file sizes
- 📊 Clean table output showing file locations
- ⚡ Fast: Size-based pre-filtering optimization
- 🔄 Supports comparing two different directories

## Install

### Option 1: Download Pre-built Binary (Recommended)

Download the latest release for your platform from [releases](https://github.com/bhargav0605/dopple/releases).

**Linux/macOS:**
```bash
# Download the binary for your platform
wget https://github.com/bhargav0605/dopple/releases/download/v0.1.1/doppel-linux-amd64
# or for macOS: doppel-darwin-amd64 / doppel-darwin-arm64

# Rename and make executable
mv doppel-* doppel
chmod +x doppel

# Install to system (makes it available from anywhere)
./install.sh
```

**Windows:**
```powershell
# Download doppel-windows-amd64.exe from releases
# Rename to doppel.exe
# Run installer
.\install.ps1
```

### Option 2: Install with Go

```bash
go install github.com/bhargav0605/dopple@latest
```

### Option 3: Build from Source

```bash
git clone https://github.com/bhargav0605/dopple.git
cd dopple
go build -o doppel

# Install to system (optional)
./install.sh
```

After installation, `doppel` will be available from anywhere in your terminal.

## Usage

**Find duplicates in a single directory:**
```bash
# Interactive mode (one by one)
doppel /path/to/directory

# Show all duplicates first, then delete all with one confirmation
doppel --show-all /path/to/directory

# Preview only
doppel --dry-run /path/to/directory

# Auto-delete (keeps first file)
doppel --auto-delete /path/to/directory

# Filter by size and extension
doppel --min-size 1048576 --extensions .jpg,.png /photos
```

**Compare two different directories:**
```bash
# Find duplicates between two directories
doppel compare /path/to/dir1 /path/to/dir2

# Preview only
doppel compare --dry-run /photos /backup

# Auto-delete from directory 1
doppel compare --delete-from-1 /backup /main

# Auto-delete from directory 2
doppel compare --delete-from-2 /main /backup

# With filters
doppel compare --extensions .jpg,.png /photos /backup
doppel compare --min-size 1048576 /downloads /archive

# Interactive mode (default) - choose which directory to delete from
doppel compare /photos /backup
```

## Options

**Detection:**
- `--exact` - Use exact byte matching for all files (disable perceptual hashing for images)
- `--threshold` - Similarity threshold for images (0-64, lower = more similar, default: 5)

**Filtering:**
- `--min-size` - Ignore files smaller than size in bytes
- `--extensions` - Filter by file extensions (comma-separated, e.g., .jpg,.png)

**Actions:**
- `--dry-run` - Show duplicates without deleting
- `--auto-delete` - Keep first file, delete others automatically
- `--show-all` - Show all duplicates first, then delete all with single confirmation

## Uninstall

**Linux/macOS:**
```bash
# If installed via install.sh
rm ~/.local/bin/doppel

# If installed via go install
rm $(go env GOPATH)/bin/doppel
```

**Windows:**
```powershell
# If installed via install.ps1
Remove-Item "$env:LOCALAPPDATA\doppel\doppel.exe"
# Optional: Remove from PATH manually via System Environment Variables

# If installed via go install
Remove-Item "$(go env GOPATH)\bin\doppel.exe"
```

## How it works

**For Images:**
1. Scans directory recursively for image files (.jpg, .png, .gif, .bmp, .tiff, .webp, .heic)
2. Creates perceptual hashes (difference hash) for each image
3. Compares images using Hamming distance to find visually similar ones
4. Groups similar images (default: 92%+ similarity)

**For Other Files:**
1. Scans directory recursively
2. Groups files by size (optimization - only hash files with matching sizes)
3. Calculates SHA-256 hash for files with matching sizes
4. Groups files by hash to find exact duplicates

**Interactive Deletion:**
- View duplicates in a clean table format showing filename, location, and size
- Choose which files to keep/delete, or use auto-delete modes

## Security

This project takes security seriously. Every release includes:

- 🔍 **Dependency scanning** with `govulncheck`
- 🛡️ **Code security analysis** with `gosec`
- ✓ **SHA-256 checksums** for all binaries
- 📋 **Detailed security reports** attached to each release

### Running Security Scans Locally

```bash
# Run all security scans
make security

# Check dependencies only
make deps-check

# Check code only
make code-scan
```

### CI/CD Security

- Security scans run automatically on every PR and push
- Weekly scheduled scans to catch new vulnerabilities
- Release builds include comprehensive security reports
- Results uploaded to GitHub Security tab (SARIF)
