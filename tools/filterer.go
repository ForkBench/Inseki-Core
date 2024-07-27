package tools

import (
	"os"
	"path/filepath"
)

type Association struct {
	Pattern string
	Nodes   []*Node
}

// This is a function that we can use with exploreFolder to filter files and folders
func FilterWithPatternMap(patterns *[]Association, callback func(filepath string, nodes []*Node)) func(path string, info os.FileInfo) error {
	return func(path string, info os.FileInfo) error {

		// Pattern could be for example :
		// { "*.c": [Node{...}, Node{...}], "*.h": [Node{...}, Node{...}] }

		// Or we could have a pattern like this:
		// { "TP*": [Node{...}, Node{...}], "TD*": [Node{...}, Node{...}] }

		// TODO: Add order

		// If it is a file, check if the file is in the patterns
		for _, association := range *patterns {
			// If the path matches the pattern
			if match, _ := filepath.Match(association.Pattern, filepath.Base(path)); match {

				// Call the callback
				callback(path, association.Nodes)

				// If it's a directory, we don't need to go deeper
				if info.IsDir() {
					return filepath.SkipDir
				}

			}
		}

		return nil
	}
}
