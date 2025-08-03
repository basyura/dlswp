# dlswp

## Overview

- Cleans up overflowing Downloads folder (`%USERPROFILE%/Downloads`) by backing up and organizing files.
- Does not run automatically - requires task scheduler or manual execution for periodic cleaning.

## Specifications

- Windows only
- Creates `__backup__` folder and moves files to date-based subdirectories
- Files or folders starting with `__` are excluded from processing
- Default target is `%USERPROFILE%/Downloads` folder
- Keeps 3 days of backups by default and removes older files

## Command Line Arguments

### Format
```
dlswp [days_offset] [target_directory]
```

### Argument Details

- **First argument (days_offset)**: Optional (defaults to "0" = today)
  - Specify as number (positive, negative, or 0)
  - Days offset from today as reference
- **Second argument (target_directory)**: Optional (defaults to `%USERPROFILE%\Downloads`)
  - Root directory path to process

### Usage Examples

```bash
# No arguments (target today's files, Downloads folder)
dlswp

# Target today's files
dlswp 0

# Target yesterday's files
dlswp -1

# Target files from 3 days ago
dlswp -3

# Specify custom directory
dlswp 0 "C:\MyFolder"

# Yesterday's files, custom directory
dlswp -1 "D:\Documents"
```