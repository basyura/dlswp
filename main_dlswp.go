package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Check for command line arguments
	diff := "0"
	root := filepath.Join(os.Getenv("USERPROFILE"), "Downloads")
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

	date := time.Now().AddDate(0, 0, days).Format("2006-01-02")

	fmt.Println("root :", root)
	fmt.Println("date :", date)
	fmt.Println("")

	// Move paths created on the target date to the new folder
	targetDir := ""
	paths := getFilePaths(root)
	for _, path := range paths {
		stat, err := os.Stat(path)
		if err != nil {
			fmt.Println(path)
			fmt.Println(err)
			continue
		}

		if stat.ModTime().Format("2006-01-02") == date {
			// Create a folder with the target date in the download folder
			if targetDir == "" {
				targetDir = filepath.Join(root, "__backup__", date)
				if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
					fmt.Println("Error creating folder:", err)
					return
				}
			}
			newPath := filepath.Join(targetDir, stat.Name())
			err := os.Rename(path, newPath)
			fmt.Println(strings.Replace(path, root+"\\", "", 1))
			fmt.Println("  →", strings.Replace(newPath, root+"\\", "", 1))
			if err != nil {
				fmt.Println("  ❗Error moving file:", err)
			}
		}
	}
}

func getFilePaths(baseDir string) []string {
	files, err := ioutil.ReadDir(baseDir)

	if err != nil {
		fmt.Println("read error :", baseDir)
		os.Exit(1)
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

	return paths
}
