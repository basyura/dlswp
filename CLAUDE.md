# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go utility called "dlswp" (Downloads Sweeper) that manages Downloads folder organization by moving files to dated backup folders and cleaning up old backups.

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
# Build the executable
go build main_dlswp.go

# Run with default Downloads folder (today's files)
go run main_dlswp.go 0

# Run with days offset (negative for past days)
go run main_dlswp.go -1

# Run with custom directory
go run main_dlswp.go 0 "C:\custom\path"
```

## Command Line Arguments

- First argument: Days offset from today (0 = today, -1 = yesterday, etc.)
- Second argument (optional): Root directory path (defaults to %USERPROFILE%\Downloads)

## Architecture Notes

- Single file Go application with no external dependencies
- Uses standard library packages: fmt, io/ioutil, os, path/filepath, regexp, strconv, strings, time
- Date validation using regex pattern `^\d{4}-\d{2}-\d{2}$`
- File organization based on modification time
- Skips files/directories prefixed with "__" to avoid interfering with its own backup structure