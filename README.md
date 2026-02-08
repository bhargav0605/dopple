# doppel

Fast CLI tool to find and remove duplicate files using SHA-256 hashing. Works even when files are renamed.

## Install

**Download Binary:**

Download the latest release for your platform from [releases](https://github.com/yourusername/doppel/releases).

**Install with Go:**

```bash
go install github.com/yourusername/doppel@latest
```

**Build from source:**

```bash
git clone https://github.com/yourusername/doppel.git
cd doppel
go build
```

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

## How it works

1. Scans directory recursively
2. Groups files by size (optimization)
3. Hashes files with matching sizes using SHA-256
4. Groups files by hash to find duplicates
5. Lets you choose which files to keep/delete
