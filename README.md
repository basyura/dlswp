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
- **File movement**: Moves ALL files in the target directory to today's date backup folder
- **Backup cleanup**: Removes backup directories older than specified days from today
- **Argument validation**: Only accepts 0 or positive values (negative values are rejected)

## Command Line Arguments

### Format
```
dlswp [days_to_keep] [target_directory]
```

### Argument Details

- **First argument (days_to_keep)**: Optional (defaults to "0")
  - Specify as number (0 or positive values only, negative values are rejected)
  - Number of days of backups to keep from today
  - Files are always moved to today's date backup folder
  - Only affects cleanup: removes backup folders older than specified days from today
  - **Special case**: 0 = treated as 4 days (keeps 4 days of backups)
  - **Note**: Does NOT filter files by modification date - ALL files are moved
- **Second argument (target_directory)**: Optional
  - Root directory path to process
  - Defaults to OS-specific Downloads folder

### Default Behavior (No Arguments)

When no arguments are provided, the program defaults to:
- First argument: "0" (which is treated as 4 days retention)
- Second argument: OS-specific Downloads folder

```bash
# No arguments specified
dlswp
# Equivalent to: dlswp 0 [default Downloads folder]
# Behavior: Move all files to today's backup, keep 4 days of backups
```

### Usage Examples

```bash
# No arguments: same as "dlswp 0" - move all files to today's backup, keep 4 days of backups
dlswp

# Explicit 0: move all files to __backup__/2025-08-03/, treated as 4 days (remove backups older than 2025-07-30)
dlswp 0

# Move all files to __backup__/2025-08-03/, keep 3 days (remove backups older than 2025-07-31)
# Keeps: 2025-08-01, 2025-08-02, 2025-08-03
dlswp 3

# Move all files to __backup__/2025-08-03/, keep 7 days (remove backups older than 2025-07-27)
# Keeps: 2025-07-28 through 2025-08-03
dlswp 7

# Custom directory - move all files to C:\MyFolder\__backup__\2025-08-03\
dlswp 0 "C:\MyFolder"

# macOS/Linux custom directory with 5 days retention
dlswp 5 "/path/to/directory"

# Explicit 4: same behavior as dlswp 0
dlswp 4

# Error: negative values are rejected
dlswp -1  # Shows error message and exits
```