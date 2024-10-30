package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

type Association struct {
	Pattern    string
	Structures []Structure
}

type Target struct {
	Filepath    string
	Association Association
}

type Response struct {
	Filepath  string
	Structure Structure
}

func (r Response) String() string {
	return fmt.Sprintf("Filepath: %s, Structure: %s", r.Filepath, r.Structure.Name)
}

func (t Target) String() string {
	str := fmt.Sprintf("Path: %s, Structures: [", t.Filepath)

	for i, structure := range t.Association.Structures {
		str += structure.Name

		if i < len(t.Association.Structures)-1 {
			str += ", "
		}
	}

	str += "]"

	return str
}

// FilterWithPatternMap : This is a function that we can use with exploreFolder to filter files and folders
func FilterWithPatternMap(patterns *[]Association, stack *Stack) func(path string, info os.FileInfo) error {
	return func(path string, info os.FileInfo) error {

		// Pattern could be for example :
		// { "*.c": [Node{...}, Node{...}], "*.h": [Node{...}, Node{...}] }

		// Or we could have a pattern like this:
		// { "TP*": [Node{...}, Node{...}], "TD*": [Node{...}, Node{...}] }

		// TODO: Add order

		for _, association := range *patterns {
			// If the path matches the pattern
			if match, _ := filepath.Match(association.Pattern, filepath.Base(path)); match {

				// Add the path to the stack
				stack.Push(Target{
					Filepath:    path,
					Association: association,
				})

				// If it's a directory, we don't need to go deeper
				if info.IsDir() {
					return filepath.SkipDir
				}

			}
		}

		return nil
	}
}
