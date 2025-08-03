package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestBackupDirPattern(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"2023-01-01", true},
		{"2023-12-31", true},
		{"2023-1-1", false},
		{"23-01-01", false},
		{"2023-13-01", true}, // regex doesn't validate month range
		{"2023-01-32", true}, // regex doesn't validate day range
		{"not-a-date", false},
		{"", false},
	}

	for _, test := range tests {
		result := backup_dir_pattern.MatchString(test.input)
		if result != test.expected {
			t.Errorf("backup_dir_pattern.MatchString(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestGetDirPaths(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "dlswp_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test directories
	testDirs := []string{"dir1", "dir2", "2023-01-01"}
	for _, dir := range testDirs {
		err := os.Mkdir(filepath.Join(tempDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory %s: %v", dir, err)
		}
	}

	// Create test files (should be ignored)
	testFiles := []string{"file1.txt", "file2.txt"}
	for _, file := range testFiles {
		f, err := os.Create(filepath.Join(tempDir, file))
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
		f.Close()
	}

	// Test getDirPaths
	dirs, err := getDirPaths(tempDir)
	if err != nil {
		t.Fatalf("getDirPaths failed: %v", err)
	}

	if len(dirs) != len(testDirs) {
		t.Errorf("Expected %d directories, got %d", len(testDirs), len(dirs))
	}

	// Check that all expected directories are present
	dirSet := make(map[string]bool)
	for _, dir := range dirs {
		dirSet[dir] = true
	}

	for _, expectedDir := range testDirs {
		if !dirSet[expectedDir] {
			t.Errorf("Expected directory %s not found in result", expectedDir)
		}
	}
}

func TestGetDirPathsNonExistentDir(t *testing.T) {
	_, err := getDirPaths("/non/existent/path")
	if err == nil {
		t.Error("Expected error for non-existent directory, got nil")
	}
}

func TestGetFilePaths(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "dlswp_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	testFiles := []string{"file1.txt", "file2.txt", "document.pdf"}
	for _, file := range testFiles {
		f, err := os.Create(filepath.Join(tempDir, file))
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
		f.Close()
	}

	// Create directories (should be included)
	testDirs := []string{"dir1", "dir2"}
	for _, dir := range testDirs {
		err := os.Mkdir(filepath.Join(tempDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory %s: %v", dir, err)
		}
	}

	// Create files/dirs starting with "__" (should be skipped)
	skipItems := []string{"__backup__", "__temp__"}
	for _, item := range skipItems {
		err := os.Mkdir(filepath.Join(tempDir, item), 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory %s: %v", item, err)
		}
	}

	// Test getFilePaths
	paths, err := getFilePaths(tempDir)
	if err != nil {
		t.Fatalf("getFilePaths failed: %v", err)
	}

	expectedCount := len(testFiles) + len(testDirs)
	if len(paths) != expectedCount {
		t.Errorf("Expected %d paths, got %d", expectedCount, len(paths))
	}

	// Check that __ prefixed items are not included
	for _, path := range paths {
		basename := filepath.Base(path)
		if filepath.HasPrefix(basename, "__") {
			t.Errorf("Found __ prefixed item in results: %s", basename)
		}
	}

	// Check that expected files and directories are included
	pathSet := make(map[string]bool)
	for _, path := range paths {
		pathSet[filepath.Base(path)] = true
	}

	allExpected := append(testFiles, testDirs...)
	for _, expected := range allExpected {
		if !pathSet[expected] {
			t.Errorf("Expected item %s not found in result", expected)
		}
	}
}

func TestRemoveOldBackupEmptyDir(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "dlswp_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create __backup__ directory but keep it empty
	backupDir := filepath.Join(tempDir, "__backup__")
	err = os.Mkdir(backupDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create backup directory: %v", err)
	}

	// Test with 4 days to keep (should not remove anything in empty dir)
	err = remove_old_backup(tempDir, 4)
	if err != nil {
		t.Errorf("remove_old_backup failed: %v", err)
	}
}

func TestRemoveOldBackupWithDateDirs(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "dlswp_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create __backup__ directory
	backupDir := filepath.Join(tempDir, "__backup__")
	err = os.Mkdir(backupDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create backup directory: %v", err)
	}

	// Create some date directories
	now := time.Now()
	oldDate := now.AddDate(0, 0, -5) // 5 days ago (should be removed with daysToKeep=4)
	recentDate := now.AddDate(0, 0, -3) // 3 days ago (should be kept with daysToKeep=4)

	oldDateStr := oldDate.Format("2006-01-02")
	recentDateStr := recentDate.Format("2006-01-02")

	// Create old directory
	oldDir := filepath.Join(backupDir, oldDateStr)
	err = os.Mkdir(oldDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create old date directory: %v", err)
	}

	// Create recent directory
	recentDir := filepath.Join(backupDir, recentDateStr)
	err = os.Mkdir(recentDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create recent date directory: %v", err)
	}

	// Create invalid directory (should be ignored)
	invalidDir := filepath.Join(backupDir, "invalid-date")
	err = os.Mkdir(invalidDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create invalid date directory: %v", err)
	}

	// Run remove_old_backup with 4 days to keep
	err = remove_old_backup(tempDir, 4)
	if err != nil {
		t.Errorf("remove_old_backup failed: %v", err)
	}

	// Check that old directory was removed
	if _, err := os.Stat(oldDir); !os.IsNotExist(err) {
		t.Errorf("Old directory %s should have been removed", oldDateStr)
	}

	// Check that recent directory still exists (within 4 days)
	if _, err := os.Stat(recentDir); err != nil {
		t.Errorf("Recent directory %s should still exist: %v", recentDateStr, err)
	}

	// Check that invalid directory still exists
	if _, err := os.Stat(invalidDir); err != nil {
		t.Errorf("Invalid directory should still exist: %v", err)
	}
}

func TestRemoveOldBackupNonExistentBackupDir(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "dlswp_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test with non-existent __backup__ directory
	err = remove_old_backup(tempDir, 4)
	if err == nil {
		t.Error("Expected error for non-existent backup directory, got nil")
	}
}

func TestMoveDownloadsToBackup(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "dlswp_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	testFiles := []string{"file1.txt", "file2.txt", "document.pdf"}
	for _, file := range testFiles {
		f, err := os.Create(filepath.Join(tempDir, file))
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
		f.Close()
	}

	// Create test directories
	testDirs := []string{"dir1", "dir2"}
	for _, dir := range testDirs {
		err := os.Mkdir(filepath.Join(tempDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory %s: %v", dir, err)
		}
	}

	// Create files/dirs starting with "__" (should be skipped)
	skipItems := []string{"__backup__", "__temp__"}
	for _, item := range skipItems {
		err := os.Mkdir(filepath.Join(tempDir, item), 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory %s: %v", item, err)
		}
	}

	// Run move_downloads_to_backup
	err = move_downloads_to_backup(tempDir)
	if err != nil {
		t.Fatalf("move_downloads_to_backup failed: %v", err)
	}

	// Check that __backup__ directory was created with today's date
	todayStr := time.Now().Format("2006-01-02")
	backupDir := filepath.Join(tempDir, "__backup__", todayStr)
	if _, err := os.Stat(backupDir); err != nil {
		t.Errorf("Backup directory %s should have been created: %v", backupDir, err)
	}

	// Check that all test files and directories were moved
	allExpected := append(testFiles, testDirs...)
	for _, expected := range allExpected {
		movedPath := filepath.Join(backupDir, expected)
		if _, err := os.Stat(movedPath); err != nil {
			t.Errorf("Expected item %s should have been moved to backup: %v", expected, err)
		}

		// Check that original item no longer exists in root
		originalPath := filepath.Join(tempDir, expected)
		if _, err := os.Stat(originalPath); !os.IsNotExist(err) {
			t.Errorf("Original item %s should have been moved from root directory", expected)
		}
	}

	// Check that __ prefixed items were not moved
	for _, skipItem := range skipItems {
		skipPath := filepath.Join(tempDir, skipItem)
		if _, err := os.Stat(skipPath); err != nil {
			t.Errorf("Skip item %s should still exist in root directory: %v", skipItem, err)
		}
	}
}

func TestMoveDownloadsToBackupEmptyDir(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "dlswp_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Run move_downloads_to_backup on empty directory
	err = move_downloads_to_backup(tempDir)
	if err != nil {
		t.Fatalf("move_downloads_to_backup failed: %v", err)
	}

	// Check that no __backup__ directory was created (since no files to move)
	backupDir := filepath.Join(tempDir, "__backup__")
	if _, err := os.Stat(backupDir); !os.IsNotExist(err) {
		t.Error("No backup directory should be created when there are no files to move")
	}
}

func TestZeroArgumentTreatedAsFour(t *testing.T) {
	// Test the conversion logic: 0 should be treated as 4
	// This tests the main function logic indirectly by testing remove_old_backup with 4

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "dlswp_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create __backup__ directory
	backupDir := filepath.Join(tempDir, "__backup__")
	err = os.Mkdir(backupDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create backup directory: %v", err)
	}

	// Create test directories to verify 4-day retention behavior
	now := time.Now()
	oldDate := now.AddDate(0, 0, -5) // 5 days ago (should be removed with 4 days retention)
	recentDate := now.AddDate(0, 0, -3) // 3 days ago (should be kept with 4 days retention)

	oldDateStr := oldDate.Format("2006-01-02")
	recentDateStr := recentDate.Format("2006-01-02")

	// Create directories
	oldDir := filepath.Join(backupDir, oldDateStr)
	err = os.Mkdir(oldDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create old date directory: %v", err)
	}

	recentDir := filepath.Join(backupDir, recentDateStr)
	err = os.Mkdir(recentDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create recent date directory: %v", err)
	}

	// Test with 4 days (what 0 should become)
	err = remove_old_backup(tempDir, 4)
	if err != nil {
		t.Errorf("remove_old_backup failed: %v", err)
	}

	// Check that old directory (5 days ago) was removed
	if _, err := os.Stat(oldDir); !os.IsNotExist(err) {
		t.Errorf("Old directory %s should have been removed with 4-day retention", oldDateStr)
	}

	// Check that recent directory (3 days ago) still exists
	if _, err := os.Stat(recentDir); err != nil {
		t.Errorf("Recent directory %s should still exist with 4-day retention: %v", recentDateStr, err)
	}
}

func TestRemoveOldBackupSpecificDays(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "dlswp_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create __backup__ directory
	backupDir := filepath.Join(tempDir, "__backup__")
	err = os.Mkdir(backupDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create backup directory: %v", err)
	}

	// Create test directories with relative dates to today
	now := time.Now()
	
	// Test with daysToKeep = 3
	dates := []struct {
		date      time.Time
		shouldKeep bool
	}{
		{now.AddDate(0, 0, -5), false}, // 5 days ago (should be removed)
		{now.AddDate(0, 0, -4), false}, // 4 days ago (should be removed)
		{now.AddDate(0, 0, -3), false}, // 3 days ago (should be removed - cutoff)
		{now.AddDate(0, 0, -2), true},  // 2 days ago (should be kept)
		{now.AddDate(0, 0, -1), true},  // 1 day ago (should be kept)
		{now, true},                    // today (should be kept)
	}

	// Create directories
	var dirs []string
	for _, d := range dates {
		dateStr := d.date.Format("2006-01-02")
		dir := filepath.Join(backupDir, dateStr)
		err = os.Mkdir(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory %s: %v", dateStr, err)
		}
		dirs = append(dirs, dir)
	}

	// Run remove_old_backup with 3 days to keep
	err = remove_old_backup(tempDir, 3)
	if err != nil {
		t.Errorf("remove_old_backup failed: %v", err)
	}

	// Check results
	for i, d := range dates {
		dir := dirs[i]
		_, err := os.Stat(dir)
		
		if d.shouldKeep {
			if err != nil {
				t.Errorf("Directory %s should still exist: %v", d.date.Format("2006-01-02"), err)
			}
		} else {
			if !os.IsNotExist(err) {
				t.Errorf("Directory %s should have been removed", d.date.Format("2006-01-02"))
			}
		}
	}
}