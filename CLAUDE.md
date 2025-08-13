# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

**言語指示**: このリポジトリでの作業時は、必ず日本語で回答してください。

## Project Overview

This is a Go utility called "dlswp" (Downloads Sweeper) that manages Downloads folder organization by moving files to dated backup folders and cleaning up old backups. The tool helps manage overflowing Downloads folders by backing up and cleaning files, supporting Windows, macOS, and Linux systems with their respective default Downloads folders. It requires manual execution via task scheduler for periodic cleaning.

## Core Functionality

The program performs two main operations:
1. **Move files to backup**: Moves ALL files from target directory to `__backup__/today's date/` folder (regardless of modification date)
2. **Clean old backups**: Removes backup directories older than specified days from today

Key functions:
- `getDefaultDownloadsPath()`: Gets OS-specific default Downloads folder path
- `move_downloads_to_backup()`: Moves ALL files to today's date backup folder (ignores modification date)
- `remove_old_backup()`: Cleans up backup directories older than specified days from today
- `getFilePaths()`: Gets file paths, skipping directories starting with "__"
- `getDirPaths()`: Gets directory paths for backup cleanup

## Build and Run Commands

**IMPORTANT**: After making code changes, always verify that the build succeeds and tests pass:
```bash
# Verify build succeeds
go build -o dlswp main_dlswp.go

# Run all tests to ensure functionality is correct
go test -v
```

### Windows
```bash
# Build the executable (creates dlswp.exe)
go build -o dlswp.exe main_dlswp.go

# Run the built executable
dlswp.exe 0
```

### macOS/Linux
```bash
# Build the executable (creates dlswp)
go build -o dlswp main_dlswp.go

# Run the built executable
./dlswp 0
```

### Common Commands (All Platforms)
```bash
# Alternative build (creates main_dlswp.exe on Windows, main_dlswp on others)
go build main_dlswp.go

# No arguments: same as "0" - move all files to today's backup, keep 4 days of backups
go run main_dlswp.go

# Explicit 0: move all files to today's backup, treated as 4 days retention
go run main_dlswp.go 0

# Keep 3 days: move all files to today's backup, keep last 3 days of backups
go run main_dlswp.go 3

# Run with custom directory
# Windows
go run main_dlswp.go 0 "C:\custom\path"
# macOS/Linux
go run main_dlswp.go 0 "/path/to/custom/directory"
```

## Command Line Arguments

- First argument: Days to keep (0 or positive values only, negative values rejected)
  - Optional: defaults to "0" when not specified
  - Used for cleanup: removes backup folders older than specified days from today
  - Files are always moved to today's date backup folder
  - Special case: 0 = treated as 4 days (keeps 4 days of backups)
  - Does NOT filter files by modification date - ALL files are moved
- Second argument (optional): Root directory path (defaults to OS-specific Downloads folder)

## Platform and Requirements

- **Cross-platform**: Supports Windows, macOS, and Linux
- **Default target folders**:
  - Windows: `%USERPROFILE%\Downloads`
  - macOS: `$HOME/Downloads`
  - Linux: `$HOME/Downloads`
- **Scheduling**: Requires external task scheduler for periodic execution (no built-in scheduling)

## Architecture Notes

- Single file Go application with no external dependencies
- Uses standard library packages: fmt, os, path/filepath, regexp, runtime, strconv, strings, time
- Cross-platform OS detection using `runtime.GOOS` for proper Downloads folder path resolution
- Date validation using regex pattern `^\d{4}-\d{2}-\d{2}$`
- File organization: moves ALL files to today's date folder regardless of modification time
- Platform-agnostic path handling using `filepath.Separator`
- Skips files/directories prefixed with "__" to avoid interfering with its own backup structure
- Creates `__backup__` folder structure with date-based subdirectories (YYYY-MM-DD format)
- Automatically removes backup directories older than specified days from today
- Default behavior: when no arguments provided, acts as if "0" was specified (treated as 4 days retention)
- Validates arguments: rejects negative values with error message