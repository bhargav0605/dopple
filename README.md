# doppel

Fast CLI tool to find and remove duplicate files using SHA-256 hashing. Works even when files are renamed.

## Install

### Option 1: Download Pre-built Binary (Recommended)

Download the latest release for your platform from [releases](https://github.com/bhargav0605/dopple/releases).

**Linux/macOS:**
```bash
# Download the binary for your platform
wget https://github.com/bhargav0605/dopple/releases/download/v0.1.0/doppel-linux-amd64
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

```bash
# Interactive mode
doppel /path/to/directory

# Preview only
doppel --dry-run /path/to/directory

# Auto-delete (keeps first file)
doppel --auto-delete /path/to/directory

# Filter by size and extension
doppel --min-size 1048576 --extensions .jpg,.png /photos
```

## Options

- `--dry-run` - Show duplicates without deleting
- `--auto-delete` - Keep first file, delete others automatically
- `--min-size` - Ignore files smaller than size in bytes
- `--extensions` - Filter by file extensions (comma-separated)

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

1. Scans directory recursively
2. Groups files by size (optimization)
3. Hashes files with matching sizes using SHA-256
4. Groups files by hash to find duplicates
5. Lets you choose which files to keep/delete
