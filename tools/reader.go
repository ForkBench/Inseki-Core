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

// This is a function that we can use with exploreFolder to filter files and folders
func FilterWithPatternMap(patterns map[string][]Node, callback func(filepath string, nodes []Node)) func(path string, info os.FileInfo) error {
	return func(path string, info os.FileInfo) error {

		// Pattern could be for example :
		// { "*.c": [Node{...}, Node{...}], "*.h": [Node{...}, Node{...}] }

		// Or we could have a pattern like this:
		// { "TP*": [Node{...}, Node{...}], "TD*": [Node{...}, Node{...}] }

		// TODO: Add order

		// If it is a file, check if the file is in the patterns
		for pattern, nodes := range patterns {
			// If the path matches the pattern
			if match, _ := filepath.Match(pattern, filepath.Base(path)); match {

				// Call the callback
				callback(path, nodes)

				// If it's a directory, we don't need to go deeper
				if info.IsDir() {
					return filepath.SkipDir
				}

			}

		}

		return nil
	}
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

func ExploreFolder(path string, insekiignore []string, callback func(path string, info os.FileInfo) error) error {

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
