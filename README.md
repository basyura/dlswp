# dlswp

## Overview

- Cleans up overflowing Downloads folder (`%USERPROFILE%/Downloads`) by backing up and organizing files.
- Does not run automatically - requires task scheduler or manual execution for periodic cleaning.

## Specifications

- Windows only
- Creates `__backup__` folder and moves files to date-based subdirectories
- Files or folders starting with `__` are excluded from processing
- Default target is `%USERPROFILE%/Downloads` folder
- **File movement**: Moves ALL files regardless of modification date or days_offset value
- **Backup cleanup**: Removes backup directories older than 4 days from the specified date (respects days_offset)

## Command Line Arguments

### Format
```
dlswp [days_offset] [target_directory]
```

### Argument Details

- **First argument (days_offset)**: Optional (defaults to "0" = today)
  - Specify as number (positive, negative, or 0)
  - Days offset from today as reference
  - **Note**: Only affects backup cleanup timing, NOT file movement criteria
- **Second argument (target_directory)**: Optional (defaults to `%USERPROFILE%\Downloads`)
  - Root directory path to process

### Usage Examples

```bash
# No arguments (move all files, cleanup from today's perspective)
dlswp

# Move all files, cleanup from today's perspective
dlswp 0

# Move all files, cleanup from yesterday's perspective (keeps 4 days from yesterday)
dlswp -1

# Move all files, cleanup from 3 days ago perspective
dlswp -3

# Specify custom directory
dlswp 0 "C:\MyFolder"

# Custom directory with different cleanup reference date
dlswp -1 "D:\Documents"
```