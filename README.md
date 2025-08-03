# dlswp

## Overview

- Cleans up overflowing Downloads folder (`%USERPROFILE%/Downloads`) by backing up and organizing files.
- Does not run automatically - requires task scheduler or manual execution for periodic cleaning.

## Specifications

- Cross-platform: Windows, macOS, Linux
- Creates `__backup__` folder and moves files to date-based subdirectories
- Files or folders starting with `__` are excluded from processing
- Default target folders:
  - Windows: `%USERPROFILE%\Downloads`
  - macOS/Linux: `$HOME/Downloads`
- **File movement**: Moves ALL files in the target directory regardless of modification date
- **Backup cleanup**: Removes backup directories older than 4 days from the date calculated by days_offset

## Command Line Arguments

### Format
```
dlswp [days_offset] [target_directory]
```

### Argument Details

- **First argument (days_offset)**: Optional (defaults to "0" = today)
  - Specify as number (positive, negative, or 0)
  - Days offset from today used for:
    1. Backup folder name (YYYY-MM-DD format)
    2. Reference date for cleanup (removes folders older than 4 days from this date)
  - **Note**: Does NOT filter files by modification date - ALL files are moved
- **Second argument (target_directory)**: Optional
  - Root directory path to process
  - Defaults to OS-specific Downloads folder

### Usage Examples

```bash
# No arguments: move all files to __backup__/2025-08-03/, cleanup folders older than 2025-07-30
dlswp

# Move all files to __backup__/2025-08-03/, cleanup folders older than 2025-07-30
dlswp 0

# Move all files to __backup__/2025-08-02/, cleanup folders older than 2025-07-29
dlswp -1

# Move all files to __backup__/2025-07-31/, cleanup folders older than 2025-07-27
dlswp -3

# Custom directory - move all files to C:\MyFolder\__backup__\2025-08-03\
dlswp 0 "C:\MyFolder"

# macOS/Linux custom directory
dlswp -1 "/path/to/directory"
```