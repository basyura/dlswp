# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go utility called "dlswp" (Downloads Sweeper) that manages Downloads folder organization by moving files to dated backup folders and cleaning up old backups. The tool helps manage overflowing Downloads folders by backing up and cleaning files, designed for Windows systems with the default Downloads folder (`%USERPROFILE%/Downloads`). It requires manual execution via task scheduler for periodic cleaning.

## Core Functionality

The program performs two main operations:
1. **Move files to backup**: Moves files from Downloads to `__backup__/YYYY-MM-DD/` folders
2. **Clean old backups**: Removes backup directories older than 4 days

Key functions:
- `move_downloads_to_backup()`: Organizes files by moving them to dated backup folders
- `remove_old_backup()`: Cleans up backup directories older than 4 days
- `getFilePaths()`: Gets file paths, skipping directories starting with "__"
- `getDirPaths()`: Gets directory paths for backup cleanup

## Build and Run Commands

```bash
# Build the executable (creates dlswp.exe)
go build -o dlswp.exe main_dlswp.go

# Alternative build (creates main_dlswp.exe)
go build main_dlswp.go

# Run with default Downloads folder (today's files)
go run main_dlswp.go 0

# Run with days offset (negative for past days)
go run main_dlswp.go -1

# Run with custom directory
go run main_dlswp.go 0 "C:\custom\path"

# Run the built executable
dlswp.exe 0
```

## Command Line Arguments

- First argument: Days offset from today (0 = today, -1 = yesterday, etc.)
- Second argument (optional): Root directory path (defaults to %USERPROFILE%\Downloads)

## Platform and Requirements

- **Windows only**: Designed specifically for Windows systems
- **Target folder**: `%USERPROFILE%/Downloads` (default Windows Downloads folder)
- **Scheduling**: Requires external task scheduler for periodic execution (no built-in scheduling)

## Architecture Notes

- Single file Go application with no external dependencies
- Uses standard library packages: fmt, os, path/filepath, regexp, strconv, strings, time
- Date validation using regex pattern `^\d{4}-\d{2}-\d{2}$`
- File organization based on modification time
- Skips files/directories prefixed with "__" to avoid interfering with its own backup structure
- Creates `__backup__` folder structure with date-based subdirectories (YYYY-MM-DD format)
- Automatically removes backup directories older than 4 days