package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var backup_dir_pattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func getDefaultDownloadsPath() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("USERPROFILE"), "Downloads")
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Downloads")
	case "linux":
		return filepath.Join(os.Getenv("HOME"), "Downloads")
	default:
		return filepath.Join(os.Getenv("HOME"), "Downloads")
	}
}

func main() {
	// Check for command line arguments
	diff := "0"
	root := getDefaultDownloadsPath()

	// Parse the argument
	if len(os.Args) == 2 {
		diff = os.Args[1]
	} else if len(os.Args) == 3 {
		diff = os.Args[1]
		root = os.Args[2]
	}

	days, err := strconv.Atoi(diff)
	if err != nil {
		fmt.Println("Invalid argument:", err)
		return
	}

	if days < 0 {
		fmt.Println("Please specify a value of 0 or greater")
		return
	}

	// 引数が0の場合は4として扱う
	if days == 0 {
		days = 4
	}

	// download → backup へ移動
	if err := move_downloads_to_backup(root); err != nil {
		fmt.Println("Error moving files to backup:", err)
		return
	}
	// backup 内の古いディレクトリを削除
	if err := remove_old_backup(root, days); err != nil {
		fmt.Println(err)
	}
}

func remove_old_backup(root string, daysToKeep int) error {
	path := filepath.Join(root, "__backup__")
	dirs, err := getDirPaths(path)
	if err != nil {
		return err
	}

	for _, v := range dirs {
		if !backup_dir_pattern.MatchString(v) {
			continue
		}

		d, err := time.Parse("2006-01-02", v)
		if err != nil {
			fmt.Println("failed to convert : "+v, err)
			continue
		}

		// 今日から指定日数分だけ過去の日付を削除対象とする
		cutoffDate := time.Now().AddDate(0, 0, -daysToKeep)
		if !d.Before(cutoffDate) {
			continue
		}

		delPath := filepath.Join(root, "__backup__", v)
		fmt.Println(delPath)

		err = os.RemoveAll(delPath)
		if err != nil {
			fmt.Printf("Failed to remove directory: %v\n", err)
		} else {
			fmt.Println("Directory successfully removed.")
		}
	}

	return nil
}

func move_downloads_to_backup(root string) error {

	date := time.Now().Format("2006-01-02")

	fmt.Println("root :", root)
	fmt.Println("date :", date)
	fmt.Println("")

	// Move paths created on the target date to the new folder
	targetDir := ""
	paths, err := getFilePaths(root)
	if err != nil {
		return fmt.Errorf("failed to get file paths: %w", err)
	}

	for _, path := range paths {
		stat, err := os.Stat(path)
		if err != nil {
			fmt.Println(path)
			fmt.Println(err)
			continue
		}

		// Create a folder with the target date in the download folder
		if targetDir == "" {
			targetDir = filepath.Join(root, "__backup__", date)
			if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
				return fmt.Errorf("error creating folder %s: %w", targetDir, err)
			}
		}

		newPath := filepath.Join(targetDir, stat.Name())
		err = os.Rename(path, newPath)

		separator := string(filepath.Separator)
		pathSeparator := root + separator
		fmt.Println(strings.Replace(path, pathSeparator, "", 1))
		fmt.Println("  →", strings.Replace(newPath, pathSeparator, "", 1))
		if err != nil {
			fmt.Println("  ❗Error moving file:", err)
		}
	}

	return nil
}

func getFilePaths(baseDir string) ([]string, error) {
	files, err := os.ReadDir(baseDir)

	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", baseDir, err)
	}

	var paths []string
	for _, file := range files {
		// __ 始まりはスキップ
		if strings.HasPrefix(file.Name(), "__") {
			continue
		}
		path := filepath.Join(baseDir, file.Name())
		paths = append(paths, path)
	}

	return paths, nil
}

func getDirPaths(baseDir string) ([]string, error) {
	files, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, err
	}

	paths := []string{}
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, file.Name())
		}
	}

	return paths, nil
}
