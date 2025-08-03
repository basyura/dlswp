# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go utility called "dlswp" (Downloads Sweeper) that manages Downloads folder organization by moving files to dated backup folders and cleaning up old backups. The tool helps manage overflowing Downloads folders by backing up and cleaning files, supporting Windows, macOS, and Linux systems with their respective default Downloads folders. It requires manual execution via task scheduler for periodic cleaning.

## Core Functionality

The program performs two main operations:
1. **Move files to backup**: Moves files from Downloads to `__backup__/YYYY-MM-DD/` folders
2. **Clean old backups**: Removes backup directories older than 4 days

Key functions:
- `getDefaultDownloadsPath()`: Gets OS-specific default Downloads folder path
- `move_downloads_to_backup()`: Organizes files by moving them to dated backup folders
- `remove_old_backup()`: Cleans up backup directories older than 4 days
- `getFilePaths()`: Gets file paths, skipping directories starting with "__"
- `getDirPaths()`: Gets directory paths for backup cleanup

## Build and Run Commands

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

# Run with default Downloads folder (today's files)
go run main_dlswp.go 0

# Run with days offset (negative for past days)
go run main_dlswp.go -1

# Run with custom directory
# Windows
go run main_dlswp.go 0 "C:\custom\path"
# macOS/Linux
go run main_dlswp.go 0 "/path/to/custom/directory"
```

## Command Line Arguments

- First argument: Days offset from today (0 = today, -1 = yesterday, etc.)
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
- File organization based on modification time
- Platform-agnostic path handling using `filepath.Separator`
- Skips files/directories prefixed with "__" to avoid interfering with its own backup structure
- Creates `__backup__` folder structure with date-based subdirectories (YYYY-MM-DD format)
- Automatically removes backup directories older than 4 days