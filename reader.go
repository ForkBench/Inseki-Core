package inseki

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// ReadInsekiIgnore Read the .insekiignore file
func ReadInsekiIgnore(config Config) (error, []string) {
	insekiIgnorePath := filepath.Join(config.InsekiPath, ".insekiignore")

	// Read the .insekiignore file
	insekiIgnore, err := os.ReadFile(insekiIgnorePath)
	if err != nil {
		// If the file does not exist, return an empty slice
		if os.IsNotExist(err) {
			return nil, []string{}
		}

		// If there is an error, return it
		return err, nil
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

	return nil, result
}

// translateDir Translate a dir name to an absolute path (if there is ~)
func translateDir(path string) string {
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

// ExploreFolder Analyze for structures
func ExploreFolder(path string, insekiIgnore []string, callback func(path string, info os.FileInfo) error, numberFilesAnalysed *int) error {

	// Translate the path
	path = translateDir(path)

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

		// Check if the file or folder name is in the .insekiIgnore
		for _, ignore := range insekiIgnore {
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
