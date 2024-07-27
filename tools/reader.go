package tools

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// Read the .insekiignore file
func ReadInsekiIgnore(config Config) []string {
	insekiIgnorePath := filepath.Join(config.InsekiPath, ".insekiignore")

	// Read the .insekiignore file
	insekiIgnore, err := os.ReadFile(insekiIgnorePath)
	if err != nil {
		// If the file does not exist, return an empty slice
		if os.IsNotExist(err) {
			return []string{}
		}

		// If there is an error, panic
		panic(err)
	}

	// Split the file by line
	lines := strings.Split(string(insekiIgnore), "\n")

	// Remove empty lines
	var result []string
	for _, line := range lines {
		if line != "" {
			result = append(result, line)
		}
	}

	return result
}

func TranslateDir(path string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir

	if path == "~" {
		// In case of "~", which won't be caught by the "else if"
		path = dir
	} else if strings.HasPrefix(path, "~/") {
		// Use strings.HasPrefix so we don't match paths like
		// "/something/~/something/"
		path = filepath.Join(dir, path[2:])
	}

	return path
}

func ExploreFolder(path string, insekiignore []string, callback func(path string, info os.FileInfo) error, numberFilesAnalysed *int) error {

	// Translate the path
	path = TranslateDir(path)

	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {

			// If the error is "operation not permitted", we can ignore it
			if os.IsPermission(err) {
				return nil
			}

			fmt.Println("Error: ", err)
			return err
		}

		if !info.IsDir() {
			*numberFilesAnalysed++
		}

		// Check if the file or folder name is in the .insekiignore
		for _, ignore := range insekiignore {
			if filepath.Base(path) == ignore {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		return callback(path, info)
	})
}
