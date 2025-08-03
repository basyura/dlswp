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
	paths := getFilePaths(tempDir)

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

	// Test with current date (should not remove anything)
	err = remove_old_backup(tempDir, time.Now())
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
	oldDate := now.AddDate(0, 0, -5) // 5 days ago (should be removed)
	recentDate := now.AddDate(0, 0, -2) // 2 days ago (should be kept)

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

	// Run remove_old_backup
	err = remove_old_backup(tempDir, now)
	if err != nil {
		t.Errorf("remove_old_backup failed: %v", err)
	}

	// Check that old directory was removed
	if _, err := os.Stat(oldDir); !os.IsNotExist(err) {
		t.Errorf("Old directory %s should have been removed", oldDateStr)
	}

	// Check that recent directory still exists
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
	err = remove_old_backup(tempDir, time.Now())
	if err == nil {
		t.Error("Expected error for non-existent backup directory, got nil")
	}
}